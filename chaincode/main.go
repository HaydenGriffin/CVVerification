package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CVTrackerChaincode implementation of Chaincode
type CVTrackerChaincode struct {
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *CVTrackerChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### CVTrackerChaincode Init ###########")

	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *CVTrackerChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### CVTrackerChaincode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()
	
	// In order to manage multiple type of request, we will check the first argument.
	if function == "saveProfile" {
		return t.saveProfile(stub, args)
	} else if function == "getProfile" {
		return t.getProfile(stub, args)
	} else if function == "updateProfileCV" {
		return t.updateProfileCV(stub, args)
	} else if function == "saveCV" {
		return t.saveCV(stub, args)
	} else if function == "getCVFromCVHash" {
		return t.getCVFromCVHash(stub, args)
	} else if function == "getCVFromProfile" {
		return t.getCVFromProfile(stub, args)
	} else if function == "getCVHashFromProfile" {
		return t.getCVHashFromProfile(stub, args)
	} else if function == "saveRating" {
		return t.saveRating(stub, args)
	} else if function == "getRatings" {
		return t.getRatings(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(CVTrackerChaincode))
	if err != nil {
		fmt.Printf("Error starting CVTracker chaincode: %s", err)
	}
}