package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const DOC_TYPE = "cvObj"

// query
// Every readonly functions in the ledger will be here
func (t *CVTrackerChaincode) queryCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### CVTrackerChaincode query ###########")

	// Check whether the number of arguments is sufficient
	if len(args) != 1 {
		return shim.Error("The number of arguments is incorrect for the method.")
	}

	// Get the state of the value matching the key hello in the ledger
	//state, err := stub.GetState("cv")
	state, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get state of cv")
	}

	if state == nil {
		return shim.Error("No CV exists for the specified key")
	}

	// Return this value in response
	return shim.Success(state)

}

// Save cv
// args: cv
func PutCV(stub shim.ChaincodeStubInterface, cv CVObject) ([]byte, bool) {

	cv.ObjectType = DOC_TYPE

	b, err := json.Marshal(cv)
	if err != nil {
		return nil, false
	}

	// Save resume status
	//err = stub.PutState("cv", b)
	err = stub.PutState(cv.CVHash, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// Add Resume Chaincode
// args: resume object
// Resume Hash is key, Resume is the value
func (t *CVTrackerChaincode) addCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var cv CVObject
	err := json.Unmarshal([]byte(args[0]), &cv)
	if err != nil {
		return shim.Error("An error occurred whilst deserialising the object")
	}

	_, bl := PutCV(stub, cv)
	if !bl {
		return shim.Error("An error occurred whilst saving the resume")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("Successfully saved the resume"))
}