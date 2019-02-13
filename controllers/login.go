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

var cuser models.User

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := &struct {
		CurrentUser models.User
		LoggedInFlag bool
	}{
		CurrentUser:models.User{},
		LoggedInFlag:true,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
		renderTemplate(w, r, "index.html", data)
	} else {
		renderTemplate(w, r, "login.html", nil)
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

	data := &struct {
		CurrentUser models.User
		LoggedInFlag bool
	}{
		CurrentUser:models.User{},
		LoggedInFlag:false,
	}

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No row found \n")
			data.CurrentUser.Username = username
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
		data.LoggedInFlag = true
		renderTemplate(w, r, "index.html", data)
	}
}

// Logout
func (app *Application) LogoutHandler(w http.ResponseWriter, r *http.Request)  {
	session := sessions.InitSession(r)
	session.Values["User"] = models.User{}
	session.Values["LoggedInFlag"] = false
	session.Save(r, w)

	renderTemplate(w, r, "login.html", nil)
}