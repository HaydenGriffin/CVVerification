/**
  author: Hayden Griffin
 */

package main

import (
	"fmt"
	"github.com/cvverification/app/web"
	"github.com/cvverification/app/web/controllers"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/chaincode/model"
	"os"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: "orderer.cvverification.com",

		// Channel parameters
		ChannelID:     "channelall",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/cvverification/fabric-network/fixtures/artifacts/cvverification.channel.tx",

		// Chaincode parameters
		ChaincodeID:     "cvverification",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/cvverification/chaincode/",
		ChaincodeVersion: "v1.0.0",
		OrgAdmin:        "Admin",
		OrdererOrgID:    "ordererorg",
		OrgMspID:        "org1.cvverification.com",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",
		CaID: 			 "ca.org1.cvverification.com",

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