package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/cvverification/app/crypto"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (c *Controller) MyCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
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

		var reviews []model.CVReview

		if data.CVInfo.CV.Status == model.CVSubmitted {
			reviews = applicant.Profile.PublicReviews[cvIDToDisplay]
		} else if applicant.Profile.Reviews[cvIDToDisplay] != nil {
			// If there is reviews on the CV that is being displayed
			var review model.CVReview
			encryptedCVReviews := applicant.Profile.Reviews[cvIDToDisplay]

			// Get private key from session
			privateKeyString := getPrivateKey(session)
			if privateKeyString != "" {
				// Convert bytes of private key to *rsa.PrivateKey
				privateKey := crypto.BytesToPrivateKey([]byte(privateKeyString))

				// Loop over each review
				for _, encryptedReview := range encryptedCVReviews {
					// Attempt to decrypt the encrypted review with the private key
					decryptedReviewByte, err := crypto.DecryptWithPrivateKey(encryptedReview, privateKey)
					if err != nil {
						fmt.Println(err)
						data.CVInfo.ReviewInfo.Status = "decrypterr"
						data.MessageWarning = "Error! It looks like you have uploaded the incorrect Private Key."
						continue
					}
					err = json.Unmarshal(decryptedReviewByte, &review)
					if err != nil {
						fmt.Println(err)
						data.CVInfo.ReviewInfo.Status = "decrypterr"
						data.MessageWarning = "Error! It looks like you have uploaded the incorrect Private Key."
						continue
					}
					reviews = append(reviews, review)
				}
			} else {
				data.CVInfo.ReviewInfo.Status = "nokey"
				data.MessageWarning = "You have reviews on this CV. Please upload your Private Key to view the reviews."
			}
		}

		data.CVInfo.ReviewInfo.Reviews = reviews
		data.CVInfo.CVHistory = allCVHistory
		data.CVInfo.CurrentCVID = cvIDToDisplay
		data.CurrentPage = "mycv"
		renderTemplate(w, r, "mycv.html", data)
	})
}

func (c *Controller) SubmitForReviewHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
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

		cvIDToUpdate := getCVID(session)
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

		for index, cvHistory := range allCVHistory {
			if cvHistory.CVID == cvIDToUpdate {
				allCVHistory[index].CV = updatedCV
			}
		}

		var reviews []model.CVReview

		if applicant.Profile.Reviews[cvIDToUpdate] != nil {
			// If there is reviews on the CV that is being displayed
			var review model.CVReview
			encryptedCVReviews := applicant.Profile.Reviews[cvIDToUpdate]

			// Get private key from session
			privateKeyString := getPrivateKey(session)
			if privateKeyString != "" {
				// Convert bytes of private key to *rsa.PrivateKey
				privateKey := crypto.BytesToPrivateKey([]byte(privateKeyString))

				// Loop over each review
				for _, encryptedReview := range encryptedCVReviews {
					// Attempt to decrypt the encrypted review with the private key
					decryptedReviewByte, err := crypto.DecryptWithPrivateKey(encryptedReview, privateKey)
					if err != nil {
						fmt.Println(err)
						data.CVInfo.ReviewInfo.Status = "decrypterr"
						data.MessageWarning = "Error! It looks like you have uploaded the incorrect Private Key."
						continue
					}
					err = json.Unmarshal(decryptedReviewByte, &review)
					if err != nil {
						fmt.Println(err)
						data.CVInfo.ReviewInfo.Status = "decrypterr"
						data.MessageWarning = "Error! It looks like you have uploaded the incorrect Private Key."
						continue
					}
					reviews = append(reviews, review)
				}
			} else {
				data.CVInfo.ReviewInfo.Status = "nokey"
				data.MessageWarning = "You have reviews on this CV. Please upload your Private Key to view the reviews."
			}
		}

		data.CVInfo.CV = updatedCV
		data.CVInfo.CurrentCVID = cvIDToUpdate
		data.CVInfo.ReviewInfo.Reviews = reviews
		data.CVInfo.CVHistory = allCVHistory
		data.MessageSuccess = "Success! Your CV can now be reviewed."
		data.CurrentPage = "mycv"
		renderTemplate(w, r, "mycv.html", data)
	})
}

