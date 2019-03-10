package controllers

import (
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func (app *Controller) AddCVView(w http.ResponseWriter, r *http.Request) {
	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentPage:"addcv",
		LoggedInFlag:true,
	}

	if sessions.IsLoggedIn(session) {
		//data.UserDetails = sessions.GetUserDetails(session)
		renderTemplate(w, r, "cvform.html", data)
		return
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to add a CV."
		renderTemplate(w, r, "index.html", data)
		return
	}
}

func (app *Controller) AddCVHandler(w http.ResponseWriter, r *http.Request) {

	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentPage:  "addcv",
		LoggedInFlag: true,
	}

	if sessions.IsLoggedIn(session) {
		//data.UserDetails = sessions.GetUserDetails(session)
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to add a CV."
		renderTemplate(w, r, "index.html", nil)
		return
	}

	//fabricUser, err := app.Fabric.LogUser(data.UserDetails.Username, data.UserDetails)

	/*cv := model.CVObject{
		Name:       r.FormValue("name"),
		Speciality: r.FormValue("speciality"),
		CV:         r.FormValue("CV"),
		CVDate:     r.FormValue("CVDate"),
	}*/

//	cvByte, err := json.Marshal(cv)

	//cvHash, err := crypto.GenerateFromByte(cvByte)

	//txid, err := app.Service.SaveCV(cv, cvHash)

	//txid, err = app.Service.UpdateProfileCV(data.CurrentUser.ProfileHash, cvHash)

/*	if err != nil {
		data.MessageWarning = err.Error()
		renderTemplate(w, r, "index.html", data)
		return
	}

	//err = database.CreateNewCV(data.CurrentUser.Id, cv.CV, cvHash)

	if err != nil {
		data.MessageWarning = "Unable to create new CV"
		renderTemplate(w, r, "mycv.html", data)
	} else {
		//	data.MessageSuccess = txid*/
		renderTemplate(w, r, "index.html", data)

}
