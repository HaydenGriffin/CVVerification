package controllers

import (
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"net/http"
)

func (app *Controller) ResultHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		//CurrentUser:  models.User{},
		CurrentPage:  "index",
		LoggedInFlag: false,
	}

	if sessions.IsLoggedIn(session) {
		data.UserDetails = sessions.GetUserDetails(session)
		data.LoggedInFlag = true
	} else {
		data.MessageWarning = "You must be logged in to view your CV."
		renderTemplate(w, r, "index.html", data)
		return
	}

	/*	_, cvHash, err := database.GetCVInfoFromID(data.CurrentUser.Id)

		if err != nil {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}*/

	/*	b, err := app.Service.GetCVFromCVHash(cvHash)

		if err != nil {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to find CV from CVHash."
			renderTemplate(w, r, "index.html", data)
			return
		}

		var cv service.CVObject

		err = json.Unmarshal(b, &cv)

		if err != nil {
			fmt.Printf(err.Error())
			data.MessageWarning = "An error occurred whilst retrieving the CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		//Retrieve ratings
		b, err = app.Service.GetRatings(data.CurrentUser.ProfileHash, cvHash)

		fmt.Println(b)

		// No ratings exist yet
		if err != nil {
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
		}

		isRatable, err := database.IsCVRatable(cvHash)

		if err != nil {
			fmt.Printf(err.Error())
			data.MessageWarning = "Unable to get status of CV."
			renderTemplate(w, r, "index.html", data)
			return
		}
		data.IsCVRatable = isRatable
		data.CV = cv
		gob.Register(cv)
		session.Values["CV"] = cv
		session.Save(r, w)*/
		renderTemplate(w, r, "mycv.html", data)
}

func (app *Controller) SubmitForReviewHandler(w http.ResponseWriter, r *http.Request) {

	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentPage:  "index",
		LoggedInFlag: false,
	}

	if sessions.IsLoggedIn(session) {
		data.UserDetails = sessions.GetUserDetails(session)
		data.LoggedInFlag = true
	} else {
		data.MessageWarning = "You must be logged in to view your CV."
		renderTemplate(w, r, "index.html", data)
		return
	}

/*	cvHash, err := app.Service.GetCVHashFromProfile(data.CurrentUser.ProfileHash)

	if err != nil {
		fmt.Println("Error GetCVHashFromProfile: " + err.Error())
	}*/


	/*err = database.UpdateCV(cvHash,1)
	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Unable to update database."
		renderTemplate(w, r, "mycv.html", data)
	} else {
		data.MessageSuccess = "Success! Your CV can now be rated."
		data.IsCVRatable = true
		data.CV = sessions.GetCV(session)
		data.Ratings = sessions.GetRatings(session)*/
		renderTemplate(w, r, "mycv.html", data)
}

func (app *Controller) WithdrawFromReviewHandler(w http.ResponseWriter, r *http.Request) {

	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentPage:  "index",
		LoggedInFlag: false,
	}

	if sessions.IsLoggedIn(session) {
		data.UserDetails = sessions.GetUserDetails(session)
		data.LoggedInFlag = true
	} else {
		data.MessageWarning = "You must be logged in to view your CV."
		renderTemplate(w, r, "index.html", data)
		return
	}

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
		data.IsCVRatable = false
		data.CV = sessions.GetCV(session)
		data.Ratings = sessions.GetRatings(session)*/
		renderTemplate(w, r, "mycv.html", data)
}