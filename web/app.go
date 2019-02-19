package web

import (
	"fmt"
	"github.com/cvtracker/controllers"
	"net/http"
)

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/index.html", app.IndexHandler)
	http.HandleFunc("/login.html", app.LoginView)
	http.HandleFunc("/loginProcess.html", app.LoginHandler)
	http.HandleFunc("/addCV.html", app.AddCVView)
	http.HandleFunc("/addCVProcess.html", app.AddCVHandler)
	http.HandleFunc("/result.html", app.ResultHandler)
	http.HandleFunc("/logout.html", app.LogoutHandler)


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}
