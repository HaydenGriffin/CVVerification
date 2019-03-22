package main

import (
	"fmt"
	"github.com/cvtracker/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// query function that handle every readonly in the ledger
func (t *CVTrackerChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Query functions")

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is invalid.")
	}

	if args[0] == "admin" {
		return t.admin(stub, args[1:])
	} else if args[0] == "applicant" {
		return t.applicant(stub, args[1:])
	} else if args[0] == "cv" {
		return t.cv(stub, args[1:])
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

func (t *CVTrackerChaincode) admin(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("# admin information")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorAdmin)
	if err != nil {
		return shim.Error(fmt.Sprintf("Only admin is allowed for the kind of request: %v", err))
	}

	adminID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to identify the ID of the request owner: %v", err))
	}
	var admin model.Admin
	err = getFromLedger(stub, model.ObjectTypeAdmin, adminID, &admin)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve admin in the ledger: %v", err))
	}
	adminAsByte, err := convertObjectToByte(admin)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable convert the admin to byte: %v", err))
	}

	return shim.Success(adminAsByte)
}

func (t *CVTrackerChaincode) applicant(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("# applicant information")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorApplicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("Only applicant is allowed for the kind of request: %v", err))
	}

	applicantID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to identify the ID of the request owner: %v", err))
	}
	var applicant model.Applicant
	err = getFromLedger(stub, model.ObjectTypeApplicant, applicantID, &applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve applicant in the ledger: %v", err))
	}
	applicantAsByte, err := convertObjectToByte(applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable convert the admin to byte: %v", err))
	}

	return shim.Success(applicantAsByte)
}

func (t *CVTrackerChaincode) cv(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("# cv detail")

	if len(args) != 1 {
		return shim.Error("The number of arguments is invalid.")
	}

	cvHash := args[0]
	if cvHash == "" {
		return shim.Error("The cvHash is empty.")
	}

	var cv CVObject

	err := getFromLedger(stub, model.ObjectTypeCV, cvHash, &cv)

	cvAsByte, err := convertObjectToByte(cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to convert the resource histories to byte: %v", err))
	}

	return shim.Success(cvAsByte)
}