package web

import (
	"fmt"
	"github.com/cvverification/app/web/controllers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Serve(app *controllers.Controller) {
	fs := http.FileServer(http.Dir("app/web/assets/"))

	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.HandleFunc("/", app.IndexHandler())
	r.HandleFunc("/logout", app.LogoutHandler())

	r.HandleFunc("/addcv", app.AddCVView())
	r.HandleFunc("/addcvprocess", app.AddCVHandler())

	r.HandleFunc("/updatecv", app.UpdateCVView())
	r.HandleFunc("/updatecvprocess", app.AddCVHandler())

	r.HandleFunc("/mycv", app.MyCVView())
	r.HandleFunc("/mycv/{requestedCVIndex}", app.MyCVView())

	r.HandleFunc("/submitcvforreview", app.SubmitForReviewHandler())
	r.HandleFunc("/withdrawcvfromreview", app.WithdrawFromReviewHandler())

	r.HandleFunc("/viewallcv", app.ViewAllCVView())
	r.HandleFunc("/viewallcv/{speciality}", app.ViewAllCVView())

	r.HandleFunc("/reviewcv/{cvID}", app.ReviewCVView())
	r.HandleFunc("/reviewcvprocess", app.ReviewCVHandler())

	r.HandleFunc("/userdetails", app.UpdateDetailsView())
	r.HandleFunc("/updatedetailsprocess", app.UpdateDetailsHandler())


	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/index.html")
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
