package controllers

import (
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"net/http"

)



func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {

	helloValue, err := app.Service.QueryHello()
	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"index",
		LoggedInFlag:false,
		MessageSuccess:helloValue,
	}

	session := sessions.InitSession(r)
	if sessions.IsLoggedIn(session) {
		user := sessions.GetUser(session)
		data.CurrentUser = user
		data.LoggedInFlag = true
		renderTemplate(w, r, "index.html", data)
	} else {
		renderTemplate(w, r, "index.html", data)
	}
}

