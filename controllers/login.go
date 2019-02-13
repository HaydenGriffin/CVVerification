package controllers

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := &struct {
		CurrentUser models.User
		CurrentPage string
		LoggedInFlag bool
		MessageWarning string
		MessageSuccess string
	}{
		CurrentUser:models.User{},
		CurrentPage:"login",
		LoggedInFlag:false,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
		data.LoggedInFlag = true
		data.CurrentPage = "index"
		data.MessageWarning = "You are already logged in!"
		renderTemplate(w, r, "index.html", data)
	} else {
		renderTemplate(w, r, "login.html", data)
	}
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	db, err := models.InitDB("root:password@tcp(localhost:3306)/verification")
	session := sessions.InitSession(r)

	var user models.User

	result := db.QueryRow("SELECT u.id, u.username, u.full_name, u.email_address, u.user_role  FROM users u WHERE username = ? AND password = ?", username, password)
	err = result.Scan(&user.Id, &user.Username, &user.FullName, &user.EmailAddress, &user.UserRole)

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"login",
		LoggedInFlag:false,
	}

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No row found \n")
			data.CurrentUser.Username = username
			data.MessageWarning = "Error! Incorrect username or password."
			renderTemplate(w, r, "login.html", data)
		} else {
			panic(err)
		}
	} else {
		gob.Register(user)
		session.Values["User"] = user
		session.Values["LoggedInFlag"] = true
		session.Save(r, w)
		data.CurrentUser = user
		data.CurrentPage = "index"
		data.LoggedInFlag = true
		data.MessageSuccess = "You have successfully logged in! Welcome, " + user.FullName
		renderTemplate(w, r, "index.html", data)
	}
}

// Logout
func (app *Application) LogoutHandler(w http.ResponseWriter, r *http.Request)  {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"login",
		LoggedInFlag:false,
	}

	session.Values["User"] = models.User{}
	session.Values["LoggedInFlag"] = false
	session.Save(r, w)

	data.MessageSuccess = "You have been successfully logged out."
	renderTemplate(w, r, "login.html", data)
}