package controllers

import (
	"encoding/gob"
	"fmt"
	"github.com/cvverification/blockchain"
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
				data.MessageWarning = "Error! Unable to update profile information in database."
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
				data.MessageSuccess = "Success! You details have been updated."
				renderTemplate(w, r, "index.html", data)
				return
			}
		} else {
			// Profile details haven't been saved to DB yet
			userID, err := u.QueryID()
			if err != nil {
				data.MessageWarning = "Error! Unable to retrieve profile ID from ledger."
			}

			// Insert row into DB
			userDetails, err := database.CreateNewUser(username, fullName, emailAddress, userID)
			if err != nil {
				data.MessageWarning = "Error! Unable to save user details."
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
			data.MessageSuccess = "Success! Your details have been saved. Welcome, " + userDetails.FullName
			renderTemplate(w, r, "index.html", data)
		}
	})
}
