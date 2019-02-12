/**
  author: Hayden Griffin
 */

package main

import (
	"fmt"
	"github.com/cvverification/blockchain"
	"github.com/cvverification/web"
	"github.com/cvverification/controllers"
	"os"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters 
		OrdererID: "orderer.cvverification.com",

		// Channel parameters
		ChannelID:     "cvverification",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/cvverification/fixtures/artifacts/cvverification.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "cvverification",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/cvverification/chaincode/",
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