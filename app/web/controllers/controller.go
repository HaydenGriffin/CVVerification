package controllers

import (
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/teris-io/shortid"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Controller struct {
	Fabric  *blockchain.FabricSetup
	ShortID *shortid.Shortid
}

// store will hold all session data
var store *sessions.FilesystemStore

// Middleware that runs every time a request to access a page is received
// basicAuth is used to provide log in credentials to authenticate and retrieve blockchain user
func (c *Controller) basicAuth(pass func(http.ResponseWriter, *http.Request, *blockchain.User)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		u, err := c.Fabric.LogUser(pair[0], pair[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("authorization failed with error: %v", err), http.StatusUnauthorized)
			return
		}

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		if session.IsNew {
			setSessionValues(session, u)
		}
		session.Options.MaxAge = 60 * 30

		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			return
		}

		pass(w, r, u)
	}
}

// Logout
func (c *Controller) LogoutHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {
		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "index",
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)

		// Remove all session values currently set
		session.Options.MaxAge = -1

		for key, _ := range session.Values {
			delete(session.Values, key)
		}

		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.MessageSuccess = "Success! You have been logged out."
		renderTemplate(w, r, "index.html", data)
	})
}

func renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	lp := filepath.Join("app", "web", "templates", "layout.html")
	ap := filepath.Join("app", "web", "templates", "alerts.html")
	tp := filepath.Join("app", "web", "templates", templateName)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(tp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	resultTemplate, err := template.ParseFiles(tp, lp, ap)
	if err != nil {
		// Log the detailed error
		fmt.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err := resultTemplate.ExecuteTemplate(w, "layout", data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func (c *Controller) InitSession() {
	var authKey = securecookie.GenerateRandomKey(64)
	var encryptionKey = securecookie.GenerateRandomKey(32)

	store = sessions.NewFilesystemStore(
		"",
		authKey,
		encryptionKey,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 30,
		HttpOnly: true,
	}
}

func setSessionValues(session *sessions.Session, u *blockchain.User) {

	// Check that there is corresponding user details stored in DB
	userDetails, err := database.GetUserDetailsFromUsername(u.Username)
	if err == nil {
		session.Values["SavedUserDetails"] = true
	} else {
		session.Values["SavedUserDetails"] = false
	}

	// accountType specifies the current userType
	var actorType string

	applicant, err := u.QueryApplicant()
	if err == nil {
		if len(applicant.Profile.CVHistory) > 0 {
			userDetails.UploadedCV = true
		}
		actorType = model.ActorApplicant
	}

	_, err = u.QueryVerifier()
	if err == nil {
		actorType = model.ActorVerifier
	}

	_, err = u.QueryAdmin()
	if err == nil {
		actorType = model.ActorAdmin
	}

	_, err = u.QueryEmployer()
	if err == nil {
		actorType = model.ActorEmployer
	}

	session.Values["AccountType"] = actorType
	gob.Register(userDetails)
	session.Values["UserDetails"] = userDetails
}

func hasSavedUserDetails(s *sessions.Session) bool {
	saved := s.Values["SavedUserDetails"]

	if saved != true {
		return false
	} else {
		return true
	}
}

func getUserDetails(s *sessions.Session) templateModel.UserDetails {
	val := s.Values["UserDetails"]

	userDetails, ok := val.(templateModel.UserDetails)
	if !ok {
		return templateModel.UserDetails{}
	}

	return userDetails
}

func getPrivateKey(s *sessions.Session) string {
	val := s.Values["PrivateKey"]

	// Type assertion - ensure that val is a string value
	privateKey, ok := val.(string)
	if !ok {
		return ""
	}
	return privateKey
}

func getAccountType(s *sessions.Session) string {
	val := s.Values["AccountType"]

	accountType, ok := val.(string)
	if !ok {
		return ""
	}
	return accountType
}

func getCVID(s *sessions.Session) string {
	val := s.Values["CVID"]

	cvID, ok := val.(string)
	if !ok {
		return ""
	}
	return cvID
}

func getApplicantFabricID(s *sessions.Session) string {
	val := s.Values["ApplicantFabricID"]

	ID, ok := val.(string)
	if !ok {
		return ""
	}
	return ID
}
