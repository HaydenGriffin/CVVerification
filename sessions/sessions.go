package sessions

import (
	"github.com/cvverification/chaincode/model"
	"github.com/cvverification/models"
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

	userUploadedCV := GetUserUploadedCV(s)

	userDetails.UploadedCV = userUploadedCV

	return userDetails
}

func GetUserUploadedCV(s *sessions.Session) bool {
	val := s.Values["UserUploadedCV"]

	uploadedCV, ok := val.(bool)
	if !ok {
		return false
	}
	return uploadedCV
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

func GetApplicantFabricID(s *sessions.Session) string {
	val := s.Values["ApplicantFabricID"]

	ID, ok := val.(string)
	if !ok {
		return ""
	}
	return ID
}

func GetAllCVList(s *sessions.Session) map[int] *model.CVObject {
	val := s.Values["AllCVList"]

	allCVList, ok := val.(map[int] *model.CVObject)
	if !ok {
		return nil
	}
	return allCVList
}

func GetReviews(s *sessions.Session) []model.CVReview {
	val := s.Values["Reviews"]

	reviews, ok := val.([]model.CVReview)
	if !ok {
		return []model.CVReview{}
	}
	return reviews
}
