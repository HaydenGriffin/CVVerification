package main

import (
	"fmt"
	"github.com/cvverification/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
)

// query function that handle every readonly in the ledger
func (t *CVVerificationChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Query functions")

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is invalid.")
	}

	if args[0] == "id" {
		return t.id(stub, args[1:])
	} else if args[0] == "admin" {
		return t.admin(stub, args[1:])
	} else if args[0] == "applicant" {
		return t.applicant(stub, args[1:])
	} else if args[0] == "verifier" {
		return t.verifier(stub, args[1:])
	} else if args[0] == "cv" {
		return t.cv(stub, args[1:])
	} else if args[0] == "cvs" {
		return t.cvs(stub, args[1:])
	} else if args[0] == "verifiercvreview" {
		return t.verifiercvreview(stub, args[1:])
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

func (t *CVVerificationChaincode) id(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("# id information")

	var actor model.Actor

	actorID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the ID of the request owner: %v", err))
	}

	actor.ID = actorID

	actorAsByte, err := convertObjectToByte(actor)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable convert the actor to byte: %v", err))
	}

	return shim.Success(actorAsByte)
}

func (t *CVVerificationChaincode) admin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("# admin information")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorAdmin)
	if err != nil {
		return shim.Error(fmt.Sprintf("only admin is allowed for the kind of request: %v", err))
	}

	adminID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the ID of the request owner: %v", err))
	}

	var admin model.Admin

	err = getFromLedger(stub, model.ObjectTypeAdmin, adminID, &admin)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve admin in the ledger: %v", err))
	}

	adminAsByte, err := convertObjectToByte(admin)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable convert the admin to byte: %v", err))
	}

	return shim.Success(adminAsByte)
}

func (t *CVVerificationChaincode) applicant(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("# applicant information")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorApplicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("only applicant is allowed for the kind of request: %v", err))
	}

	applicantID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the ID of the request owner: %v", err))
	}

	var applicant model.Applicant

	err = getFromLedger(stub, model.ObjectTypeApplicant, applicantID, &applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve applicant in the ledger: %v", err))
	}

	applicantAsByte, err := convertObjectToByte(applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable convert the admin to byte: %v", err))
	}

	return shim.Success(applicantAsByte)
}

func (t *CVVerificationChaincode) verifier(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("# verifier information")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorVerifier)
	if err != nil {
		return shim.Error(fmt.Sprintf("only verifier is allowed for the kind of request: %v", err))
	}

	verifierID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the ID of the request owner: %v", err))
	}

	var verifier model.Verifier
	err = getFromLedger(stub, model.ObjectTypeVerifier, verifierID, &verifier)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve verifier in the ledger: %v", err))
	}

	verifierAsByte, err := convertObjectToByte(verifier)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable convert the verifier to byte: %v", err))
	}

	return shim.Success(verifierAsByte)
}

func (t *CVVerificationChaincode) cv(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("# cv detail")

	if len(args) != 1 {
		return shim.Error("The number of arguments is invalid.")
	}

	cvID := args[0]
	if cvID == "" {
		return shim.Error("The cvID is empty.")
	}

	var cv model.CVObject

	err := getFromLedger(stub, model.ObjectTypeCV, cvID, &cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve CV in the ledger: %v", err))
	}

	cvAsByte, err := convertObjectToByte(cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to convert the resource histories to byte: %v", err))
	}

	return shim.Success(cvAsByte)
}

func (t *CVVerificationChaincode) cvs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("# cv list detail")

	if len(args) != 2 {
		return shim.Error("The number of arguments is invalid.")
	}

	status := args[0]
	if status == "" {
		return shim.Error("The status is empty.")
	}

	filter := args[1]

	iterator, err := stub.GetStateByPartialCompositeKey(model.ObjectTypeCV, []string{})
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve the list of resource in the ledger: %v", err))
	}

	actorType, found, err := cid.GetAttributeValue(stub, model.ActorAttribute)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the type of the request owner: %v", err))
	}
	if !found {
		return shim.Error("The type of the request owner is not present")
	}

	cvList := make(map[string] model.CVObject)
	var cv model.CVObject
	for iterator.HasNext() {
		cvKeyValue, err := iterator.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("unable to retrieve CV from ledger: %v", err))
		}
		err = convertByteToObject(cvKeyValue.Value, &cv)
		if err != nil {
			return shim.Error(fmt.Sprintf("unable to convert CV byte to object: %v", err))
		}
		// The key value contains null characters for some reason, and it has 'cv' appended to the front
		id := cvKeyValue.Key
		// Remove null characters
		id = strings.Replace(id, "\x00","",-1)
		// Remove 'cv' from front of string
		id = id[2:]
		if returnCV(actorType, filter, cv) {
			fmt.Println("RETURNED CV:")
			fmt.Println(cv)
			cvList[id] = cv
		}
	}

	cvListByte, err := convertObjectToByte(cvList)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to convert CV list to byte: %v", err))
	}

	return shim.Success(cvListByte)
}

func (t *CVVerificationChaincode) verifiercvreview(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("# verifier cv review")

	if len(args) != 2 {
		return shim.Error("The number of arguments is invalid.")
	}

	var applicant model.Applicant
	verifierReview := model.CVReview{}
	applicantID := args[0]
	cvID := args[1]

	if applicantID == "" {
		return shim.Error("The applicant ID is empty.")
	}

	if cvID == "" {
		return shim.Error("The CV ID is empty.")
	}

	verifierID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error("Unable to retrieve user identity.")
	}

	err = getFromLedger(stub, model.ObjectTypeApplicant, applicantID, &applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve applicant profile in the ledger: %v", err))
	}

	if applicant.Profile.Reviews[cvID] != nil {
		for _, review := range applicant.Profile.Reviews[cvID] {
			if review.VerifierID == verifierID {
				verifierReview = review
			}
		}
	}

	reviewAsByte, err := convertObjectToByte(verifierReview)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to convert the review to byte: %v", err))
	}

	return shim.Success(reviewAsByte)
}
