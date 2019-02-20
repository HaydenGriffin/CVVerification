package controllers

import (
	"github.com/cvtracker/models"
	"net/http"
)



func (app *Application) ResultHandler(w http.ResponseWriter, r *http.Request) {


	data := models.TemplateData{
		CurrentUser:models.User{},
		CurrentPage:"result",
		LoggedInFlag:false,
	}

	helloValue, err := app.Service.QueryCVByHash()
	if err != nil {
		data.MessageWarning = "Unable to query the blockchain"
	} else {
		data.MessageSuccess = "helloValue: " + helloValue
	}
	renderTemplate(w, r, "mycv.html", data)
}

