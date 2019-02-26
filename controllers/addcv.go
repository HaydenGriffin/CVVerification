package controllers

import (
	"encoding/json"
	"fmt"
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
		renderTemplate(w, r, "addcv.html", data)
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

	txid, err = app.Service.UpdateProfile(data.CurrentUser.ProfileHash, cvHash)


	db, err := database.InitDB("root:password@tcp(localhost:3306)/verification")

	res, err := db.Exec("INSERT INTO user_cvs(user_id, timestamp, cv) VALUES (?, CURRENT_TIMESTAMP, ?)", data.CurrentUser.Id, cv.CV)


	if err != nil {
		panic(err)
	} else {
		fmt.Println(res)
	}

	if err != nil {
		data.MessageWarning = "Unable to invoke"
		renderTemplate(w, r, "mycv.html", data)
	} else {
		data.MessageSuccess = txid
		renderTemplate(w, r, "index.html", data)
	}
}
