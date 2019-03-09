package controllers

import (
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"net/http"
)

func (app *Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {

	data := models.TemplateData{
		CurrentPage:  "index",
		LoggedInFlag: false,
	}

	session := sessions.InitSession(r)
	if sessions.IsLoggedIn(session) {
		data.UserDetails = sessions.GetUserDetails(session)
		data.LoggedInFlag = true
		renderTemplate(w, r, "index.html", data)
	} else {
		renderTemplate(w, r, "index.html", data)
	}
}

