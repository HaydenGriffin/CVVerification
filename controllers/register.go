package controllers

import (
	"encoding/gob"
	"fmt"
	"github.com/cvtracker/crypto"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/service"
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

	data := models.TemplateData{
		//CurrentUser:  user,
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

	var user models.User

	user.Username = r.FormValue("username")
	user.FullName = r.FormValue("fullName")
	user.EmailAddress = r.FormValue("emailAddress")
	hashedPassword, err := crypto.GenerateFromString(r.FormValue("password"))

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Error! Something went wrong. Please try again"
		renderTemplate(w, r, "register.html", data)
		return
	}
	user.Password = hashedPassword
	user.UserRole = "APPLICANT"
	profileHash, err := crypto.GenerateFromString(user.Username)
	user.ProfileHash = profileHash

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Error! Something went wrong. Please try again"
		renderTemplate(w, r, "register.html", data)
		return
	}

	err = database.CreateNewUser(user.Username, user.FullName, user.Password, user.EmailAddress, user.UserRole, profileHash)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Failed to create new account."
		renderTemplate(w, r, "register.html", data)
		return
	}

	profile := service.UserProfile{
		Username:user.Username,
	}
	//user, err = database.GetUserFromUsername(username)
	txid, err := app.Service.SaveProfile(profile, user.ProfileHash)

	gob.Register(user)
	session.Values["User"] = user
	session.Values["LoggedInFlag"] = true
	err = session.Save(r, w)
	data.CurrentUser = user
	data.CurrentPage = "index"
	data.LoggedInFlag = true
	data.MessageSuccess = "You have successfully created an account! Welcome, " + user.FullName + txid
	renderTemplate(w, r, "index.html", data)
}