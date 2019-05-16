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

// Invoke chaincode
// All future requests named invoke will arrive here.
func (t *CVVerificationChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("CVVerificationChaincode Invoke")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error(fmt.Sprintf("Unknown function call: %v", function))
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is invalid.")
	}

	// The chaincode request operation is stored within the first argument
	operationType := args[0]

	// All query operations
	if operationType == "query" {
		return t.query(stub, args[1:])
	}

	// The update argument will manage all update in the ledger
	if operationType == "update" {
		return t.update(stub, args[1:])
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
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