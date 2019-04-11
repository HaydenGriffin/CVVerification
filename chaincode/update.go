package main

import (
	"fmt"
	"github.com/cvverification/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// query function that handle every readonly in the ledger
func (t *CVVerificationChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Update functions")

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is invalid.")
	}

	function := args[0]

	if function == "register" {
		return t.register(stub, args[1:])
	} else if function == "savecv" {
		return t.saveCV(stub, args[1:])
	} else if function == "saveprofilecv" {
		return t.saveProfileCV(stub, args[1:])
	} else if function == "saverating" {
		return t.saveRating(stub, args[1:])
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

func (t *CVVerificationChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Register new user")

	actorType, found, err := cid.GetAttributeValue(stub, model.ActorAttribute)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to identify the account type to register: %v", err))
	}
	if !found {
		return shim.Error("The account type to register could not be found")
	}

	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	actorID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to identify the ID of the request owner: %v", err))
	}

	switch actorType {
	case model.ActorAdmin:
		newAdmin := model.Admin{
			Actor: model.Actor{
				ID:   actorID,
				Username: args[0],
			},
			Profile: model.AdminProfile{},
		}
		err = updateInLedger(stub, model.ObjectTypeAdmin, actorID, newAdmin)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable to register the new admin in the ledger: %v", err))
		}
		var newAdminAsByte []byte
		newAdminAsByte, err = convertObjectToByte(newAdmin)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable convert the new admin to byte: %v", err))
		}

		fmt.Printf("Admin:\n  ID -> %s\n  Name -> %s\n", actorID, args[0])

		return shim.Success(newAdminAsByte)
	case model.ActorApplicant:
		newApplicant := model.Applicant{
			Actor: model.Actor{
				ID:   actorID,
				Username: args[0],
			},
			Profile: model.ApplicantProfile{},
		}
		err = updateInLedger(stub, model.ObjectTypeApplicant, actorID, newApplicant)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable to register the new applicant in the ledger: %v", err))
		}
		newApplicantAsByte, err := convertObjectToByte(newApplicant)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable convert the new applicant to byte: %v", err))
		}

		fmt.Printf("Applicant:\n  ID -> %s\n  Name -> %s\n", actorID, args[0])

		return shim.Success(newApplicantAsByte)
	case model.ActorVerifier:
		newVerifier := model.Verifier{
			Actor: model.Actor{
				ID:   actorID,
				Username: args[0],
			},
			Profile: model.VerifierProfile{},
		}
		err = updateInLedger(stub, model.ObjectTypeVerifier, actorID, newVerifier)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable to register the new verifier in the ledger: %v", err))
		}
		newVerifierAsByte, err := convertObjectToByte(newVerifier)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable convert the new verifier to byte: %v", err))
		}

		fmt.Printf("Verifier:\n  ID -> %s\n  Name -> %s\n", actorID, args[0])

		return shim.Success(newVerifierAsByte)
	case model.ActorEmployer:
		newEmployer := model.Employer{
			Actor: model.Actor{
				ID:   actorID,
				Username: args[0],
			},
			Profile: model.EmployerProfile{},
		}
		err = updateInLedger(stub, model.ObjectTypeEmployer, actorID, newEmployer)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable to register the new employer in the ledger: %v", err))
		}
		newEmployerAsByte, err := convertObjectToByte(newEmployer)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable convert the new employer to byte: %v", err))
		}

		fmt.Printf("Employer:\n  ID -> %s\n  Name -> %s\n", actorID, args[0])

		return shim.Success(newEmployerAsByte)
	default:
		return shim.Error("The type of the request owner is unknown")
	}
}

func (t *CVVerificationChaincode) saveCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Add new CV")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorApplicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("Only applicant users are allowed to add a CV: %v", err))
	}

	if len(args) != 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	var cv model.CVObject

	cvByte := args[0]
	if len(cvByte) == 0 {
		return shim.Error("There is no CV to be saved.")
	}

	err = convertByteToObject([]byte(cvByte), &cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to convert cv byte to object: %v", err))
	}

	cvHash := args[1]
	if cvHash == "" {
		return shim.Error("The CV hash is empty.")
	}

	err = updateInLedger(stub, model.ObjectTypeCV, cvHash, cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to create the CV in the ledger: %v", err))
	}

	cvAsByte, err := convertObjectToByte(cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable convert the CV to byte: %v", err))
	}

	fmt.Printf("CV created:\n  ID -> %s\n  Description -> %s\n", model.ObjectTypeCV, cvHash)

	return shim.Success(cvAsByte)
}

