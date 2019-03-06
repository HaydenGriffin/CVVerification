/**
  author: Hayden Griffin
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/chaincode/model"
	"github.com/cvtracker/controllers"
	"github.com/cvtracker/crypto"
	"github.com/cvtracker/database"
	"github.com/cvtracker/service"
	"github.com/cvtracker/web"
	"os"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters 
		OrdererID: "orderer.cvtracker.com",

		// Channel parameters
		ChannelID:     "cvtracker",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/cvtracker/fixtures/artifacts/cvtracker.channel.tx",

		// Chaincode parameters
		ChaincodeID:     "cvtracker",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/cvtracker/chaincode/",
		ChaincodeVersion: "v1.0.0",
		OrgAdmin:        "Admin",
		OrdererOrgID:    "ordererorg",
		OrgMspID:        "org1.cvtracker.com",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",
		CaID: 			 "ca.org1.cvtracker.com",

		// User parameters
		UserName: "User1",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()	

	// Install and instantiate the chaincode
	channelClient, err := fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

	err = fSetup.RegisterUser("admin1", "password", model.ActorAdmin)
	if err != nil {
		fmt.Printf("Unable to register the user 'admin1': %v\n", err)
		return
	}

	serviceSetup := service.ServiceSetup{
		ChaincodeID:"cvtracker",
		Client:channelClient,
	}



	//Init a dummy user and test chaincode methods
	profile := service.UserProfile{
		Username: "testUser",
	}

	var profile1 service.UserProfile

	profile.Ratings = make(map[string][]service.CVRating)

	userHash, err := crypto.GenerateFromString("testUser")

	passwordHash, err := crypto.GenerateFromString("password")

	user, err := database.CreateNewUser("testUser", "Test User", passwordHash, "test@user.com", "APPLICANT", userHash)

	result, err := serviceSetup.SaveProfile(profile, userHash)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully saved profile: " + result)
	}

	cv := service.CVObject{
		Name:"Test User",
		Speciality: "Test Speciality",
		CV:"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}

	cvByte, err := json.Marshal(cv)

	cvHash, err := crypto.GenerateFromByte(cvByte)

	result, err = serviceSetup.SaveCV(cv, cvHash)
	err = database.CreateNewCV(user.Id, cv.CV, cvHash)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully saved CV: " + result)
	}

	result, err = serviceSetup.UpdateProfileCV(userHash, cvHash)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully updated profile: " + result)


		profilebyte, err := serviceSetup.GetProfile(userHash)

		json.Unmarshal(profilebyte, &profile)
		fmt.Println(profile)

		b, err := serviceSetup.GetCVFromProfile(userHash)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			var cv1 service.CVObject
			json.Unmarshal(b, &cv1)
			fmt.Println(cv1)
		}
	}

	rating :=  service.CVRating{
		Name: "Test Rater",
		Comment: "Test Comment",
		Rating: 4,
	}

	txid, err := serviceSetup.SaveRating(userHash, cvHash, rating)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("TXID: " + txid)
	}

	profilebyte, err := serviceSetup.GetProfile(userHash)

	err = json.Unmarshal(profilebyte, &profile1)

	if err != nil {
		fmt.Println("An error occurred whilst retrieving the profile: " + err.Error())
	} else {
		fmt.Println(profile1)
	}

	// Launch the web application listening
	app := &controllers.Application{
		Service: &serviceSetup,
	}
	web.Serve(app)
}