func (c *Controller) WithdrawFromReviewHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
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

		cvIDToUpdate := getCVID(session)
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
			// If there is reviews on the CV that is being displayed
			var review model.CVReview
			encryptedCVReviews := applicant.Profile.Reviews[cvIDToUpdate]

			// Get private key from session
			privateKeyString := getPrivateKey(session)
			if privateKeyString != "" {
				// Convert bytes of private key to *rsa.PrivateKey
				privateKey := crypto.BytesToPrivateKey([]byte(privateKeyString))

				// Loop over each review
				for _, encryptedReview := range encryptedCVReviews {
					// Attempt to decrypt the encrypted review with the private key
					decryptedReviewByte, err := crypto.DecryptWithPrivateKey(encryptedReview, privateKey)
					if err != nil {
						fmt.Println(err)
						data.CVInfo.ReviewInfo.Status = "decrypterr"
						data.MessageWarning = "Error! It looks like you have uploaded the incorrect Private Key."
						continue
					}
					err = json.Unmarshal(decryptedReviewByte, &review)
					if err != nil {
						fmt.Println(err)
						data.CVInfo.ReviewInfo.Status = "decrypterr"
						data.MessageWarning = "Error! It looks like you have uploaded the incorrect Private Key."
						continue
					}
					reviews = append(reviews, review)
				}
			} else {
				data.CVInfo.ReviewInfo.Status = "nokey"
				data.MessageWarning = "You have reviews on this CV. Please upload your Private Key to view the reviews."
			}
		}

		data.CVInfo.CV = updatedCV
		data.CVInfo.CurrentCVID = cvIDToUpdate
		data.CVInfo.ReviewInfo.Reviews = reviews
		data.CVInfo.CVHistory = allCVHistory
		data.MessageSuccess = "Success! Your CV has been withdrawn from review."
		data.CurrentPage = "mycv"
		renderTemplate(w, r, "mycv.html", data)
	})
}

func (c *Controller) SubmitToEmployerHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
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

		cvIDToUpdate := getCVID(session)
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

		var reviews []model.CVReview

		if applicant.Profile.Reviews[cvIDToUpdate] != nil {
			// If there is reviews on the CV that is being displayed
			var review model.CVReview
			encryptedCVReviews := applicant.Profile.Reviews[cvIDToUpdate]

			// Get private key from session
			privateKeyString := getPrivateKey(session)
			if privateKeyString != "" {
				// Convert bytes of private key to *rsa.PrivateKey
				privateKey := crypto.BytesToPrivateKey([]byte(privateKeyString))

				// Loop over each review
				for _, encryptedReview := range encryptedCVReviews {
					// Attempt to decrypt the encrypted review with the private key
					decryptedReviewByte, err := crypto.DecryptWithPrivateKey(encryptedReview, privateKey)
					if err != nil {
						fmt.Println(err)
						data.CVInfo.ReviewInfo.Status = "decrypterr"
						data.MessageWarning = "Error! It looks like you have uploaded the incorrect Private Key."
						continue
					}
					err = json.Unmarshal(decryptedReviewByte, &review)
					if err != nil {
						fmt.Println(err)
						data.CVInfo.ReviewInfo.Status = "decrypterr"
						data.MessageWarning = "Error! It looks like you have uploaded the incorrect Private Key."
						continue
					}
					reviews = append(reviews, review)
				}
			} else {
				data.CVInfo.ReviewInfo.Status = "nokey"
				data.MessageWarning = "Error! Please upload your Private Key to submit your CV to employers."
				renderTemplate(w, r, "index.html", data)
				return
			}
		} else {
			data.MessageWarning = "Error! You must have at least one review to submit your CV to employers."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if len(reviews) == 0 {
			data.MessageWarning = "Error! Something went wrong whilst decrypting reviews. Please check your Private Key and try again."
			renderTemplate(w, r, "index.html", data)
			return
		}

		reviewsByte, err := json.Marshal(reviews)
		if err != nil {
			data.MessageWarning = "Error! Unable to publish CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		err = u.UpdatePublishReviews(cvIDToUpdate, reviewsByte)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to publish CV in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		updatedCV, err := u.UpdateTransitionCV(cvIDToUpdate, model.CVSubmitted)
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

		data.CVInfo.CV = updatedCV
		data.CVInfo.CurrentCVID = cvIDToUpdate
		data.CVInfo.ReviewInfo.Reviews = reviews
		data.CVInfo.CVHistory = allCVHistory
		data.MessageSuccess = "Success! Your CV has been submitted to employers."
		data.CurrentPage = "mycv"
		renderTemplate(w, r, "mycv.html", data)
	})
}
