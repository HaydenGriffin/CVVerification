package controllers

import (
	"github.com/cvtracker/models"
	"github.com/cvtracker/sessions"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)


func (app *Controller) UpdateCVView(w http.ResponseWriter, r *http.Request) {
	//session := sessions.InitSession(r)

	data := models.TemplateData{
		//CurrentUser:models.User{},
		CurrentPage:"addcv",
	}

	/*if sessions.IsLoggedIn(session) {
		data.CurrentUser = sessions.GetUser(session)
	} else {
		data.LoggedInFlag = false
		data.MessageWarning = "Error! Please log in to update your CV."
		renderTemplate(w, r, "index.html", nil)
		return
	}

	b, err := app.Service.GetCVFromProfile(data.CurrentUser.ProfileHash)

	if err != nil {
		fmt.Printf(err.Error())
		data.MessageWarning = "Unable to find CV from hash."
		renderTemplate(w, r, "index.html", data)
	} else {
		var cv= service.CVObject{}
		err = json.Unmarshal(b, &cv)
		data.CV = cv*/
		renderTemplate(w, r, "cvform.html", data)
}

func (app *Controller) UpdateCVHandler(w http.ResponseWriter, r *http.Request) {

	session := sessions.InitSession(r)

	data := models.TemplateData{
		CurrentPage:"addcv",
	}

	if sessions.IsLoggedIn(session) {
		//data.UserDetails = sessions.GetUserDetails(session)
	} else {
		data.MessageWarning = "Error! Please log in to add a CV."
		renderTemplate(w, r, "index.html", nil)
		return
	}


	/*cv := model.CVObject{
		Name:r.FormValue("name"),
		Speciality:r.FormValue("speciality"),
		CV:r.FormValue("CV"),
		CVDate:r.FormValue("CVDate"),
	}*/


	/*cvByte, err := json.Marshal(cv)

	cvHash, err := crypto.GenerateFromByte(cvByte)

	txid, err := app.Service.SaveCV(cv, cvHash)

	txid, err = app.Service.UpdateProfileCV(data.CurrentUser.ProfileHash, cvHash)

	err = database.CreateNewCV(data.CurrentUser.Id, cv.CV, cvHash)

	if err != nil {
		data.MessageWarning = "Unable to invoke"
		renderTemplate(w, r, "index.html", data)
	} else {
		data.MessageSuccess = txid
		data.CV = cv*/
		renderTemplate(w, r, "mycv.html", data)
}