// Add CV Chaincode
// args: CV object
// CV Hash is key, CVObject is the value
func (t *CVVerificationChaincode) saveProfileCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Save profile CV")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorApplicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("Only applicant users are allowed to add a CV: %v", err))
	}

	// Check whether the number of arguments is sufficient
	if len(args) != 1 {
		return shim.Error("The number of arguments is invalid.")
	}

	cvHash := args[0]

	if cvHash == "" {
		return shim.Error("The cv hash is empty.")
	}

	var applicant model.Applicant

	applicantID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to identify the ID of the request owner: %v", err))
	}

	err = getFromLedger(stub, model.ObjectTypeApplicant, applicantID, &applicant)

	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve applicant profile from the ledger: %v", err))
	}

	if applicant.ID != applicantID {
		return shim.Error("Unable to update profile as applicantID differs from profile ID")
	}

	applicant.Profile.CVHistory = append(applicant.Profile.CVHistory, cvHash)

	err = updateInLedger(stub, model.ObjectTypeApplicant, applicantID, applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to create the CV in the ledger: %v", err))
	}

	applicantAsByte, err := convertObjectToByte(applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable convert the profile to byte: %v", err))
	}

	fmt.Printf("Resource updated:\n  ID -> %s\n  Description -> %s\n", model.ObjectTypeApplicant, applicantID)

	return shim.Success(applicantAsByte)
}

// Add CV Chaincode
// args: CV object
// CV Hash is key, CVObject is the value
func (t *CVVerificationChaincode) saveRating(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Save rating")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorVerifier)
	if err != nil {
		return shim.Error(fmt.Sprintf("Only verifier users are allowed to add a CV: %v", err))
	}

	// Check whether the number of arguments is sufficient
	if len(args) != 3 {
		return shim.Error("The number of arguments is invalid.")
	}

	applicantID := args[0]
	cvHash := args[1]
	ratingString := args[2]
	var verifierRating model.CVReview
	var applicant model.Applicant

	if applicantID == "" {
		return shim.Error("The profile hash is empty.")
	}

	if cvHash == "" {
		return shim.Error("The cv hash is empty.")
	}

	if len(ratingString) == 0 {
		return shim.Error("There is no rating to be saved.")
	}

	err = convertByteToObject([]byte(ratingString), &verifierRating)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to convert rating byte to object: %v", err))
	}

	verifierRating.VerifierID, err = cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to get invoking chaincode identity: %v", err))
	}

	err = getFromLedger(stub, model.ObjectTypeApplicant, applicantID, &applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve applicant profile in the ledger: %v", err))
	}

	var reviews = make(map[string][]model.CVReview)
	existingRatingFound := false

	if applicant.Profile.Reviews != nil {
		reviews = applicant.Profile.Reviews

		for i, rating := range applicant.Profile.Reviews[cvHash] {
			// If the verifier has already rated the CV
			if rating.VerifierID == verifierRating.VerifierID {
				reviews[cvHash][i] = verifierRating
				existingRatingFound = true
			}
		}
	}

	// No existing rating from verifier - append new review
	if existingRatingFound == false {
		reviews[cvHash] = append(reviews[cvHash], verifierRating)
	}

	applicant.Profile.Reviews = make(map[string][]model.CVReview)

	for cvHash, cvReview := range reviews {
		applicant.Profile.Reviews[cvHash] = cvReview
	}

	// Put the updated profile back to the ledger
	err = updateInLedger(stub, model.ObjectTypeApplicant, applicantID, applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to save the review in the ledger: %v", err))
	}

	return shim.Success([]byte("Successfully saved the rating"))
}
