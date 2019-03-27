package controllers

import (
	"fmt"
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/chaincode/model"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"net/http"
)


func (c *Controller) ViewAllCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {
		session := sessions.InitSession(r)

		data := models.TemplateData{
			CurrentPage: "index",
		}

		if sessions.IsLoggedIn(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.MessageWarning = "You must be logged in to view the CVs."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Check that the user connected is an admin
		_, err := u.QueryAdmin()
		if err != nil {
			fmt.Println(err)
			data.CurrentPage = "index"
			data.MessageWarning = "You must be an admin user to view CVs."
			renderTemplate(w, r, "index.html", data)
			return
		}

		ratableCVs := make(map[int]string)

		ratableCVs, err = database.GetAllRatableCVHashes()
		fmt.Println(ratableCVs)

		if err != nil {
			data.MessageWarning = err.Error()
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CVList = make(map[int]*model.CVObject)

		for userID, cvHash := range ratableCVs {
			fmt.Println("profileHash: " + string(userID))
			fmt.Println("cvHash: " + cvHash)
			cv, err := u.QueryCV(cvHash)

			if err != nil {
				fmt.Println(err)
				data.MessageWarning = "Unable to retrieve CV detail from ledger."
				renderTemplate(w, r, "index.html", data)
				return
			}

			data.CVInfo.CVList[userID] = cv
		}

		if len(data.CVInfo.CVList) == 0 {
			data.MessageWarning = "There are no CVs to be rated at this time."
			renderTemplate(w, r, "viewallcv.html", data)
			return
		}

		renderTemplate(w, r, "viewallcv.html", data)
	})
}
