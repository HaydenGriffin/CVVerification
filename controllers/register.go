package controllers

import (
	"encoding/gob"
	"fmt"
	"github.com/cvtracker/chaincode/model"
	"github.com/cvtracker/crypto"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (app *Controller) RegisterView(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		UserDetails:  models.UserDetails{},
		CurrentPage:  "register",
		LoggedInFlag: false,
	}

	if sessions.IsLoggedIn(session) {
		data.UserDetails = sessions.GetUserDetails(session)
		data.LoggedInFlag = true
		data.CurrentPage = "index"
		data.MessageWarning = "You are already logged in!"
		renderTemplate(w, r, "index.html", data)
	} else {
		renderTemplate(w, r, "register.html", data)
	}
}

func (app *Controller) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		//CurrentUser:  user,
		CurrentPage:  "register",
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

	// Required for CA registration
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Additional user info
	fullName := r.FormValue("fullName")
	emailAddress := r.FormValue("emailAddress")

	err := app.Fabric.RegisterUser(username, password, model.ActorApplicant)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Error! Something went wrong whilst registering. Please try again"
		renderTemplate(w, r, "register.html", data)
		return
	}

	fabricUser, err := app.Fabric.LogUser(username, password)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Failed to retrieve user details. Please try again"
		renderTemplate(w, r, "register.html", data)
		return
	}

	// Generate a unique profile hash. This is made up from the users signing identity + their chosen username
	profileHash, err := crypto.GenerateFromString(fabricUser.SigningIdentity.Identifier().ID+fabricUser.Username)
	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Error! Something went wrong. Please try again"
		renderTemplate(w, r, "register.html", data)
		return
	}

	userDetails, err := database.CreateNewUser(fabricUser.Username,fullName, emailAddress, profileHash)
	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Error! Something went wrong whilst saving the user details to the DB. Please try again"
		renderTemplate(w, r, "register.html", data)
		return
	}


	/*profile := model.UserProfile{
		Username: user.Username,
	}
*/

// STILL NEED TO SAVE THE PROFILE
	//_, err = app.Service.SaveProfile(profile, user.ProfileHash)

	gob.Register(userDetails)
	session.Values["UserDetails"] = userDetails
	session.Values["LoggedInFlag"] = true
	err = session.Save(r, w)
	if err != nil {
		fmt.Printf(err.Error())
	}
	data.CurrentPage = "index"
	data.LoggedInFlag = true
	data.MessageSuccess = "You have successfully created an account! Welcome, " + fabricUser.Username
	renderTemplate(w, r, "index.html", data)
}
