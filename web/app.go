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
	r.HandleFunc("/addcv.html", app.AddCVView())
	r.HandleFunc("/addcvprocess.html", app.AddCVHandler())
	r.HandleFunc("/updatecv.html", app.UpdateCVView)
	r.HandleFunc("/updatecvprocess.html", app.UpdateCVHandler)
	r.HandleFunc("/mycv.html", app.ViewCVHandler())

	//r.HandleFunc("/submitcvforreview.html", app.SubmitForReviewHandler)
	r.HandleFunc("/withdrawcvfromreview.html", app.WithdrawFromReviewHandler)

	r.HandleFunc("/viewallcv.html", app.ViewAllCVView)

	r.HandleFunc("/ratecv/{userID}.html", app.RateCVView)
	r.HandleFunc("/ratecvprocess.html", app.RateCVHandler)

	r.HandleFunc("/logout.html", app.LogoutHandler)
	r.HandleFunc("/userdetails.html", app.UpdateDetailsView())
	r.HandleFunc("/updatedetailsprocess.html", app.UpdateDetailsHandler())


	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/index.html")
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
