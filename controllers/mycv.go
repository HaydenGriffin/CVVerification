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
			for i := range allCVHistoryInfo {
				fmt.Println("first info")
				fmt.Println(allCVHistoryInfo[i])

				if allCVHistoryInfo[i].Index == cvToDisplayID {
					cvToDisplayCVHash = allCVHistoryInfo[i].CVHash
					if allCVHistoryInfo[i].CVInReview == 1 {
						fmt.Println("cv is in review")
						foundCVInReview = true
						data.CurrentCVInReview = true
					}
				}
			}
		} else {
			for i := range allCVHistoryInfo {
				fmt.Println("first info")
				fmt.Println(allCVHistoryInfo[i])

				if allCVHistoryInfo[i].CVInReview == 1 {
					cvToDisplayCVHash = allCVHistoryInfo[i].CVHash
					fmt.Println("cv is in review")
					foundCVInReview = true
					data.CurrentCVInReview = true
				}
			}
		}

		// No CV in review - display the most recent CV version to the user
		if cvToDisplayID == 0 && foundCVInReview == false {
			cvToDisplayCVHash = allCVHistoryInfo[len(allCVHistoryInfo)-1].CVHash
		}

		cv, err := u.QueryCV(cvToDisplayCVHash)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Unable to retrieve CV detail from ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}
		data.CV = cv

		data.UserHasCVInReview = database.UserHasCVInReview(data.UserDetails.Id)


		data.CVHistory = allCVHistoryInfo
		data.CurrentPage = "mycv"

		gob.Register(data.CV)
		session.Values["CV"] = data.CV
		session.Values["CVHash"] = cvToDisplayCVHash

		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err.Error())
			renderTemplate(w, r, "index.html", data)
			return
		}

		renderTemplate(w, r, "mycv.html", data)

		//_, cvHash, err := database.GetCVInfoFromID(data.UserDetails.Id)

		/*	if err != nil {
				fmt.Printf(err.Error())
				data.MessageWarning = "Unable to find CV info in database."
				renderTemplate(w, r, "index.html", data)
				return
			}

			cv, err := u.QueryCV(cvHash)

			if err != nil {
				fmt.Printf(err.Error())
				data.MessageWarning = "Unable to retrieve CV details from ledger."
				renderTemplate(w, r, "index.html", data)
				return
			}*/

		//Retrieve ratings
		//b, err = app.Service.GetRatings(data.CurrentUser.ProfileHash, cvHash)

		//fmt.Println(b)

		// No ratings exist yet
		/*if err != nil {
			fmt.Printf(err.Error())
		} else {
			var ratings []service.CVRating
			err = json.Unmarshal(b, &ratings)
			if err != nil {
				fmt.Printf(err.Error())
				data.MessageWarning = "Unable to retrieve ratings from ledger"
				renderTemplate(w, r, "index.html", data)
			} else {
				session.Values["Ratings"] = ratings
				session.Save(r, w)
				data.Ratings = ratings
				fmt.Println("CV Ratings:")
				fmt.Println(data.Ratings)
			}
		}*/

		/*isRatable, err := database.IsCVInReview(cvHash)

		if err != nil {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to get status of CV."
			renderTemplate(w, r, "index.html", data)
			return
		}
		data.IsCVInReview = isRatable
		data.CV = cv
		gob.Register(cv)
		renderTemplate(w, r, "mycv.html", data)*/
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
		if cv == nil || cvHash == "" {
			data.MessageWarning = "Unable to update status of CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CV = cv

		data.UserHasCVInReview = database.UserHasCVInReview(data.UserDetails.Id)

		if data.UserHasCVInReview {
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
			data.MessageSuccess = "Success! Your CV can now be reviwed."
			data.CurrentCVInReview = true
			data.CurrentPage = "mycv"
			data.UserHasCVInReview = true
		}

		allCVHistoryInfo, err := database.GetUserCVDetails(data.UserDetails.Id)

		if err != nil || len(allCVHistoryInfo) == 0 {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVHistory = allCVHistoryInfo
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
		if cv == nil || cvHash == "" {
			data.MessageWarning = "Unable to update status of CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CV = cv

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

		data.CVHistory = allCVHistoryInfo
		data.MessageSuccess = "Success! Your CV has been withdrawn from review."
		data.CurrentCVInReview = false
		//data.Ratings = sessions.GetRatings(session)
		renderTemplate(w, r, "mycv.html", data)

		//cvHash, err := app.Service.GetCVHashFromProfile(data.CurrentUser.ProfileHash)
		/*
		if err != nil {
			fmt.Println("Error GetCVHashFromProfile: " + err.Error())
		}


		err = database.UpdateCV(cvHash,0)
		if err != nil {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to update database."
			renderTemplate(w, r, "mycv.html", data)
		} else {
			data.MessageSuccess = "Success! Your CV can now be edited."
			data.IsCVInReview = false
			data.CV = sessions.GetCV(session)
			data.Ratings = sessions.GetRatings(session)*/
	})
}
