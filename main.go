/**
  author: Hayden Griffin
 */

package main

import (
	"fmt"
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/chaincode/model"
	"github.com/cvtracker/controllers"
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
	_, err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

	_, err = fSetup.LogUser("admin", "adminpw")
	if err != nil {
		fmt.Printf("failed to enroll identity 'admin': %v", err)
		return
	}

	err = fSetup.RegisterUser("admin1", "password", model.ActorAdmin)
	if err != nil {
		fmt.Printf("Unable to register the user 'admin1': %v\n", err)
		return
	}

	err = fSetup.RegisterUser("applicant1", "password", model.ActorApplicant)
	if err != nil {
		fmt.Printf("Unable to register the user 'applicant1': %v\n", err)
		return
	}

	err = fSetup.RegisterUser("verifier1", "password", model.ActorVerifier)
	if err != nil {
		fmt.Printf("Unable to register the user 'verifier1': %v\n", err)
		return
	}

	// Launch the web application listening
	app := &controllers.Controller{
		Fabric: &fSetup,
	}
	web.Serve(app)
}