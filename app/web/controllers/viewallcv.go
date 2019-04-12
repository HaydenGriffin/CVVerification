package controllers

import (
	"encoding/gob"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	"net/http"
)

func (c *Controller) ViewAllCVView() func(http.ResponseWriter, *http.Request) {
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

		// Check that the user connected is a verifier
		_, err := u.QueryVerifier()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "You must be a verifier user to review CVs."
			renderTemplate(w, r, "index.html", data)
			return
		}

		specialityFilter:= r.FormValue("speciality")

		cvList, err := u.QueryCVs(model.CVInReview, specialityFilter)
		if err != nil {
			data.MessageWarning = "An error occurred whilst retrieving CVs to review."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CVList = cvList

		if len(data.CVInfo.CVList) == 0 {
			data.MessageWarning = "There are no CVs to be reviewed at this time."
			renderTemplate(w, r, "index.html", data)
			return
		}

		gob.Register(data.CVInfo.CVList)
		session.Values["AllCVList"] = data.CVInfo.CVList
		err = session.Save(r, w)
		if err != nil {
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if specialityFilter!= "" {
			data.MessageSuccess = "Showing results for " +specialityFilter
		}
		data.CurrentPage = "viewallcv"
		renderTemplate(w, r, "viewallcv.html", data)
	})
}
