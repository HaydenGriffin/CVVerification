package controllers

import (
	"fmt"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"net/http"
)

func (app *Application) ResultHandler(w http.ResponseWriter, r *http.Request) {

	session := sessions.InitSession(r)

	data := &struct {
		CurrentUser    models.User
		CurrentPage    string
		LoggedInFlag   bool
		MessageWarning string
		MessageSuccess string
	}{
		CurrentUser:  models.User{},
		CurrentPage:  "index",
		LoggedInFlag: false,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
		data.LoggedInFlag = true
	} else {
		data.MessageWarning = "You must be logged in to view your CV."
		renderTemplate(w, r, "index.html", data)
	}

	cvHash, err := database.GetCVHashFromUserID(data.CurrentUser.Id)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Unable to find CV from hash."
	} else {
		cv, err := app.Service.QueryCVByHash(cvHash)
		if err != nil {
			data.MessageWarning = "Unable to query the blockchain"
		} else {
			data.MessageSuccess = "cv:" + cv
		}
		renderTemplate(w, r, "mycv.html", data)
	}
}
