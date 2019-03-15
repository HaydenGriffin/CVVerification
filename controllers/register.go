package controllers

import (
	"encoding/gob"
	"fmt"
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/crypto"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (c *Controller) RegisterView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {
		data := models.TemplateData{
			CurrentPage:  "register",
			LoggedInFlag: false,
		}

		session := sessions.InitSession(r)
		if sessions.IsLoggedIn(session) {
			data.LoggedInFlag = true
			if sessions.HasSavedUserDetails(session) {
				userDetails := sessions.GetUserDetails(session)
				data.UserDetails = userDetails
				renderTemplate(w, r, "register.html", data)
			} else {
				data.CurrentPage = "register"
				data.UserDetails.Username = u.Username
				renderTemplate(w, r, "register.html", data)
			}
		}
	})
}
func (c *Controller) RegisterHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		data := models.TemplateData{
			CurrentPage:  "index",
			LoggedInFlag: false,
		}

		session := sessions.InitSession(r)
		if sessions.IsLoggedIn(session) {
			data.LoggedInFlag = true

			// Form values to insert into DB
			username := r.FormValue("username")
			fullName := r.FormValue("fullName")
			emailAddress := r.FormValue("emailAddress")

			if sessions.HasSavedUserDetails(session) {
				// Logic to update a profile
				renderTemplate(w, r, "index.html", data)
			} else {
				// Profile details haven't been saved to DB yet
				// Generate a unique profile hash. This is made up from the users signing identity + their chosen username
				profileHash, err := crypto.GenerateFromString(u.SigningIdentity.Identifier().ID + u.Username)
				if err != nil {
					fmt.Printf(err.Error())
					data.MessageWarning = "Error! Something went wrong. Please try again"
					renderTemplate(w, r, "register.html", data)
					return
				}

				// Insert row into DB
				userDetails, err := database.CreateNewUser(username, fullName, emailAddress, profileHash)
				if err != nil {
					fmt.Printf(err.Error())
					data.MessageWarning = "Error! Something went wrong whilst saving the user details to the DB. Please try again"
					renderTemplate(w, r, "register.html", data)
					return
				}

				// Register the userDetails gob to be used as a session value
				gob.Register(userDetails)
				session.Values["UserDetails"] = userDetails
				err = session.Save(r, w)
				if err != nil {
					fmt.Printf(err.Error())
				}
				data.UserDetails = userDetails
				data.MessageSuccess = "You have successfully registered your details! Welcome, " + userDetails.FullName
				renderTemplate(w, r, "index.html", data)
			}
		}


		/*profile := model.UserProfile{
			Username: user.Username,
		}
	*/

		// STILL NEED TO SAVE THE PROFILE
		//_, err = app.Service.SaveProfile(profile, user.ProfileHash)

	})
}
