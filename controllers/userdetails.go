package controllers

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	"github.com/cvverification/crypto"
	"github.com/cvverification/database"
	"github.com/cvverification/models"
	"github.com/cvverification/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (c *Controller) UpdateDetailsView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		data := models.TemplateData{
			CurrentPage:  "userdetails",
		}

		session := sessions.InitSession(r)
		if sessions.HasSavedUserDetails(session) {
			userDetails := sessions.GetUserDetails(session)
			data.UserDetails = userDetails
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
		}
		renderTemplate(w, r, "userdetails.html", data)
	})
}
func (c *Controller) UpdateDetailsHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		data := models.TemplateData{
			CurrentPage:  "index",
		}

		session := sessions.InitSession(r)


		// Form values to insert into DB
		username := r.FormValue("username")
		fullName := r.FormValue("fullName")
		emailAddress := r.FormValue("emailAddress")

		if sessions.HasSavedUserDetails(session) {
			// Logic to update a profile
			userDetails, err := database.UpdateUser(username, fullName, emailAddress)
			if err != nil {
				fmt.Println(err)
				data.MessageWarning = "An error occurred whilst trying to update your profile information."
				data.CurrentPage = "userdetails"
				renderTemplate(w, r, "userdetails.html", data)
				return
			} else {
				// Successfully updated user
				// Update the session values and save session
				gob.Register(userDetails)
				session.Values["UserDetails"] = userDetails
				err = session.Save(r, w)
				if err != nil {
					fmt.Printf(err.Error())
				}
				data.UserDetails = userDetails
				data.MessageSuccess = "You have successfully updated your details!"
				renderTemplate(w, r, "index.html", data)
				return
			}
		} else {
			fmt.Printf("Saving profile for first time")
			// Profile details haven't been saved to DB yet
			// Generate a unique profile hash. This is made up from the users signing identity + their chosen username
			profileHash, err := crypto.GenerateFromString(u.SigningIdentity.Identifier().ID + u.Username)
			if err != nil {
				fmt.Printf(err.Error())
				data.MessageWarning = "Error! Something went wrong. Please try again"
				renderTemplate(w, r, "userdetails.html", data)
				return
			}

			profile := model.UserProfile{
				Username:u.Username,
			}

			profileByte, err := json.Marshal(profile)

			err = u.UpdateSaveProfile(profileByte, profileHash)

			if err != nil {
				fmt.Println("Unable to save profile to ledger")
			}

			// Insert row into DB
			userDetails, err := database.CreateNewUser(username, fullName, emailAddress, profileHash)
			if err != nil {
				fmt.Printf(err.Error())
				data.MessageWarning = "Error! Something went wrong whilst saving the user details. Please try again"
				renderTemplate(w, r, "userdetails.html", data)
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
	})
}
