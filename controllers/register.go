package controllers

import (
	"encoding/gob"
	"fmt"
	"github.com/cvtracker/crypto"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (app *Application) RegisterView(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"register",
		LoggedInFlag:false,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
		data.LoggedInFlag = true
		data.CurrentPage = "index"
		data.MessageWarning = "You are already logged in!"
		renderTemplate(w, r, "index.html", data)
	} else {
		renderTemplate(w, r, "register.html", data)
	}
}

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	username := r.FormValue("username")
	fullName := r.FormValue("fullName")
	emailAddress := r.FormValue("emailAddress")
	password := r.FormValue("password")

	data := models.TemplateData{
		CurrentUser:  models.User{},
		CurrentPage:  "register",
		LoggedInFlag: false,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
		data.LoggedInFlag = true
		data.CurrentPage = "index"
		data.MessageWarning = "You are already logged in!"
		renderTemplate(w, r, "index.html", data)
		return
	}

	hashedPassword, err := crypto.Generate(password)

	if err != nil {
		fmt.Printf(err.Error())
	}

	err = database.CreateNewUser(username,fullName,hashedPassword,emailAddress)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Failed to create new account."
		renderTemplate(w, r, "register.html", data)
		return
	}

	var user models.User

	user, err = database.GetUserFromUsername(username)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Failed to retrieve account from database."
		renderTemplate(w, r, "register.html", data)
		return
	} else {
		gob.Register(user)
		session.Values["User"] = user
		session.Values["LoggedInFlag"] = true
		err := session.Save(r, w)
		if err != nil {
			fmt.Println(err.Error())
		}
		data.CurrentUser = user
		data.CurrentPage = "index"
		data.LoggedInFlag = true
		data.MessageSuccess = "You have successfully created an account! Welcome, " + user.FullName
		renderTemplate(w, r, "index.html", data)
	}
}