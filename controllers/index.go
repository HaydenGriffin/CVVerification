package controllers

import (
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"net/http"

)



func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {

	data := &struct {
		CurrentUser models.User
		LoggedInFlag bool
	}{
		CurrentUser:models.User{},
		LoggedInFlag:false,
	}

	session := sessions.InitSession(r)
	if sessions.IsLoggedIn(session) {
		user := sessions.GetUser(session)
		data.CurrentUser = user
		data.LoggedInFlag = true
		renderTemplate(w, r, "index.html", data)
	} else {
		renderTemplate(w, r, "index.html", nil)
	}
}

