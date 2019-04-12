package controllers

import (
	"encoding/gob"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (c *Controller) UpdateDetailsView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		data := templateModel.Data{
			CurrentPage:  "userdetails",
		}

		session := sessions.InitSession(r)

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
		}
		renderTemplate(w, r, "userdetails.html", data)
	})
}
func (c *Controller) UpdateDetailsHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		data := templateModel.Data{
			CurrentPage:  "index",
		}

		session := sessions.InitSession(r)


		// Form values to insert into DB
		username := r.FormValue("username")
		fullName := r.FormValue("fullName")
		emailAddress := r.FormValue("emailAddress")

		// Retrieve user details
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
				data.UserDetails = userDetails
				data.MessageSuccess = "Success! You details have been updated."
				renderTemplate(w, r, "index.html", data)
				return
			}
		} else {
			// Profile details haven't been saved to DB yet
			fabricID, err := u.QueryID()
			if err != nil {
				data.MessageWarning = "Error! Unable to retrieve profile ID from ledger."
			}

			// Insert row into DB
			userDetails, err := database.CreateNewUser(username, fullName, emailAddress, fabricID)
			if err != nil {
				data.MessageWarning = "Error! Unable to save user details."
				renderTemplate(w, r, "userdetails.html", data)
				return
			}
			// Register the userDetails gob to be used as a session value
			gob.Register(userDetails)
			session.Values["UserDetails"] = userDetails
			data.UserDetails = userDetails
			data.MessageSuccess = "Success! Your details have been saved. Welcome, " + userDetails.FullName
			renderTemplate(w, r, "index.html", data)
		}

		err := session.Save(r, w)
		if err != nil {
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

	})
}
