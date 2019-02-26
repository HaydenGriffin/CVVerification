package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/cvtracker/database"
	"github.com/cvtracker/models"
	"github.com/cvtracker/service"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)


func (app *Application) UpdateCVView(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"addcv",
		LoggedInFlag:true,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to update a CV."
		renderTemplate(w, r, "index.html", nil)
		return
	}

	cvHash, err := database.GetCVHashFromUserID(data.CurrentUser.Id)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Unable to find CV from hash."
	} else {
		cvString, err := app.Service.QueryCVByHash(cvHash)
		if err != nil {
			data.MessageWarning = "Unable to query the blockchain"
		} else {
			var cv = service.CVObject{}
			err := json.Unmarshal(cvString, &cv)
			if err != nil {
				data.MessageWarning = "Unable to unmarshal CV object"
			}
			data.CV = cv
		}
		renderTemplate(w, r, "addcv.html", data)
	}

}

func (app *Application) UpdateCVHandler(w http.ResponseWriter, r *http.Request) {

	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"addcv",
		LoggedInFlag:true,
	}

	if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to add a CV."
		renderTemplate(w, r, "index.html", nil)
	}


	cv := service.CVObject{
		Name:r.FormValue("name"),
		Speciality:r.FormValue("speciality"),
		CV:r.FormValue("CV"),
		CVDate:r.FormValue("CVDate"),
	}

	txid, err := app.Service.ModifyCV(cv)


	db, err := database.InitDB("root:password@tcp(localhost:3306)/verification")

	res, err := db.Exec("UPDATE user_cvs SET timestamp = CURRENT_TIMESTAMP, cv = ? WHERE user_id = ?", cv.CV, data.CurrentUser.Id)


	if err != nil {
		panic(err)
	} else {
		fmt.Println(res)
	}

	if err != nil {
		data.MessageWarning = "Unable to invoke"
		renderTemplate(w, r, "mycv.html", data)
	} else {
		data.MessageSuccess = txid
		renderTemplate(w, r, "index.html", data)
	}
}