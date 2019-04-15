package sessions

import (
	"encoding/gob"
	"fmt"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/blockchain"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
)


type SessionSetup struct {
	Store *sessions.CookieStore
	Session *sessions.Session
}

func (s *SessionSetup) InitSession() {

	var authKey = securecookie.GenerateRandomKey(64)
	var encryptionKey = securecookie.GenerateRandomKey(32)

	s.Store = sessions.NewCookieStore(authKey, encryptionKey)

	s.Store.Options = &sessions.Options{
		MaxAge:   60 * 30,
		HttpOnly: true,
	}
}

func (s *SessionSetup) GetSession(r *http.Request, w http.ResponseWriter, u *blockchain.User) *sessions.Session {
	// Calling store.Get will return an empty session if the session does not already exist
	session, err := s.Store.Get(r, "userSession")
	if err != nil {
		fmt.Println(err)
	}

	if session.IsNew {
		session.Options.Domain = "localhost"
		session.Options.Path = "/"
	}

	if session.Values["ActiveSession"] != true {
		fmt.Println("active session false")
		// Check that there is corresponding user details stored in DB
		userDetails, err := database.GetUserDetailsFromUsername(u.Username)
		if err == nil {
			session.Values["SavedUserDetails"] = true
		} else {
			session.Values["SavedUserDetails"] = false
		}

		var accountType string

		applicant, err := u.QueryApplicant()
		if err == nil {
			if len(applicant.Profile.CVHistory) > 0 {
				userDetails.UploadedCV = true
			}
			accountType = "applicant"
		}

		_, err = u.QueryVerifier()
		if err == nil {
			accountType = "verifier"
		}

		_, err = u.QueryAdmin()
		if err == nil {
			accountType = "admin"
		}

		session.Values["AccountType"] = accountType
		gob.Register(userDetails)
		session.Values["UserDetails"] = userDetails
		session.Values["ActiveSession"] = true
	}

	fmt.Println(session.Values)

	session.Options.MaxAge = 60 * 30

	err = session.Save(r, w)
	if err != nil {
		fmt.Println("session save error (from controller")
		fmt.Println(err.Error())
	}

	return session
}



func HasSavedUserDetails(s *sessions.Session) bool {
	saved := s.Values["SavedUserDetails"]

	if saved != true {
		return false
	} else {
		return true
	}
}

func GetUserDetails(s *sessions.Session) templateModel.UserDetails {
	val := s.Values["UserDetails"]

	userDetails, ok := val.(templateModel.UserDetails)
	if !ok {
		return templateModel.UserDetails{}
	}

	return userDetails
}

func GetPrivateKey(s *sessions.Session) string {
	val := s.Values["PrivateKey"]

	privateKey, ok := val.(string)
	if !ok {
		return ""
	}
	return privateKey
}

func GetAccountType(s *sessions.Session) string {
	val := s.Values["AccountType"]

	accountType, ok := val.(string)
	if !ok {
		return ""
	}

	return accountType
}

func GetCVID(s *sessions.Session) string {
	val := s.Values["CVID"]

	cvID, ok := val.(string)
	if !ok {
		return ""
	}
	return cvID
}

func GetApplicantFabricID(s *sessions.Session) string {
	val := s.Values["ApplicantFabricID"]

	ID, ok := val.(string)
	if !ok {
		return ""
	}
	return ID
}
