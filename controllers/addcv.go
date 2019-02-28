package controllers

import (
	"encoding/json"
	"github.com/cvtracker/crypto"
	"github.com/cvtracker/database"
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
		renderTemplate(w, r, "cvform.html", data)
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to add a CV."
		renderTemplate(w, r, "index.html", nil)
	}
}

func (app *Application) AddCVHandler(w http.ResponseWriter, r *http.Request) {

	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"addcv",
		LoggedInFlag:true,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to add a CV."
		renderTemplate(w, r, "index.html", nil)
		return
	}


	cv := service.CVObject{
		Name:r.FormValue("name"),
		Speciality:r.FormValue("speciality"),
		CV:r.FormValue("CV"),
		CVDate:r.FormValue("CVDate"),
	}

	cvByte, err := json.Marshal(cv)

	cvHash, err := crypto.GenerateFromByte(cvByte)

	txid, err := app.Service.SaveCV(cv, cvHash)

	txid, err = app.Service.UpdateProfileCV(data.CurrentUser.ProfileHash, cvHash)

	if err != nil {
		data.MessageWarning = err.Error()
		renderTemplate(w, r, "index.html", data)
		return
	}

	err = database.CreateNewCV(data.CurrentUser.Id, cv.CV, cvHash)

	if err != nil {
		data.MessageWarning = "Unable to create new CV"
		renderTemplate(w, r, "mycv.html", data)
	} else {
		data.MessageSuccess = txid
		renderTemplate(w, r, "index.html", data)
	}
}
