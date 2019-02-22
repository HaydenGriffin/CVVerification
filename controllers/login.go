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

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"login",
		LoggedInFlag:false,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
		data.LoggedInFlag = true
		data.CurrentPage = "index"
		data.MessageWarning = "You are already logged in!"
		renderTemplate(w, r, "index.html", data)
	} else {
		renderTemplate(w, r, "login.html", data)
	}
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	session := sessions.InitSession(r)

	var user models.User
	var passwordCorrect = false

	user, err := database.GetUserFromUsername(username)

	data := models.TemplateData{
		CurrentUser:  models.User{},
		CurrentPage:  "login",
		LoggedInFlag: false,
	}
	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Error! Incorrect username or password."
	} else {
		passwordCorrect = crypto.Compare(user.Password, password)
	}

	if passwordCorrect {
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
		data.MessageSuccess = "You have successfully logged in! Welcome, " + user.FullName
		renderTemplate(w, r, "index.html", data)
	} else {
		data.CurrentUser.Username = username
		data.MessageWarning = "Error! Incorrect username or password."
		renderTemplate(w, r, "login.html", data)
	}
}

// Logout
func (app *Application) LogoutHandler(w http.ResponseWriter, r *http.Request)  {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"login",
		LoggedInFlag:false,
	}

	session.Values["User"] = models.User{}
	session.Values["LoggedInFlag"] = false
	session.Save(r, w)

	data.MessageSuccess = "You have been successfully logged out."
	renderTemplate(w, r, "login.html", data)
}