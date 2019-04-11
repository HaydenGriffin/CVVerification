package controllers

import (
	"encoding/gob"
	"encoding/json"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (c *Controller) ReviewCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "viewallcv",
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

		// Check that the user connected is an admin
		_, err := u.QueryVerifier()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "You must be a verifier user to rate a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		//data.CVInfo.CVList = sessions.GetAllCVList(session)

		result, success := mux.Vars(r)["userID"]
		if !success {
			data.MessageWarning = "Error! No CV to be retrieved."
			renderTemplate(w, r, "viewallcv.html", data)
			return
		}

		userID, err := strconv.Atoi(result)
		if err != nil {
			data.MessageWarning = "Error! Invalid CV ID."
			renderTemplate(w, r, "viewallcv.html", data)
			return
		}

		applicantFabricID, cvHash, err := database.GetCVInfoFromID(userID)
		if err != nil {
			data.MessageWarning = "Error! Unable to find CV info in database."
			renderTemplate(w, r, "viewallcv.html", data)
			return
		}

		verifierReview, err := u.QueryVerifierCVReview(applicantFabricID, cvHash)
		if err != nil {
			data.MessageWarning = "Error! Unable to find CV review information in ledger."
			renderTemplate(w, r, "viewallcv.html", data)
			return
		}

		data.CVInfo.Review = verifierReview

		cv, err := u.QueryCV(cvHash)
		if err != nil {
			data.MessageWarning = "Error! Unable to find CV from hash."
			renderTemplate(w, r, "viewallcv.html", data)
			return
		}

		data.CVInfo.CV = cv
		gob.Register(cv)
		session.Values["ApplicantFabricID"] = applicantFabricID
		session.Values["CV"] = cv

		err = session.Save(r, w)
		if err != nil {
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "viewallcv.html", data)
			return
		}

		renderTemplate(w, r, "reviewcv.html", data)
	})
}



func (c *Controller) ReviewCVHandler() func(http.ResponseWriter, *http.Request) {
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

		// Check that the user connected is an admin
		_, err := u.QueryVerifier()
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

		rating := model.CVReview{
			Name:    r.FormValue("name"),
			Comment: r.FormValue("comment"),
			Rating:  ratingInt,
		}

		applicantID := sessions.GetApplicantFabricID(session)
		cvHash := sessions.GetCVHash(session)

		if applicantID == "" || cvHash == "" {
			data.MessageWarning = "Error! Unable to retrieve cv information"
			renderTemplate(w, r, "index.html", data)
			return
		}

		reviewByte, err := json.Marshal(rating)
		if err != nil {
			data.MessageWarning = "Error! Unable to save review in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		err = u.UpdateSaveRating(applicantID, cvHash, reviewByte)
		if err != nil {
			data.MessageWarning = "Error! Unable to save rating in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = sessions.GetCV(session)

		data.CurrentPage = "viewallcv"
		data.MessageSuccess = "Success! Your review has been saved."

		renderTemplate(w, r, "index.html", data)
	})
}
