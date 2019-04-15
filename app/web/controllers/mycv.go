package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/cvverification/app/crypto"
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

		session := sessions.GetSession(r)

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = sessions.GetAccountType(session)
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
		var cvHistory templateModel.CVHistoryInfo

		for index, cvID := range applicant.Profile.CVHistory {

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

		session.Values["CVID"] = cvIDToDisplay
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

		fmt.Println("reviews")
		fmt.Println(applicant.Profile.Reviews)

		var reviews []model.CVReview

		if applicant.Profile.Reviews[cvIDToDisplay] != nil {
			privateKeyString := sessions.GetPrivateKey(session)
			privateKey := crypto.BytesToPrivateKey([]byte(privateKeyString))
			encryptedCVReviews := applicant.Profile.Reviews[cvIDToDisplay]

			var review model.CVReview

			for _, encryptedReview := range encryptedCVReviews {
				decryptedReviewByte := crypto.DecryptWithPrivateKey(encryptedReview, privateKey)
				err = json.Unmarshal(decryptedReviewByte, &review)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(review)
				reviews = append(reviews, review)
			}
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

		session := sessions.GetSession(r)

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = sessions.GetAccountType(session)
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

		cvIDToUpdate := sessions.GetCVID(session)
		var allCVHistory []templateModel.CVHistoryInfo
		var cvHistory templateModel.CVHistoryInfo

		for index, cvID := range applicant.Profile.CVHistory {

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

		if cvIDToUpdate == "" || len(allCVHistory) == 0 {
			data.MessageWarning = "Error! Unable to retrieve CV info."
			renderTemplate(w, r, "index.html", data)
			return
		}

		updatedCV, err := u.UpdateTransitionCV(cvIDToUpdate, model.CVInReview)
		if err != nil {
			data.MessageWarning = "Error! Unable to transition CV status in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		var reviews []model.CVReview

		if applicant.Profile.Reviews[cvIDToUpdate] != nil {
			privateKeyString := sessions.GetPrivateKey(session)
			privateKey := crypto.BytesToPrivateKey([]byte(privateKeyString))
			encryptedCVReviews := applicant.Profile.Reviews[cvIDToUpdate]

			var review model.CVReview

			for _, encryptedReview := range encryptedCVReviews {
				decryptedReviewByte := crypto.DecryptWithPrivateKey(encryptedReview, privateKey)
				err = json.Unmarshal(decryptedReviewByte, &review)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(review)
				reviews = append(reviews, review)
			}
		}

		data.CVInfo.CV = updatedCV
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

		session := sessions.GetSession(r)

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = sessions.GetAccountType(session)
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

		cvIDToUpdate := sessions.GetCVID(session)
		var allCVHistory []templateModel.CVHistoryInfo
		var cvHistory templateModel.CVHistoryInfo

		for index, cvID := range applicant.Profile.CVHistory {

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

		if cvIDToUpdate == "" || len(allCVHistory) == 0 {
			data.MessageWarning = "Error! Unable to retrieve CV info."
			renderTemplate(w, r, "index.html", data)
			return
		}

		updatedCV, err := u.UpdateTransitionCV(cvIDToUpdate, model.CVInDraft)
		if err != nil {
			data.MessageWarning = "Error! Unable to transition CV status in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		for index, cvHistory := range allCVHistory {
			if cvHistory.CVID == cvIDToUpdate {
				allCVHistory[index].CV = updatedCV
			}
		}

		var reviews []model.CVReview

		if applicant.Profile.Reviews[cvIDToUpdate] != nil {
			privateKeyString := sessions.GetPrivateKey(session)
			privateKey := crypto.BytesToPrivateKey([]byte(privateKeyString))
			encryptedCVReviews := applicant.Profile.Reviews[cvIDToUpdate]

			var review model.CVReview

			for _, encryptedReview := range encryptedCVReviews {
				decryptedReviewByte := crypto.DecryptWithPrivateKey(encryptedReview, privateKey)
				err = json.Unmarshal(decryptedReviewByte, &review)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(review)
				reviews = append(reviews, review)
			}
		}


		data.CVInfo.CV = updatedCV
		data.CVInfo.CurrentCVID = cvIDToUpdate
		data.CVInfo.Reviews = reviews
		data.CVInfo.CVHistory = allCVHistory
		data.MessageSuccess = "Success! Your CV has been withdrawn from review."
		data.CurrentPage = "mycv"
		renderTemplate(w, r, "mycv.html", data)
	})
}
