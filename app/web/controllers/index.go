package controllers

import (
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	"net/http"
)

func (c *Controller) IndexHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		if c.UserSession.Session == nil {
			c.UserSession.Session = c.UserSession.GetSession(r,w,u)
		}

		data := templateModel.Data{
			CurrentPage:  "index",
		}

		// Retrieve user details
		data.AccountType = sessions.GetAccountType(c.UserSession.Session)
		if sessions.HasSavedUserDetails(c.UserSession.Session) {
			data.UserDetails = sessions.GetUserDetails(c.UserSession.Session)
			renderTemplate(w, r, "index.html", data)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
		}
	})
}
