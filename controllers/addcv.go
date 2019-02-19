package controllers

import (
	"fmt"
	"github.com/cvtracker/models"
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

	helloValue := r.FormValue("speciality")
	fmt.Println(helloValue)
	txid, err := app.Service.InvokeHello(helloValue)
	if err != nil {
		data.MessageWarning = "Unable to invoke"
		renderTemplate(w, r, "result.html", data)
	} else {
		data.MessageSuccess = txid
		renderTemplate(w, r, "result.html", data)
	}
}
