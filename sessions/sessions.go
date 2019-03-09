package sessions

import (
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/chaincode/model"
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

func GetFabricUser(s *sessions.Session) *blockchain.User {
	val := s.Values["FabricUser"]

	user, ok := val.(*blockchain.User)
	if !ok || user == nil {
		return nil
	}
	return user
}

func GetUserDetails(s *sessions.Session) models.UserDetails {
	val := s.Values["FabricUser"]

	user, ok := val.(models.UserDetails)
	if !ok {
		return models.UserDetails{}
	}
	return user
}

func GetCV(s *sessions.Session) model.CVObject {
	val := s.Values["CV"]

	cv, ok := val.(model.CVObject)
	if !ok {
		return model.CVObject{}
	}
	return cv
}

func GetRatings(s *sessions.Session) []model.CVRating {
	val := s.Values["CV"]

	ratings, ok := val.([]model.CVRating)
	if !ok {
		return []model.CVRating{}
	}
	return ratings
}
