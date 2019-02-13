/**
  author: Hayden Griffin
 */

package main

import (
	"fmt"
	"github.com/cvtracker/blockchain"
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
		ChainCodeID:     "cvtracker",
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
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

	// Launch the web application listening
	app := &controllers.Application{
		Fabric: &fSetup,
	}
	web.Serve(app)
}