package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
  _ "github.com/go-sql-driver/mysql"
	"github.com/cvverification/models"
)

var cuser models.User

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "login.html", nil)
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	db, err := models.InitDB("root:password@tcp(localhost:3306)/verification")


	var user models.User
	var loggedInFlag bool

	result := db.QueryRow("SELECT u.id, u.username, u.full_name, u.email_address, u.user_role  FROM users u WHERE username = ? AND password = ?", username, password)
	err = result.Scan(&user.Id, &user.Username, &user.FullName, &user.EmailAddress, &user.UserRole)

	if err != nil {

		if err == sql.ErrNoRows {
			fmt.Printf("No row found \n")
		} else {
			panic(err)
		}
	} else {
		loggedInFlag = true
	}

	data := &struct {
		CurrentUser models.User
		LoggedInFlag bool
	}{
		CurrentUser:user,
		LoggedInFlag:false,
	}

	if loggedInFlag {
		// Login successful
		data.LoggedInFlag = true
		renderTemplate(w, r, "index.html", data)
	}else{
		// Login failed
		data.CurrentUser.Username = username
		renderTemplate(w, r, "login.html", data)
	}
}

// Logout
func (app *Application) LogOut(w http.ResponseWriter, r *http.Request)  {
	cuser = models.User{}
	renderTemplate(w, r, "login.html", nil)
}