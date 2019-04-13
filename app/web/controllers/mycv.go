package controllers

import (
	"encoding/gob"
	"fmt"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
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
			renderTemplate(w, r, "registerdetails.html", data)
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

		var requestedCVIndex int
		// Check for URL parameter (corresponds to index of a CV)
		// There may not be a passed URL parameter
		result, success := mux.Vars(r)["requestedCVIndex"]
		if success {
			// Convert string to int. If there is an error, set requestedCVIndex to 0 (default handling)
			requestedCVIndex, err = strconv.Atoi(result)
			if err != nil {
				requestedCVIndex = 0
			}
		}

		var allCVHistory []templateModel.CVHistoryInfo

		for index, cvID := range applicant.Profile.CVHistory {
			var cvHistory templateModel.CVHistoryInfo
			cv, err := u.QueryCV(cvID)
			if err != nil {
				data.MessageWarning = "Error! Unable to find CV info in ledger."
				renderTemplate(w, r, "index.html", data)
				return
			}
			cvHistory.Index = index + 1
			cvHistory.CVID = cvID
			cvHistory.CV = cv
			allCVHistory = append(allCVHistory, cvHistory)
		}

		var cvIDToDisplay string

		if requestedCVIndex != 0 {
			data.CVInfo.CV = allCVHistory[requestedCVIndex-1].CV
			cvIDToDisplay = allCVHistory[requestedCVIndex-1].CVID
		} else {
			data.CVInfo.CV = allCVHistory[len(allCVHistory)-1].CV
			cvIDToDisplay = allCVHistory[len(allCVHistory)-1].CVID
		}

		// Retrieve reviews for the CV that is to be displayed
		reviews := applicant.Profile.Reviews[cvIDToDisplay]

		gob.Register(data.CVInfo.CV)
		gob.Register(reviews)
		gob.Register(allCVHistory)


		session.Values["CVHistory"] = allCVHistory
		fmt.Println(allCVHistory)
		session.Values["CV"] = data.CVInfo.CV
		fmt.Println(data.CVInfo.CV)
		session.Values["CVID"] = cvIDToDisplay
		fmt.Println(cvIDToDisplay)
		session.Values["Reviews"] = reviews
		fmt.Println(reviews)
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.Reviews = reviews
		data.CVInfo.CVHistory = allCVHistory
		data.CVInfo.CurrentCVID = cvIDToDisplay
		data.CurrentPage = "mycv"
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
			renderTemplate(w, r, "registerdetails.html", data)
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
		cvIDToUpdate := sessions.GetCVID(session)
		reviews := sessions.GetReviews(session)
		allCVHistory := sessions.GetCVHistory(session)
		if cv == nil || cvIDToUpdate == "" || len(allCVHistory) == 0 {
			data.MessageWarning = "Error! Unable to retrieve CV info."
			renderTemplate(w, r, "index.html", data)
			return
		}

		err = u.UpdateTransitionCV(cvIDToUpdate, model.CVInReview)
		if err != nil {
			data.MessageWarning = "Error! Unable to transition CV status in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cv.Status = model.CVInReview
		for index, cvHistory := range allCVHistory {
			if cvHistory.CVID == cvIDToUpdate {
				allCVHistory[index].CV = cv
			}
		}

		session.Values["CVHistory"] = allCVHistory
		session.Values["CV"] = cv
		err = session.Save(r, w)
		if err != nil {
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = cv
		data.CVInfo.CurrentCVID = cvIDToUpdate
		data.CVInfo.Reviews = reviews
		data.CVInfo.CVHistory = allCVHistory
		data.MessageSuccess = "Success! Your CV can now be reviewed."
		data.CurrentPage = "mycv"
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
			renderTemplate(w, r, "registerdetails.html", data)
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
		cvID := sessions.GetCVID(session)
		reviews := sessions.GetReviews(session)
		allCVHistory := sessions.GetCVHistory(session)
		if cv == nil || cvID == "" || len(allCVHistory) == 0 {
			data.MessageWarning = "Error! Unable to retrieve CV info."
			renderTemplate(w, r, "index.html", data)
			return
		}

		err = u.UpdateTransitionCV(cvID, model.CVInDraft)
		if err != nil {
			data.MessageWarning = "Error! Unable to transition CV status in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cv.Status = model.CVInDraft
		for index, cvHistory := range allCVHistory {
			if cvHistory.CVID == cvID {
				allCVHistory[index].CV = cv
			}
		}

		session.Values["CVHistory"] = allCVHistory
		session.Values["CV"] = cv
		err = session.Save(r, w)
		if err != nil {
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = cv
		data.CVInfo.CurrentCVID = cvID
		data.CVInfo.Reviews = reviews
		data.CVInfo.CVHistory = allCVHistory
		data.MessageSuccess = "Success! Your CV has been withdrawn from review."
		data.CurrentPage = "mycv"
		renderTemplate(w, r, "mycv.html", data)
	})
}
