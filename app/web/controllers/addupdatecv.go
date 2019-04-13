package controllers

import (
	"encoding/json"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/app/sessions"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
)

func (c *Controller) AddCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "addcv",
		}

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to upload a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}
		renderTemplate(w, r, "addcv.html", data)
	})
}

func (c *Controller) UpdateCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		// Check that the user connected is an applicant
		applicant, err := u.QueryApplicant()
		if err != nil {
			data.MessageWarning = "Error! You must be an applicant user to update your CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// If the user hasn't uploaded a CV then take them to addcv page
		if data.UserDetails.UploadedCV == false {
			data.CurrentPage = "addcv"
			data.MessageWarning = "Error! You must add a CV before you can update it."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		cvToDisplay := sessions.GetCV(session)

		// User is able to update a selected CV if specified. Otherwise update latest
		if cvToDisplay == nil {
			if len(applicant.Profile.CVHistory) != 0 {
				cvToDisplayCVID := applicant.Profile.CVHistory[len(applicant.Profile.CVHistory)-1]
				cvToDisplay, err = u.QueryCV(cvToDisplayCVID)
				if err != nil {
					data.MessageWarning = "Error! Something went wrong whilst retrieving CV details from ledger."
					renderTemplate(w, r, "index.html", data)
					return
				}
			} else {
				data.MessageWarning = "Error! Unable to retrieve CV from ledger."
				renderTemplate(w, r, "index.html", data)
				return
			}

		}

		data.CVInfo.CV = cvToDisplay
		data.CurrentPage = "updatecv"
		renderTemplate(w, r, "updatecv.html", data)
	})
}

func (c *Controller) AddCVHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session := sessions.InitSession(r)

		data := templateModel.Data{
			CurrentPage: "addcv",
		}

		// Retrieve user details
		if sessions.HasSavedUserDetails(session) {
			data.UserDetails = sessions.GetUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		// Check that the user connected is an applicant
		_, err := u.QueryApplicant()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to upload a CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Extract form values and create new object
		cv := model.CVObject{
			DocType: "cv",
			Name: r.FormValue("name"),
			Speciality: r.FormValue("speciality"),
			CVDate: r.FormValue("cvDate"),
			CV: r.FormValue("mainCVSectionValue"),
		}

		// Additional sections are stored in a map
		cv.CVSections = make(map[string]string)

		// Parse the form to retreive all the form data
		if err := r.ParseForm(); err != nil {
			data.MessageWarning = "Error! Something went wrong whilst processing the request."
			renderTemplate(w, r, "index.html", data)
			return
		}

		var listOfIndexes []string

		// Additional sections can be added and deleted dynamically
		// The key here is the name of the form field
		for key, _ := range r.PostForm {
			// Retrieve all additional section form fields
			if strings.Contains(key, "additionalCVSectionSubject") {
				// Extract the last character in the string (this is a number set by jQuery)
				index := key[len(key)-1:]
				// Add the number to the list of indexes
				listOfIndexes = append(listOfIndexes, index)
			}
		}

		// For each number in the list of indexes, there is a form value
		// for the subject and the value
		for _, index := range listOfIndexes {
			// Extract the additional section subject and values
			key := r.PostForm.Get("additionalCVSectionSubject"+index)
			value := r.PostForm.Get("additionalCVSectionValue"+index)
			// Add the key and value to the CVSections map
			cv.CVSections[key] = value
		}

		// Convert the CV object to byte (for storage on ledger)
		cvByte, err := json.Marshal(cv)
		if err != nil {
			data.MessageWarning = "Error! Failed to save CV to ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		// Generate hash for the CV object
		// The hash is used as the key to access the object on the ledger
		cvID, err := c.ShortID.Generate()
		if err != nil {
			data.MessageWarning = "Error! Failed to save CV to ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		// Save the CV object to ledger
		err = u.UpdateSaveCV(cvByte, cvID)
		if err != nil {
			data.MessageWarning = "Error! Failed to save CV to ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		// Update the users profile CV history to include the latest CV ID
		err = u.UpdateSaveProfileCV(cvID)
		if err != nil {
			data.MessageWarning = "Error! Unable to update profile information in ledger."
			renderTemplate(w, r, "addcv.html", data)
			return
		}

		// Create a new DB table row with info about the CV saved
		err = database.CreateNewCV(data.UserDetails.Id, cvID)
		if err != nil {
			data.MessageWarning = "Error! Unable to save CV details to database."
			renderTemplate(w, r, "addcv.html", data)
		} else {
			// Set session values
			session.Values["UserUploadedCV"] = true
			err = session.Save(r, w)
			if err != nil {
				data.MessageWarning = "Error! Unable to save session values."
				renderTemplate(w, r, "index.html", data)
				return
			}
			data.CurrentPage = "index"
			data.MessageSuccess = "Success! Your CV has been saved to the ledger."
			renderTemplate(w, r, "index.html", data)

		}
	})
}
