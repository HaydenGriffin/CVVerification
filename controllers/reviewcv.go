package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	"github.com/cvverification/database"
	"github.com/cvverification/models"
	"github.com/cvverification/sessions"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (c *Controller) ReviewCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		fmt.Println("ReviewCVView")

		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "addcv",
		}

		if sessions.IsLoggedIn(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.MessageWarning = "You must be logged in to view the CVs."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Check that the user connected is an admin
		_, err := u.QueryVerifier()
		if err != nil {
			fmt.Println(err)
			data.CurrentPage = "index"
			data.MessageWarning = "You must be a verifier user to rate a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		result, success := mux.Vars(r)["userID"]

		if !success {
			data.MessageWarning = "Error! No CV to be retrieved."
			renderTemplate(w, r, "index.html", data)
			return
		}

		userID, err := strconv.Atoi(result)

		if err != nil {
			data.MessageWarning = "Error! Invalid CV ID."
			renderTemplate(w, r, "index.html", data)
			return
		}

		profileHash, cvHash, err := database.GetCVInfoFromID(userID)

		if err != nil {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		verifierReview, err := u.QueryCVReviewable(profileHash, cvHash)
		if err != nil {
			fmt.Printf(err.Error())
			fmt.Printf("cuck")

		}

		if (model.CVReview{}) == verifierReview {
			fmt.Printf("Verifier hasn't reviewed yet")
		} else {
			fmt.Printf("Verifier has reviewed!!!")
		}
		fmt.Println(verifierReview)

		data.CVInfo.Review = verifierReview

		cv, err := u.QueryCV(cvHash)

		if err != nil {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to find CV from hash."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = cv

		session.Values["ProfileHash"] = profileHash
		session.Values["CVHash"] = cvHash

		err = session.Save(r, w)
		if err != nil {
			data.MessageWarning = err.Error()
			fmt.Println(err.Error())
			renderTemplate(w, r, "index.html", data)
			return
		}

		renderTemplate(w, r, "reviewcv.html", data)
	})
}



func (c *Controller) ReviewCVHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		fmt.Println("ReviewCVHandler")

		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "addcv",
		}

		if sessions.IsLoggedIn(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.MessageWarning = "You must be logged in to view the CVs."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Check that the user connected is an admin
		_, err := u.QueryVerifier()
		if err != nil {
			fmt.Println(err)
			data.CurrentPage = "index"
			data.MessageWarning = "You must be a verifier user to review a CV."
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

		profileHash := sessions.GetProfileHash(session)
		cvHash := sessions.GetCVHash(session)

		if profileHash == "" || cvHash == "" {
			data.MessageWarning = "Couldn't retrieve cv information"
			renderTemplate(w, r, "index.html", data)
			return

		}

		ratingByte, err := json.Marshal(rating)

		err = u.UpdateSaveRating(profileHash, cvHash, ratingByte)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "An error occurred whilst saving rating in ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		//data.MessageSuccess = txid
		renderTemplate(w, r, "reviewcv.html", data)
	})
}
