package sessions

import (
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

func HasSavedUserDetails(s *sessions.Session) bool {
	saved := s.Values["SavedUserDetails"]

	if saved != true {
		return false
	} else {
		return true
	}
}

func GetUserDetails(s *sessions.Session) models.UserDetails {
	val := s.Values["UserDetails"]

	userDetails, ok := val.(models.UserDetails)
	if !ok {
		return models.UserDetails{}
	}
	return userDetails
}

func GetCV(s *sessions.Session) *model.CVObject {
	val := s.Values["CV"]

	cv, ok := val.(*model.CVObject)
	if !ok {
		return nil
	}
	return cv
}

func GetCVHash(s *sessions.Session) string {
	val := s.Values["CVHash"]

	cvHash, ok := val.(string)
	if !ok {
		return ""
	}
	return cvHash
}

func GetProfileHash(s *sessions.Session) string {
	val := s.Values["ProfileHash"]

	profileHash, ok := val.(string)
	if !ok {
		return ""
	}
	return profileHash
}

func GetInReviewCVHash(s *sessions.Session) string {
	val := s.Values["InReviewCVHash"]

	cvHash, ok := val.(string)
	if !ok {
		return ""
	}
	return cvHash
}

func GetRatings(s *sessions.Session) []model.CVRating {
	val := s.Values["CV"]

	ratings, ok := val.([]model.CVRating)
	if !ok {
		return []model.CVRating{}
	}
	return ratings
}
