package controllers

import (
	"encoding/gob"
	"fmt"
	"github.com/cvverification/app/crypto"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/blockchain"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (c *Controller) RegisterDetailsView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "userdetails",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
			renderTemplate(w, r, "updatedetails.html", data)
		} else {
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
		}
	})
}

func (c *Controller) RegisterDetailsHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "userdetails",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
			data.MessageWarning = "Error! Unable to register - user already registered."
			renderTemplate(w, r, "updatedetails.html", data)
			return
		}

		// Form values to insert into DB
		username := r.FormValue("username")
		title := r.FormValue("title")
		firstName := r.FormValue("firstName")
		surname := r.FormValue("surname")
		dateOfBirth := r.FormValue("dateOfBirth")
		emailAddress := r.FormValue("emailAddress")

		fabricID, err := u.QueryID()
		if err != nil {
			data.MessageWarning = "Error! Unable to retrieve profile ID from ledger."
		}

		// If the user is an applicant
		applicant, applicantErr := u.QueryApplicant()
		// Generate a new public and private key for the user
		if applicantErr == nil {
			privateKey, publicKey := crypto.GenerateKeyPair(2048)
			privateKeyBytes := crypto.PrivateKeyToBytes(privateKey)
			privateKeyString := string(privateKeyBytes)
			publicKeyBytes := crypto.PublicKeyToBytes(publicKey)
			session.Values["PrivateKey"] = privateKeyString
			data.PrivateKey = string(privateKeyBytes)

			err := u.UpdateSaveProfileKey(string(publicKeyBytes))
			if err != nil {
				fmt.Println(err)
				data.MessageWarning = "Error! Unable to update profile in ledger."
				renderTemplate(w, r, "registerdetails.html", data)
			}
			applicant.Profile.PublicKey = string(publicKeyBytes)
		}

		// Insert row into DB
		userDetails, err := database.CreateNewUser(username, title, firstName, surname, emailAddress, dateOfBirth, fabricID)
		if err != nil {
			data.UserDetails.Username = u.Username
			data.MessageWarning = "Error! Unable to save user details."
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		// If the user is a verifier
		_, verifierErr := u.QueryVerifier()
		// Update verifier org
		if verifierErr == nil {
			newOrganisation := r.FormValue("organisation")
			if len(newOrganisation) != 0 {
				err := u.UpdateVerifierSaveOrganisation(newOrganisation)
				if err != nil {
					data.MessageWarning = "Error! Unable to update organisation name."
				}
			}
		}

		// Register the userDetails gob to be used as a session value
		gob.Register(userDetails)
		session.Values["SavedUserDetails"] = true
		session.Values["UserDetails"] = userDetails
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.UserDetails = userDetails
		data.MessageSuccess = "Success! Your details have been saved. Welcome, " + userDetails.FirstName + " " + userDetails.Surname + "."

		if applicantErr == nil {
			data.MessageWarning = "Before using the system, please make a copy of your Private Key."
			renderTemplate(w, r, "displaykey.html", data)
		} else {
			renderTemplate(w, r, "index.html", data)
		}
	})
}

func (c *Controller) UpdateDetailsView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "userdetails",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
			renderTemplate(w, r, "updatedetails.html", data)
		} else {
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
		}
	})
}

func (c *Controller) UpdateDetailsHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if !hasSavedUserDetails(session) {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		// Form values to insert into DB
		username := r.FormValue("username")
		title := r.FormValue("title")
		firstName := r.FormValue("firstName")
		surname := r.FormValue("surname")
		dateOfBirth := r.FormValue("dateOfBirth")
		emailAddress := r.FormValue("emailAddress")

		// Logic to update a profile
		userDetails, err := database.UpdateUser(username, title, firstName, surname, dateOfBirth, emailAddress)
		if err != nil {
			data.MessageWarning = "Error! Unable to update profile information in database."
			data.CurrentPage = "userdetails"
			renderTemplate(w, r, "updatedetails.html", data)
		} else {
			// Successfully updated user
			// Update the session values and save session
			gob.Register(userDetails)
			session.Values["UserDetails"] = userDetails
			data.UserDetails = userDetails
			data.MessageSuccess = "Success! You details have been updated."
			renderTemplate(w, r, "index.html", data)
		}

	})
}
