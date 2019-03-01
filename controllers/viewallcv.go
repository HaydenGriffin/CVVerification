package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/service"
	"github.com/cvtracker/sessions"
	"net/http"
)

func (app *Application) ViewAllView(w http.ResponseWriter, r *http.Request) {

	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentUser:  models.User{},
		CurrentPage:  "index",
		LoggedInFlag: false,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
		data.LoggedInFlag = true
	} else {
		data.MessageWarning = "You must be logged in to view the CVs."
		renderTemplate(w, r, "index.html", data)
		return
	}

	ratableCVs := make(map[int] string)

	ratableCVs, err := database.GetAllRatableCVHashes()
	fmt.Println(ratableCVs)

	if err != nil {
		data.MessageWarning = err.Error()
		renderTemplate(w, r, "index.html", data)
		return
	}

	data.CVList = make(map[int] service.CVObject)


	for userID, cvHash := range ratableCVs {
		fmt.Println("profileHash: " + string(userID))
		fmt.Println("cvHash: " + cvHash)
		b, err := app.Service.GetCVFromCVHash(cvHash)

		if err != nil {
			data.MessageWarning = err.Error()
			renderTemplate(w, r, "index.html", data)
			return
		}
		var cv service.CVObject
		err = json.Unmarshal(b, &cv)
		if err != nil {
			data.MessageWarning = err.Error()
			renderTemplate(w, r, "index.html", data)
			return
		}
		data.CVList[userID] = cv
	}

		if len(data.CVList) == 0 {
			data.MessageWarning = "There are no CVs to be rated at this time."
			renderTemplate(w, r, "viewall.html", data)
			return
		}

		renderTemplate(w, r, "viewall.html", data)
}
