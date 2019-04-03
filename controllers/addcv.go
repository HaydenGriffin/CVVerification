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
			data.MessageWarning = "You must be an applicant user to upload a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
			renderTemplate(w, r, "cvform.html", data)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
			data.MessageWarning = "You must register your user details before using the system."
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
			data.MessageWarning = "You must be an applicant user to upload a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
			data.MessageWarning = "You must register your user details before using the system."
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

		cvHash, err := crypto.GenerateFromByte(cvByte)

		err = u.UpdateSaveCV(cvByte, cvHash)

		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "An error occurred whilst saving the CV to ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		err = u.UpdateSaveProfileCV(data.UserDetails.ProfileHash, cvHash)

		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "An error occurred whilst updating profile information in ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		err = database.CreateNewCV(data.UserDetails.Id, cvHash)

		if err != nil {
			data.MessageWarning = "An error occurred whilst saving CV details to database."
			renderTemplate(w, r, "addcv.html", data)
		} else {
			data.MessageSuccess = "You have successfully saved your CV to the ledger."
			renderTemplate(w, r, "index.html", data)

		}
	})
}
