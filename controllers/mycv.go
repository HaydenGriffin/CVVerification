package controllers

import (
	"encoding/gob"
	"fmt"
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (c *Controller) MyCVHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "index",
		}

		if sessions.IsLoggedIn(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.MessageWarning = "You must be logged in to view your CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			data.MessageWarning = "You must be an applicant user to upload a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		fmt.Println("Query profile")
		fmt.Println(u.QueryProfile(data.UserDetails.ProfileHash))

		result, success := mux.Vars(r)["cvToDisplayID"]
		var cvToDisplayID int

		if !success {
			fmt.Println("No passed CV ID")
		} else {
			cvToDisplayID, err = strconv.Atoi(result)
			if err != nil {
				fmt.Println("Unable to convert CV to display ID to number")
			} else {
				fmt.Println(cvToDisplayID)
			}
		}

		allCVHistoryInfo, err := database.GetUserCVDetails(data.UserDetails.Id)

		if err != nil || len(allCVHistoryInfo) == 0 {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		var foundCVInReview = false
		var cvToDisplayCVHash string

		if cvToDisplayID != 0 {
			// Specific CV ID requested
			for i := range allCVHistoryInfo {
				fmt.Println("first info")
				fmt.Println(allCVHistoryInfo[i])

				if allCVHistoryInfo[i].Index == cvToDisplayID {
					cvToDisplayCVHash = allCVHistoryInfo[i].CVHash
					if allCVHistoryInfo[i].CVInReview == 1 {
						fmt.Println("cv is in review")
						foundCVInReview = true
						data.CVInfo.CurrentCVInReview = true
					}
					data.CVInfo.CurrentCVHash = allCVHistoryInfo[i].CVHash
				}
			}
		} else {
			// No CV ID requested - retrieve any CV in review
			for i := range allCVHistoryInfo {
				fmt.Println("first info")
				fmt.Println(allCVHistoryInfo[i])

				if allCVHistoryInfo[i].CVInReview == 1 {
					cvToDisplayCVHash = allCVHistoryInfo[i].CVHash
					fmt.Println("cv is in review")
					foundCVInReview = true
					data.CVInfo.CurrentCVInReview = true
					data.CVInfo.CurrentCVHash = allCVHistoryInfo[i].CVHash
				}
			}
		}

		// No CV in review - display the most recent CV version to the user
		if cvToDisplayID == 0 && foundCVInReview == false {
			cvToDisplayCVHash = allCVHistoryInfo[len(allCVHistoryInfo)-1].CVHash
			data.CVInfo.CurrentCVHash = allCVHistoryInfo[len(allCVHistoryInfo)-1].CVHash
		}

		// Retrieve CV from ledger
		cv, err := u.QueryCV(cvToDisplayCVHash)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Unable to retrieve CV detail from ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}
		data.CVInfo.CV = cv

		data.CVInfo.UserHasCVInReview = database.UserHasCVInReview(data.UserDetails.Id)


		data.CVInfo.CVHistory = allCVHistoryInfo
		data.CurrentPage = "mycv"


		reviews, err := u.QueryCVReviews(data.UserDetails.ProfileHash, cvToDisplayCVHash)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "An error occurred whilst retrieving ratings for the CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.Reviews = reviews
		fmt.Println(reviews)

		gob.Register(cv)
		gob.Register(reviews)
		session.Values["CV"] = cv
		session.Values["CVHash"] = cvToDisplayCVHash
		session.Values["Reviews"] = reviews

		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err.Error())
			renderTemplate(w, r, "index.html", data)
			return
		}

		renderTemplate(w, r, "mycv.html", data)
	})
}

func (c *Controller) SubmitForReviewHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "index",
		}

		if sessions.IsLoggedIn(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.MessageWarning = "You must be logged in to view your CV."
			renderTemplate(w, r, "index.html", data)
			return
		}


		cv := sessions.GetCV(session)
		cvHash := sessions.GetCVHash(session)
		reviews := sessions.GetReviews(session)
		if cv == nil || cvHash == "" {
			data.MessageWarning = "Unable to update status of CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = cv
		data.CVInfo.CurrentCVHash = cvHash
		data.CVInfo.Reviews = reviews

		data.CVInfo.UserHasCVInReview = database.UserHasCVInReview(data.UserDetails.Id)

		if data.CVInfo.UserHasCVInReview {
			fmt.Println("Only allowed one CV in review")
			data.MessageWarning = "Error! You are only allowed one version of your CV in review at a time."
		} else {
			err := database.UpdateCV(cvHash, 1)
			if err != nil {
				fmt.Printf(err.Error())
				data.MessageWarning = "Unable to update CV info in database."
				renderTemplate(w, r, "index.html", data)
				return
			}
			data.MessageSuccess = "Success! Your CV can now be reviewed."
			data.CVInfo.CurrentCVInReview = true
			data.CurrentPage = "mycv"
			data.CVInfo.UserHasCVInReview = true
		}

		allCVHistoryInfo, err := database.GetUserCVDetails(data.UserDetails.Id)

		if err != nil || len(allCVHistoryInfo) == 0 {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CVHistory = allCVHistoryInfo
		//data.Ratings = sessions.GetRatings(session)
		renderTemplate(w, r, "mycv.html", data)

	})
}


func (c *Controller) WithdrawFromReviewHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "index",
		}

		if sessions.IsLoggedIn(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.MessageWarning = "You must be logged in to view your CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cv := sessions.GetCV(session)
		cvHash := sessions.GetCVHash(session)
		reviews := sessions.GetReviews(session)
		if cv == nil || cvHash == "" {
			data.MessageWarning = "Unable to update status of CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CV = cv
		data.CVInfo.CurrentCVHash = cvHash
		data.CVInfo.Reviews = reviews

		err := database.UpdateCV(cvHash, 0)
		if err != nil {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to update CV info in database."
			renderTemplate(w, r, "index.html", data)
		}

		allCVHistoryInfo, err := database.GetUserCVDetails(data.UserDetails.Id)

		if err != nil || len(allCVHistoryInfo) == 0 {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CVHistory = allCVHistoryInfo
		data.MessageSuccess = "Success! Your CV has been withdrawn from review."
		data.CVInfo.CurrentCVInReview = false
		//data.Ratings = sessions.GetRatings(session)
		renderTemplate(w, r, "mycv.html", data)

	})
}
