package controllers

import (
	"fmt"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)


func (app *Controller) RateCVView(w http.ResponseWriter, r *http.Request) {

	fmt.Println("RateCVView")

	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentPage:  "addcv",
		LoggedInFlag: true,
	}

	if sessions.IsLoggedIn(session) {
		//data.UserDetails = sessions.GetUserDetails(session)
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to update your CV."
		renderTemplate(w, r, "index.html", data)
		return
	}

	res, success := mux.Vars(r)["userID"]

	if !success {
		data.MessageWarning = "Error! Please log in to update your CV."
		renderTemplate(w, r, "index.html", data)
		return
	}

	userID, err := strconv.Atoi(res)

	if err != nil {
		data.MessageWarning = "Error! Please log in to update your CV."
		renderTemplate(w, r, "index.html", data)
		return
	}

	profileHash, cvHash, err := database.GetCVInfoFromID(userID)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Unable to find CV info in database."
		renderTemplate(w, r, "index.html", data)
		return
	}

	//b, err := app.Service.GetCVFromCVHash(cvHash)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Unable to find CV from hash."
		renderTemplate(w, r, "index.html", data)
		return
	}

	//var cv= service.CVObject{}
	//err = json.Unmarshal(b, &cv)
	//data.CV = cv

	session.Values["selectedProfileHash"] = profileHash
	session.Values["selectedCVHash"] = cvHash

	err = session.Save(r, w)
	if err != nil {
		data.MessageWarning = err.Error()
		fmt.Println(err.Error())
		renderTemplate(w, r, "index.html", data)
		return
	}

	renderTemplate(w, r, "rateCV.html", data)
}

func (app *Controller) RateCVHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("RateCVHandler")

	session := sessions.InitSession(r)

	data := models.TemplateData{
		//CurrentUser:  models.User{},
		CurrentPage:  "addcv",
		LoggedInFlag: true,
	}

	if sessions.IsLoggedIn(session) {
		//data.CurrentUser = sessions.GetUser(session)
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to update your CV."
		renderTemplate(w, r, "index.html", nil)
		return
	}

	/*ratingInt, err := strconv.Atoi(r.FormValue("rating"))

if err != nil {
	data.MessageWarning = "Error! Rating must be a number."
	renderTemplate(w, r, "index.html", data)
	return
}

rating := model.CVRating{
	Name:r.FormValue("name"),
	Comment:r.FormValue("comment"),
	Rating:ratingInt,
}

profileHash := session.Values["selectedProfileHash"].(string)
cvHash := session.Values["selectedCVHash"].(string)
// some handling required to ensure profile is returned

txid, err := app.Service.SaveRating(profileHash, cvHash, rating)

if err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println("Successfully saved rating: " + txid)
}
*/
	//data.MessageSuccess = txid
	renderTemplate(w, r, "rateCV.html", data)
}