package web

import (
	"fmt"
	"github.com/cvtracker/controllers"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))

	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.HandleFunc("/", app.IndexHandler)
	r.HandleFunc("/login.html", app.LoginView)
	r.HandleFunc("/loginProcess.html", app.LoginHandler)
	r.HandleFunc("/addCV.html", app.AddCVView)
	r.HandleFunc("/addCVProcess.html", app.AddCVHandler)
	r.HandleFunc("/updateCV.html", app.UpdateCVView)
	r.HandleFunc("/updateCVProcess.html", app.UpdateCVHandler)
	r.HandleFunc("/mycv.html", app.ResultHandler)

	r.HandleFunc("/submitForReview.html", app.SubmitForReviewHandler)
	r.HandleFunc("/withdrawFromReview.html", app.WithdrawFromReviewHandler)

	r.HandleFunc("/viewall.html", app.ViewAllView)
	r.HandleFunc("/logout.html", app.LogoutHandler)
	r.HandleFunc("/register.html", app.RegisterView)
	r.HandleFunc("/registerProcess.html", app.RegisterHandler)


	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/index.html")
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
