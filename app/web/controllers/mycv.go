package controllers

import (
	"encoding/gob"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (c *Controller) MyCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "userdetails.html", data)
			return
		}

		// Check that the user connected is an applicant
		applicant, err := u.QueryApplicant()
		if err != nil {
			data.MessageWarning = "Error! You must be an applicant user to view your CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Check that the user has uploaded at least one CV
		if len(applicant.Profile.CVHistory) == 0 {
			data.MessageWarning = "Error! You have not uploaded a CV yet. Please fill in the following form and add your CV."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		// Check for URL parameter (corresponds to index of a CV)
		// There may not be a passed URL parameter
		result, success := mux.Vars(r)["requestedCVIndex"]
		var requestedCVIndex int

		if success {
			// Convert string to int. If there is an error, set requestedCVIndex to 0 (default handling)
			requestedCVIndex, err = strconv.Atoi(result)
			if err != nil {
				requestedCVIndex = 0
			}
		}

		// Retrieve historical CV information, including current version
		// If there is an error or there is no CV history, exit
		historicalCVHistoryInfo, err := database.GetUserCVDetails(data.UserDetails.Id)
		if err != nil || len(historicalCVHistoryInfo) == 0 {
			data.MessageWarning = "Error! Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		var cvToDisplayCVHash string

		if requestedCVIndex != 0 {
			// User has requested a CV to display
			for i := range historicalCVHistoryInfo {
				if historicalCVHistoryInfo[i].Index == requestedCVIndex {
					// Found match for requestedCVIndex
					cvToDisplayCVHash = historicalCVHistoryInfo[i].CVHash
					if historicalCVHistoryInfo[i].CVInReview == 1 {
						data.CVInfo.UserHasCVInReview = true
						data.CVInfo.CurrentCVInReview = true
					}
					data.CVInfo.CurrentCVHash = historicalCVHistoryInfo[i].CVHash
				}
			}
		} else {
			// No CV to display requested. If the user has a CV in review, display this with priority
			for i := range historicalCVHistoryInfo {
				if historicalCVHistoryInfo[i].CVInReview == 1 {
					data.CVInfo.UserHasCVInReview = true
					cvToDisplayCVHash = historicalCVHistoryInfo[i].CVHash
					data.CVInfo.CurrentCVInReview = true
					data.CVInfo.CurrentCVHash = historicalCVHistoryInfo[i].CVHash
				}
			}
		}

		// No CV requested and no CV found in review - display the most recent version of the CV
		if cvToDisplayCVHash == "" {
			cvToDisplayCVHash = historicalCVHistoryInfo[len(historicalCVHistoryInfo)-1].CVHash
			data.CVInfo.CurrentCVHash = historicalCVHistoryInfo[len(historicalCVHistoryInfo)-1].CVHash
		}

		// Retrieve CV details from ledger
		cv, err := u.QueryCV(cvToDisplayCVHash)
		if err != nil {
			data.MessageWarning = "Error! Unable to retrieve CV details from ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = cv
		data.CVInfo.CVHistory = historicalCVHistoryInfo
		data.CurrentPage = "mycv"

		// Retrieve reviews for the CV that is to be displayed
		reviews := applicant.Profile.Reviews[cvToDisplayCVHash]
		data.CVInfo.Reviews = reviews

		gob.Register(cv)
		gob.Register(reviews)
		session.Values["CV"] = cv
		session.Values["CVHash"] = cvToDisplayCVHash
		session.Values["Reviews"] = reviews

		err = session.Save(r, w)
		if err != nil {
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

		renderTemplate(w, r, "mycv.html", data)
	})
}

func (c *Controller) SubmitForReviewHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "userdetails.html", data)
			return
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			data.MessageWarning = "Error! You must be an applicant user to view your CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cv := sessions.GetCV(session)
		cvHash := sessions.GetCVHash(session)
		reviews := sessions.GetReviews(session)
		if cv == nil || cvHash == "" {
			data.MessageWarning = "Error! Unable to update status of CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = cv
		data.CVInfo.CurrentCVHash = cvHash
		data.CVInfo.Reviews = reviews

		data.CVInfo.UserHasCVInReview = database.UserHasCVInReview(data.UserDetails.Id)

		if data.CVInfo.UserHasCVInReview {
			data.MessageWarning = "Error! You are only allowed one version of your CV in review at a time."
		} else {
			err := database.UpdateCV(cvHash, 1)
			if err != nil {
				data.MessageWarning = "Error! Unable to update CV info in database."
				renderTemplate(w, r, "index.html", data)
				return
			}
			data.MessageSuccess = "Success! Your CV can now be reviewed."
			data.CVInfo.CurrentCVInReview = true
			data.CVInfo.UserHasCVInReview = true
			data.CurrentPage = "mycv"
		}

		historicalCVHistoryInfo, err := database.GetUserCVDetails(data.UserDetails.Id)

		if err != nil || len(historicalCVHistoryInfo) == 0 {
			data.MessageWarning = "Error! Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CVHistory = historicalCVHistoryInfo
		renderTemplate(w, r, "mycv.html", data)
	})
}

func (c *Controller) WithdrawFromReviewHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "userdetails.html", data)
			return
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			data.MessageWarning = "Error! You must be an applicant user to view your CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cv := sessions.GetCV(session)
		cvHash := sessions.GetCVHash(session)
		reviews := sessions.GetReviews(session)
		if cv == nil || cvHash == "" {
			data.MessageWarning = "Error! Unable to update status of CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = cv
		data.CVInfo.CurrentCVHash = cvHash
		data.CVInfo.Reviews = reviews

		err = database.UpdateCV(cvHash, 0)
		if err != nil {
			data.MessageWarning = "Error! Unable to update CV info in database."
			renderTemplate(w, r, "index.html", data)
		}

		historicalCVHistoryInfo, err := database.GetUserCVDetails(data.UserDetails.Id)

		if err != nil || len(historicalCVHistoryInfo) == 0 {
			data.MessageWarning = "Error! Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CVHistory = historicalCVHistoryInfo
		data.MessageSuccess = "Success! Your CV has been withdrawn from review."
		data.CVInfo.CurrentCVInReview = false
		renderTemplate(w, r, "mycv.html", data)
	})
}
