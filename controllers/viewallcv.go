package controllers

import (
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"net/http"
)

func (app *Application) ViewAllHandler(w http.ResponseWriter, r *http.Request) {

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
	}

/*	var cvHashList []string

	cvHashList, err := database.GetAllCVHashes()

	for

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Unable to find CV from hash."
	} else {
		cvString, err := app.Service.QueryCVByHash(cvHash)
		if err != nil {
			data.MessageWarning = "Unable to query the blockchain"
		} else {
			var cv = service.CVObject{}
			err := json.Unmarshal(cvString, &cv)
			if err != nil {
				data.MessageWarning = "Unable to unmarshal CV object"
			}
			data.CV = cv
		}
		renderTemplate(w, r, "mycv.html", data)
	}*/
}
