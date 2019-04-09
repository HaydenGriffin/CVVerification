package controllers

import (
	"encoding/json"
	"github.com/cvverification/app/crypto"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (c *Controller) AddCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "addcv",
		}

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "userdetails.html", data)
			return
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to upload a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}
		renderTemplate(w, r, "addcv.html", data)
	})
}

func (c *Controller) UpdateCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "userdetails.html", data)
			return
		}

		// Check that the user connected is an applicant
		applicant, err := u.QueryApplicant()
		if err != nil {
			data.MessageWarning = "Error! You must be an applicant user to update your CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if data.UserDetails.UploadedCV == false {
			data.CurrentPage = "addcv"
			data.MessageWarning = "Error! You must add a CV before you can update it."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		cvToDisplay := sessions.GetCV(session)

		if cvToDisplay == nil {
			cvToDisplayCVHash := applicant.Profile.CVHistory[len(applicant.Profile.CVHistory)-1]
			cvToDisplay, err = u.QueryCV(cvToDisplayCVHash)
			if err != nil {
				data.MessageWarning = "Error! Something went wrong whilst retrieving CV details from ledger."
				renderTemplate(w, r, "index.html", data)
				return
			}
		}

		data.CVInfo.CV = cvToDisplay
		data.CurrentPage = "updatecv"
		renderTemplate(w, r, "updatecv.html", data)
	})
}

func (c *Controller) AddCVHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "addcv",
		}

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "userdetails.html", data)
			return
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to upload a CV."
			renderTemplate(w, r, "index.html", data)
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
			session.Values["UserUploadedCV"] = true
			err = session.Save(r, w)
			if err != nil {
				data.MessageWarning = "Error! Unable to save session values."
				renderTemplate(w, r, "index.html", data)
				return
			}
			data.CurrentPage = "index"
			data.MessageSuccess = "Success! Your CV has been saved to the ledger."
			renderTemplate(w, r, "index.html", data)

		}
	})
}
