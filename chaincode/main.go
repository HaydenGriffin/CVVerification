package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CVVerificationChaincode implementation of Chaincode
type CVVerificationChaincode struct {
	testing bool
}

// Init of the chaincode
// Function is called to instantiate the chaincode
func (t *CVVerificationChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("CVVerificationChaincode Init")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error(fmt.Sprintf("Unknown function call: %v", function))
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke - All requests from Fabric SDK Go are processed through the Invoke function
// Operations are split into update and query functions
func (t *CVVerificationChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("CVVerificationChaincode Invoke")

	// Retrieve the transaction proposal function and parameters
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error(fmt.Sprintf("Unknown function call: %v", function))
	}

	// Ensure at least one argument has been provided; otherwise return error
	if len(args) < 1 {
		return shim.Error("At least one argument parameter should be specified.")
	}

	// The chaincode request operation is stored within the first argument
	operationType := args[0]

	// All query operations
	if operationType == "query" {
		return t.query(stub, args[1:])
	}

	// If the operationType is update, return the update function
	if operationType == "update" {
		return t.update(stub, args[1:])
	}

	// If the first argument is not query or update, something has gone wrong
	return shim.Error("Invalid first argument supplied")
}

func main() {
	// Start the chaincode and make it ready for futures requests
	cvvc := new(CVVerificationChaincode)
	cvvc.testing = false
	err := shim.Start(cvvc)
	if err != nil {
		fmt.Printf("Error starting cvverification chaincode: %s", err)
	}
}