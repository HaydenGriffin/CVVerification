package controllers

import (
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

func (c *Controller) AddCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "addcv",
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			fmt.Println(err)
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to upload a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
			renderTemplate(w, r, "addcv.html", data)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
			data.MessageWarning = "Error! You must register your user details before using the system."
			renderTemplate(w, r, "userdetails.html", data)
		}
	})
}

func (c *Controller) UpdateCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "addcv",
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			data.MessageWarning = "Error! You must be an applicant user to update your CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
			data.CVInfo.CV = sessions.GetCV(session)
			fmt.Println(data.CVInfo.CV)
			renderTemplate(w, r, "updatecv.html", data)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
			data.MessageWarning = "Error! You must register your user details before using the system."
			renderTemplate(w, r, "userdetails.html", data)
		}

	})
}

func (c *Controller) AddCVHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "addcv",
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to upload a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
			data.MessageWarning = "Error! You must register your user details before using the system."
			renderTemplate(w, r, "userdetails.html", data)
			return
		}

		cv := model.CVObject{
			Name:       r.FormValue("name"),
			Speciality: r.FormValue("speciality"),
			CV:         r.FormValue("cv"),
			CVDate:     r.FormValue("cvDate"),
		}

		cvByte, err := json.Marshal(cv)
		if err != nil {
			data.MessageWarning = "Error! Failed to save CV to ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		cvHash, err := crypto.GenerateFromByte(cvByte)
		if err != nil {
			data.MessageWarning = "Error! Failed to save CV to ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		err = u.UpdateSaveCV(cvByte, cvHash)
		if err != nil {
			data.MessageWarning = "Error! Unable to save CV to ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		err = u.UpdateSaveProfileCV(cvHash)
		if err != nil {
			data.MessageWarning = "Error! Unable to update profile information in ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		err = database.CreateNewCV(data.UserDetails.Id, cvHash)
		if err != nil {
			data.MessageWarning = "Error! Unable to save CV details to database."
			renderTemplate(w, r, "addcv.html", data)
		} else {
			data.MessageSuccess = "Success! Your CV has been saved to the ledger."
			renderTemplate(w, r, "index.html", data)

		}
	})
}
