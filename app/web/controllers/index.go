package controllers

import (
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	"net/http"
)

func (c *Controller) IndexHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		data := templateModel.Data{
			CurrentPage:  "index",
		}

		session := sessions.GetSession(r)

		// Retrieve user details
		data.AccountType = sessions.GetAccountType(session)
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
			renderTemplate(w, r, "index.html", data)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
		}
	})
}
