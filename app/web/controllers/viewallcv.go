package controllers

import (
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	"net/http"
)

func (c *Controller) ViewAllCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		if c.UserSession.Session == nil {
			c.UserSession.Session = c.UserSession.GetSession(r,w,u)
		}

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = sessions.GetAccountType(c.UserSession.Session)
		if sessions.HasSavedUserDetails(c.UserSession.Session) {
			data.UserDetails = sessions.GetUserDetails(c.UserSession.Session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
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

		industryFilter := r.FormValue("industry")

		cvList, err := u.QueryCVs(model.CVInReview, industryFilter)
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


		if industryFilter!= "" {
			data.MessageSuccess = "Showing results for " +industryFilter
		}

		data.CurrentPage = "viewallcv"
		renderTemplate(w, r, "viewallcv.html", data)
	})
}
