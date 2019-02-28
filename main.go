/**
  author: Hayden Griffin
 */

package main

import (
	"fmt"
	"github.com/cvtracker/blockchain"
	"github.com/cvtracker/controllers"
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

	// Launch the web application listening
	app := &controllers.Application{
		Service: &serviceSetup,
	}
	web.Serve(app)
}