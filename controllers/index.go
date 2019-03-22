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
		}

		session := sessions.InitSession(r)
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
			renderTemplate(w, r, "index.html", data)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "userdetails.html", data)
		}
	})
}
