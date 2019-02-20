package controllers

import (
	"github.com/cvtracker/models"
	"github.com/cvtracker/service"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)


func (app *Application) AddCVView(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"addcv",
		LoggedInFlag:true,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
		renderTemplate(w, r, "addcv.html", data)
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to add a CV."
		renderTemplate(w, r, "index.html", nil)
	}
}

func (app *Application) AddCVHandler(w http.ResponseWriter, r *http.Request) {

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"addcv",
		LoggedInFlag:true,
	}

	cv := service.CVObject{
		Name:r.FormValue("name"),
		Speciality:r.FormValue("speciality"),
		CVHash:r.FormValue("CV"),
		CVDate:r.FormValue("CVDate"),
	}

	txid, err := app.Service.SaveCV(cv)

	if err != nil {
		data.MessageWarning = "Unable to invoke"
		renderTemplate(w, r, "mycv.html", data)
	} else {
		data.MessageSuccess = txid
		renderTemplate(w, r, "mycv.html", data)
	}
}
