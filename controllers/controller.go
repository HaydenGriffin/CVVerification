package controllers

import (
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Controller struct {
	Fabric *blockchain.FabricSetup
}

// basicAuth used to check the authentication (using basic auth) and retrieve the blockchain user
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

		session := sessions.InitSession(r)
		session.Values["LoggedInFlag"] = true

		// Check that there is corresponding user details stored in DB
		userDetails, err := database.GetUserDetailsFromUsername(pair[0])
		if err != nil {
			session.Values["SavedUserDetails"] = false
			err = session.Save(r, w)
			if err != nil {
				fmt.Println(err.Error())
			}
			pass(w, r, u)
		} else {
			gob.Register(userDetails)
			session.Values["SavedUserDetails"] = true
			session.Values["UserDetails"] = userDetails
			err = session.Save(r, w)
			if err != nil {
				fmt.Println(err.Error())
			}
			pass(w, r, u)
		}
	}
}

// Logout
func (c *Controller) LogoutHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {
		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "index",
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)

		session.Values["UserDetails"] = models.UserDetails{}
		session.Values["LoggedInFlag"] = false
		err := session.Save(r, w)
		if err != nil {
			fmt.Println(err.Error())
		}

		data.MessageSuccess = "You have been successfully logged out."
		renderTemplate(w, r, "index.html", data)
	})
}


func renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	lp := filepath.Join("web", "templates", "layout.html")
	ap := filepath.Join("web", "templates", "alerts.html")
	tp := filepath.Join("web", "templates", templateName)

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
