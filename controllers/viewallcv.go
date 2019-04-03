package controllers

import (
	"fmt"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	"github.com/cvverification/database"
	"github.com/cvverification/models"
	"github.com/cvverification/sessions"
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

		// Check that the user connected is a verifier
		_, err := u.QueryVerifier()
		if err != nil {
			fmt.Println(err)
			data.CurrentPage = "index"
			data.MessageWarning = "You must be a verifier user to review CVs."
			renderTemplate(w, r, "index.html", data)
			return
		}

		reviewableCVs := make(map[int]string)

		reviewableCVs, err = database.GetAllReviewableCVHashes()
		fmt.Println(reviewableCVs)

		if err != nil {
			data.MessageWarning = err.Error()
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.CVInfo.CVList = make(map[int]*model.CVObject)

		for userID, cvHash := range reviewableCVs {
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
			data.MessageWarning = "There are no CVs to be reviewed at this time."
			renderTemplate(w, r, "viewallcv.html", data)
			return
		}

		renderTemplate(w, r, "viewallcv.html", data)
	})
}
