package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/cvverification/app/crypto"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (c *Controller) CVsToReviewView() func(http.ResponseWriter, *http.Request) {
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

		// Check that the user connected is a verifier
		_, err = u.QueryVerifier()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "You must be a verifier user to review CVs."
			renderTemplate(w, r, "index.html", data)
			return
		}

		industryFilter := r.FormValue("industry")

		cvList, err := u.QueryCVs(model.CVInReview, industryFilter)
		if err != nil {
			data.MessageWarning = "An error occurred whilst retrieving CVs to review."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CVList = cvList

		if len(data.CVInfo.CVList) == 0 {
			data.MessageWarning = "There are no CVs to be reviewed at this time."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if industryFilter != "" {
			data.MessageSuccess = "Showing results for " + industryFilter
		}

		data.CurrentPage = "viewcvs"
		renderTemplate(w, r, "cvstoreview.html", data)
	})
}

func (c *Controller) ReviewCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "cvstoreview",
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

		// Check that the user connected is an admin
		verifier, err := u.QueryVerifier()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "You must be a verifier user to rate a CV."
			renderTemplate(w, r, "index.html", data)
			return
		} else {
			data.UserDetails.Organisation = verifier.Profile.Organisation
		}

		cvID, success := mux.Vars(r)["cvID"]
		if !success {
			data.MessageWarning = "Error! No CV to be retrieved."
			renderTemplate(w, r, "index.html", data)
			return
		}

		applicantFabricID, err := database.GetFabricIDFromCVID(cvID)
		if err != nil {
			data.MessageWarning = "Error! Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cv, err := u.QueryCV(cvID)
		if err != nil {
			data.MessageWarning = "Error! Unable to find CV from hash."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if cv.Status != model.CVInReview {
			data.MessageWarning = "Error! Unable to review CV: CV is not in review."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = cv
		session.Values["ApplicantFabricID"] = applicantFabricID
		session.Values["CVID"] = cvID
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}
		data.MessageWarning = "If you have already reviewed this CV, your review will be overwritten."
		renderTemplate(w, r, "reviewcv.html", data)
	})
}

func (c *Controller) ReviewCVHandler() func(http.ResponseWriter, *http.Request) {
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

		// Check that the user connected is an admin
		_, err = u.QueryVerifier()
		if err != nil {
			data.MessageWarning = "Error! You must be a verifier user to review a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		ratingInt, err := strconv.Atoi(r.FormValue("rating"))
		if err != nil {
			data.MessageWarning = "Error! Rating must be a number."
			renderTemplate(w, r, "index.html", data)
			return
		}

		review := model.CVReview{
			Name:    r.FormValue("name"),
			Organisation: r.FormValue("organisation"),
			Comment: r.FormValue("comment"),
			Type:    r.FormValue("type"),
			Rating:  ratingInt,
		}

		applicantID := getApplicantFabricID(session)
		cvID := getCVID(session)

		if applicantID == "" || cvID == "" {
			data.MessageWarning = "Error! Unable to retrieve CV information"
			renderTemplate(w, r, "index.html", data)
			return
		}

		reviewByte, err := json.Marshal(review)
		if err != nil {
			data.MessageWarning = "Error! Unable to save review in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Retrieve the applicants public key
		applicantKeyString, err := u.QueryApplicantKey(applicantID)

		// Convert the byte representation of the key to *rsa.PublicKey
		applicantKey := crypto.BytesToPublicKey([]byte(applicantKeyString))

		// Encrypt the reviewByte with the applicants public key
		encryptedReview := crypto.EncryptWithPublicKey(reviewByte, applicantKey)

		// Save the rating to the applicants profile
		err = u.UpdateVerifierSaveReview(applicantID, cvID, encryptedReview)
		if err != nil {
			data.MessageWarning = "Error! Unable to save rating in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Query all CVs to display
		cvList, err := u.QueryCVs(model.CVInReview, "")
		if err != nil {
			data.MessageWarning = "An error occurred whilst retrieving CVs to review."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CVList = cvList
		data.CurrentPage = "cvstoreview"
		data.MessageSuccess = "Success! Your review has been saved."

		renderTemplate(w, r, "cvstoreview.html", data)
	})
}
