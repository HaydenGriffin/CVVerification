package controllers

import (
	"fmt"
	"github.com/cvverification/app/database"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	"github.com/gorilla/mux"
	"net/http"
)

func (c *Controller) ViewCVApplications() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "index",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		// Check that the user connected is an employer
		employer, err := u.QueryEmployer()
		if err != nil {
			fmt.Println(err)
			data.CurrentPage = "index"
			data.MessageWarning = "You must be an employer user to view CV applications."
			renderTemplate(w, r, "index.html", data)
			return
		}

		industryFilter := r.FormValue("industry")
		registeredCVs := r.FormValue("registeredcvs")

		cvList, err := u.QueryCVs(model.CVSubmitted, industryFilter)
		if err != nil {
			data.MessageWarning = "An error occurred whilst retrieving CVs to view."
			renderTemplate(w, r, "index.html", data)
			return
		}

		var finalCVList = make(map[string]model.CVObject)

		if registeredCVs == "on" && len(employer.Profile.ProspectiveCVs) > 0 {
			for _, registeredCVID := range employer.Profile.ProspectiveCVs {
				for cvID, cv := range cvList {
					if registeredCVID == cvID {
						finalCVList[cvID] = cv
					}
				}
			}
			data.CVInfo.CVList = finalCVList
		} else if registeredCVs == "on" && len(employer.Profile.ProspectiveCVs) == 0  {
			data.CVInfo.CVList = finalCVList
		} else {
			data.CVInfo.CVList = cvList
		}

		if len(data.CVInfo.CVList) == 0 {
			data.MessageWarning = "There are no CVs to be viewed at this time."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if industryFilter != "" {
			data.MessageSuccess = "Showing results for: " + industryFilter
		}

		data.CurrentPage = "viewcvs"
		renderTemplate(w, r, "viewcvapplications.html", data)
	})
}

func (c *Controller) ViewCVView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "viewcvs",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		// Check that the user connected is an employer
		employer, err := u.QueryEmployer()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "You must be an employer user to view CVs."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cvID, success := mux.Vars(r)["cvID"]
		if !success {
			data.MessageWarning = "Error! No CV to be retrieved."
			renderTemplate(w, r, "index.html", data)
			return
		}

		applicantFabricID, err := database.GetFabricIDFromCVID(cvID)
		if err != nil {
			data.MessageWarning = "Error! Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cv, err := u.QueryCV(cvID)
		if err != nil {
			data.MessageWarning = "Error! Unable to find CV from ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		reviews, err := u.QueryCVReviews(applicantFabricID, cvID)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to find CV reviews from ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if cv.Status != model.CVSubmitted {
			data.MessageWarning = "Error! Unable to view CV: CV is not submitted."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if len(employer.Profile.ProspectiveCVs) > 0 {
			var cvFound = false
			for _, interestedCVID := range employer.Profile.ProspectiveCVs {
				if interestedCVID == cvID {
					cvFound = true
				}
			}
			if cvFound == true {
				applicantUserDetails, err := database.GetUserDetailsFromFabricID(applicantFabricID)
				if err != nil {
					fmt.Println(err)
					data.MessageWarning = "Error! Unable to retreive user details."
					renderTemplate(w, r, "index.html", data)
					return
				}
				data.MessageSuccess = "You have registered your interest with the CV. Please contact the applicant to make further arrangements."
				data.CVInfo.ApplicantDetails = applicantUserDetails
			} else {
				data.MessageSuccess = "Register your interest to view the contact details for the applicant."
			}
		} else {
			data.MessageSuccess = "Register your interest to view the contact details for the applicant."
		}

		data.CVInfo.ReviewInfo.Reviews = reviews
		data.CVInfo.CV = cv
		session.Values["ApplicantFabricID"] = applicantFabricID
		session.Values["CVID"] = cvID
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}
		renderTemplate(w, r, "viewcv.html", data)
	})
}

func (c *Controller) RegisterInterestHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "viewcvs",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
		} else {
			data.CurrentPage = "userdetails"
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		// Check that the user connected is an employer
		employer, err := u.QueryEmployer()
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "You must be an employer user to view CVs."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cvID := getCVID(session)
		if cvID == "" {
			data.MessageWarning = "Error! CV ID could not be found."
			renderTemplate(w, r, "index.html", data)
			return
		}

		applicantFabricID, err := database.GetFabricIDFromCVID(cvID)
		if err != nil {
			data.MessageWarning = "Error! Unable to find CV info in database."
			renderTemplate(w, r, "index.html", data)
			return
		}

		cv, err := u.QueryCV(cvID)
		if err != nil {
			data.MessageWarning = "Error! Unable to find CV from ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		reviews, err := u.QueryCVReviews(applicantFabricID, cvID)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to find CV reviews from ledger."
			renderTemplate(w, r, "index.html", data)
			return
		}

		if cv.Status != model.CVSubmitted {
			data.MessageWarning = "Error! Unable to view CV: CV is not submitted."
			renderTemplate(w, r, "index.html", data)
			return
		}

		fmt.Println(employer.Profile.ProspectiveCVs)

		err = u.UpdateEmployerSaveCV(cvID)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to register interest for CV."
			renderTemplate(w, r, "index.html", data)
			return
		}

		applicantUserDetails, err := database.GetUserDetailsFromFabricID(applicantFabricID)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to retreive user details."
			renderTemplate(w, r, "index.html", data)
			return
		}
		data.CVInfo.ApplicantDetails = applicantUserDetails

		data.CVInfo.ReviewInfo.Reviews = reviews
		data.CVInfo.CV = cv
		session.Values["ApplicantFabricID"] = applicantFabricID
		session.Values["CVID"] = cvID
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}
		data.MessageSuccess = "Success! You have registered your interest in the CV."
		renderTemplate(w, r, "viewcv.html", data)
	})
}
