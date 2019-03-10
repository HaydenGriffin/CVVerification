package controllers

import (
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	"net/http"
)

func (c *Controller) IndexHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		data := models.TemplateData{
			CurrentPage:  "index",
			LoggedInFlag: false,
		}

		session := sessions.InitSession(r)
		if sessions.IsLoggedIn(session) {
			data.LoggedInFlag = true

			if sessions.HasSavedUserDetails(session) {
				userDetails := sessions.GetUserDetails(session)
				data.UserDetails = userDetails
				renderTemplate(w, r, "index.html", data)
			} else {
				data.CurrentPage = "register"
				data.UserDetails.Username = u.Username
				renderTemplate(w, r, "register.html", data)
			}
		}
	})
}
