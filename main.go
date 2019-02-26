/**
  author: Hayden Griffin
 */

package main

import (
	"encoding/json"
	"fmt"
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/crypto"
	"github.com/cvtracker/service"
	"github.com/cvtracker/web"
	"github.com/cvtracker/controllers"
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
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

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

	serviceSetup := service.ServiceSetup{
		ChaincodeID:"cvtracker",
		Client:channelClient,
	}

	profile := service.UserProfile{
		Username: "admin",
	}

	userHash, err := crypto.GenerateFromString("admin")

	result, err := serviceSetup.SaveProfile(profile, userHash)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully saved profile: " + result)
	}

	cv := service.CVObject{
		Speciality: "Test",
		CV:"The greatest CV of them all! Bow down to me peasants",
	}

	cvByte, err := json.Marshal(cv)

	cvHash, err := crypto.GenerateFromByte(cvByte)

	result, err = serviceSetup.SaveCV(cv, cvHash)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully saved CV: " + result)
	}

	result, err = serviceSetup.UpdateProfile(userHash, cvHash)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully updated profile: " + result)


	profilebyte, err := serviceSetup.GetProfile(userHash)

	json.Unmarshal(profilebyte, &profile)
	fmt.Println(profile)

	b, err := serviceSetup.QueryCVFromProfileHash(userHash)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			var cv1 service.CVObject
			json.Unmarshal(b, &cv1)
			fmt.Println(cv1)
		}
	}



	// Launch the web application listening
	app := &controllers.Application{
		Service: &serviceSetup,
	}
	web.Serve(app)
}