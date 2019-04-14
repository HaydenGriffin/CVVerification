package sessions

import (
	templateModel "github.com/cvverification/app/model"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
)

var authKey = securecookie.GenerateRandomKey(64)
var encryptionKey = securecookie.GenerateRandomKey(32)

// store will hold all session data
var store = sessions.NewCookieStore(authKey, encryptionKey)

func GetSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "userSession")
	if session.IsNew {
		session.Options.Domain = "localhost"
		session.Options.Path = "/"
		// Session max age is in seconds. Max Age is set to 15 mins
		session.Options.MaxAge = 60*30
		session.Options.HttpOnly = false
		session.Options.Secure = false
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

