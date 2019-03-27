package web

import (
	"fmt"
	"github.com/cvtracker/controllers"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func Serve(app *controllers.Controller) {
	fs := http.FileServer(http.Dir("web/assets"))

	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.HandleFunc("/", app.IndexHandler())
	r.HandleFunc("/addcv", app.AddCVView())
	r.HandleFunc("/addcvprocess", app.AddCVHandler())
	r.HandleFunc("/updatecv", app.UpdateCVView)
	r.HandleFunc("/updatecvprocess", app.UpdateCVHandler)
	r.HandleFunc("/mycv", app.MyCVHandler())
	r.HandleFunc("/mycv/{cvToDisplayID}", app.MyCVHandler())

	r.HandleFunc("/submitcvforreview", app.SubmitForReviewHandler())
	r.HandleFunc("/withdrawcvfromreview", app.WithdrawFromReviewHandler())

	r.HandleFunc("/viewallcv", app.ViewAllCVView())

	r.HandleFunc("/ratecv/{userID}", app.RateCVView())
	r.HandleFunc("/ratecvprocess", app.RateCVHandler())

	r.HandleFunc("/logout", app.LogoutHandler)
	r.HandleFunc("/userdetails", app.UpdateDetailsView())
	r.HandleFunc("/updatedetailsprocess", app.UpdateDetailsHandler())


	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/index.html")
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
