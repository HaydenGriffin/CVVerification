package controllers

import (
	"fmt"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/blockchain"
	"net/http"
)

func (c *Controller) IndexHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
			renderTemplate(w, r, "index.html", data)
		} else {
			data.CurrentPage = "userdetails"
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
		}
	})
}
