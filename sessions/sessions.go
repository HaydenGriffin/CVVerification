package sessions

import (
	"github.com/cvtracker/models"
	"github.com/gorilla/securecookie"
	"net/http"

	"github.com/gorilla/sessions"
)

var authKey = securecookie.GenerateRandomKey(64)
var encryptionKey = securecookie.GenerateRandomKey(32)

// store will hold all session data
var store = sessions.NewCookieStore(authKey, encryptionKey)

func InitSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "userSession")
	if session.IsNew {
		session.Options.Domain = "localhost"
		session.Options.Path = "/"
		session.Options.MaxAge = 0
		session.Options.HttpOnly = false
		session.Options.Secure = false
		session.Values["LoggedInFlag"] = false
	}
	return session
}

func IsLoggedIn(s *sessions.Session) bool {
	loggedIn := s.Values["LoggedInFlag"]

	if loggedIn != true {
		return false
	} else {
		return true
	}
}

func GetUser(s *sessions.Session) models.User {
	val := s.Values["User"]
	var user = models.User{}
	user, ok := val.(models.User)
	if !ok {
		return models.User{}
	}
	return user
}