package controllers

import (
	"encoding/gob"
	"fmt"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (app *Controller) LoginView(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentPage:  "login",
		LoggedInFlag: false,
	}

	if sessions.IsLoggedIn(session) {
		data.UserDetails = sessions.GetUserDetails(session)
		data.LoggedInFlag = true
		data.CurrentPage = "index"
		data.MessageWarning = "You are already logged in!"
		renderTemplate(w, r, "index.html", data)
	} else {
		renderTemplate(w, r, "login.html", data)
	}
}

func (app *Controller) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentPage:  "login",
		LoggedInFlag: false,
	}

	if sessions.IsLoggedIn(session) {
		data.UserDetails = sessions.GetUserDetails(session)
		data.LoggedInFlag = true
		data.CurrentPage = "index"
		data.MessageWarning = "You are already logged in!"
		renderTemplate(w, r, "index.html", data)
		return
	}

	fabricUser, err := app.Fabric.LogUser(username, password)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Error! Incorrect username or password. Please try again"
		renderTemplate(w, r, "login.html", data)
		return
	}

	// Successfully logged in, retrieve user details from DB
	userDetails, err := database.GetUserDetailsFromUsername(fabricUser.Username)
	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Failed to retrieve user profile details from DB."
	} else {
		data.UserDetails = userDetails
	}

	gob.Register(userDetails)
	session.Values["UserDetails"] = userDetails
	session.Values["LoggedInFlag"] = true
	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err.Error())
	}
	data.CurrentPage = "index"
	data.LoggedInFlag = true
	data.MessageSuccess = "You have successfully logged in! Welcome, " + fabricUser.Username
	renderTemplate(w, r, "index.html", data)
}

// Logout
func (app *Controller) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentPage:  "login",
		LoggedInFlag: false,
	}

	session.Values["UserDetails"] = models.UserDetails{}
	session.Values["LoggedInFlag"] = false
	err := session.Save(r, w)
	if err != nil {
		fmt.Println(err.Error())
	}

	data.MessageSuccess = "You have been successfully logged out."
	renderTemplate(w, r, "login.html", data)
}
