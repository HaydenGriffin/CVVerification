package controllers

import (
	"bytes"
	"fmt"
	"github.com/cvverification/app/crypto"
	templateModel "github.com/cvverification/app/model"
	"github.com/cvverification/blockchain"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func (c *Controller) ApplicantKeyView() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "userdetails",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
		} else {
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		_, err = u.QueryApplicant()
		// User is not an applicant
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to access this page."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.PrivateKey = getPrivateKey(session)

		renderTemplate(w, r, "displaykey.html", data)
	})
}

func (c *Controller) UploadPrivateKeyHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "userdetails",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
		} else {
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		_, err = u.QueryApplicant()
		// User is not an applicant
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to access this page."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Initialise private key, in case something goes wrong during upload
		data.PrivateKey = getPrivateKey(session)

		const MAX_MEMORY = 1 * 1024 * 1024
		privateKeyBuffer := bytes.NewBuffer(nil)

		if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
			data.MessageWarning = "Error! Uploaded file is too large."
			renderTemplate(w, r, "displaykey.html", data)
			return
		}

		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			data.MessageWarning = "Error! No file selected for upload."
			renderTemplate(w, r, "displaykey.html", data)
			return
		}
		defer file.Close()

		if _, err := io.Copy(privateKeyBuffer, file); err != nil {
			data.MessageWarning = "Error! Something went wrong whilst downloading the file."
			renderTemplate(w, r, "displaykey.html", data)
			return
		}

		fileType := http.DetectContentType(privateKeyBuffer.Bytes())

		switch fileType {
		case "text/plain; charset=utf-8":
			data.PrivateKey = string(privateKeyBuffer.Bytes())
			session.Values["PrivateKey"] = string(privateKeyBuffer.Bytes())
			data.MessageSuccess = "Success! Private key has been updated."
		default:
			data.MessageWarning = "Error! Unsupported filetype uploaded."
		}

		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "displaykey.html", data)
			return
		}

		renderTemplate(w, r, "displaykey.html", data)
	})
}

func (c *Controller) DownloadPrivateKeyHandler() func(http.ResponseWriter, *http.Request) {
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
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		_, err = u.QueryApplicant()
		// User is not an applicant
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to access this page."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Form values
		privateKey := r.FormValue("privateKey")

		var filename = u.Username + "-privatekey.pem"

		err = ioutil.WriteFile(filename, []byte(privateKey), 0755)
		if err != nil {
			data.MessageWarning = "Error! Something went wrong whilst downloading the Private Key file."
			renderTemplate(w, r, "displaykey.html", data)
			return
		}

		file, err := ioutil.ReadFile(filename)
		if err != nil {
			data.MessageWarning = "Error! Something went wrong whilst downloading the Private Key file."
			renderTemplate(w, r, "displaykey.html", data)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Expires", "0")
		http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(file))

		err = os.Remove(filename)
		if err != nil {
			fmt.Println(err)
		}
	})
}

func (c *Controller) GenerateNewKeysHandler() func(http.ResponseWriter, *http.Request) {
	return c.basicAuth(func(w http.ResponseWriter, r *http.Request, u *blockchain.User) {

		session, err := store.Get(r, "userSession")
		if err != nil {
			fmt.Println(err)
		}

		data := templateModel.Data{
			CurrentPage: "userdetails",
		}

		// Retrieve user details
		data.AccountType = getAccountType(session)
		if hasSavedUserDetails(session) {
			data.UserDetails = getUserDetails(session)
		} else {
			data.MessageWarning = "Error! You must register your user details before using the system."
			data.UserDetails.Username = u.Username
			renderTemplate(w, r, "registerdetails.html", data)
			return
		}

		applicant, err := u.QueryApplicant()
		// User is not an applicant
		if err != nil {
			data.CurrentPage = "index"
			data.MessageWarning = "Error! You must be an applicant user to access this page."
			renderTemplate(w, r, "index.html", data)
			return
		}

		// Initialise private key display, in case something goes wrong during generation
		data.PrivateKey = getPrivateKey(session)

		privateKey, publicKey := crypto.GenerateKeyPair(2048)
		privateKeyBytes := crypto.PrivateKeyToBytes(privateKey)
		privateKeyString := string(privateKeyBytes)
		publicKeyBytes := crypto.PublicKeyToBytes(publicKey)
		applicant.Profile.PublicKey = string(publicKeyBytes)

		err = u.UpdateSaveProfileKey(string(publicKeyBytes))
		if err != nil {
			data.MessageWarning = "Error! Unable to update profile in ledger."
			renderTemplate(w, r, "displaykey.html", data)
			return
		}

		session.Values["PrivateKey"] = privateKeyString
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			data.MessageWarning = "Error! Unable to save session values."
			renderTemplate(w, r, "index.html", data)
			return
		}

		data.PrivateKey = privateKeyString
		data.MessageSuccess = "Success! You have generated a new Private Key."
		renderTemplate(w, r, "displaykey.html", data)
	})
}
