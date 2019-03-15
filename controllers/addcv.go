package controllers

import (
	"encoding/json"
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/chaincode/model"
	"github.com/cvtracker/crypto"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (c *Controller) AddCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {
		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage:  "addcv",
			LoggedInFlag: true,
		}

		if sessions.IsLoggedIn(session) {
			if sessions.HasSavedUserDetails(session) {
				data.UserDetails = sessions.GetUserDetails(session)
				renderTemplate(w, r, "cvform.html", data)
			} else {
				data.CurrentPage = "register"
				data.UserDetails.Username = u.Username
				renderTemplate(w, r, "register.html", data)
			}
		}
	})
}

func (c *Controller) AddCVHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage:  "addcv",
			LoggedInFlag: true,
		}

		if sessions.IsLoggedIn(session) {
			if sessions.HasSavedUserDetails(session) {
				data.UserDetails = sessions.GetUserDetails(session)
				renderTemplate(w, r, "cvform.html", data)
			} else {
				data.CurrentPage = "register"
				data.UserDetails.Username = u.Username
				data.MessageWarning = "You must enter your user details before adding your CV."
				renderTemplate(w, r, "register.html", data)
				return
			}
		}


		//fabricUser, err := app.Fabric.LogUser(data.UserDetails.Username, data.UserDetails)
		cv := model.CVObject{
		Name:       r.FormValue("name"),
		Speciality: r.FormValue("speciality"),
		CV:         r.FormValue("CV"),
		CVDate:     r.FormValue("CVDate"),
	}

			cvByte, err := json.Marshal(cv)

		cvHash, err := crypto.GenerateFromByte(cvByte)

		txid, err := u.UpdateAddCV()

		//txid, err = app.Service.UpdateProfileCV(data.CurrentUser.ProfileHash, cvHash)

		/*	if err != nil {
		data.MessageWarning = err.Error()
		renderTemplate(w, r, "index.html", data)
		return
	}

	//err = database.CreateNewCV(data.CurrentUser.Id, cv.CV, cvHash)

	if err != nil {
		data.MessageWarning = "Unable to create new CV"
		renderTemplate(w, r, "mycv.html", data)
	} else {
		//	data.MessageSuccess = txid*/
		renderTemplate(w, r, "index.html", data)

	})
}
