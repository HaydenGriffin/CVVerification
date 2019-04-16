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
	} else if function == "transitioncv" {
		return t.transitionCV(stub, args[1:])
	} else if function == "saveprofilekey" {
		return t.saveProfileKey(stub, args[1:])
	} else if function == "saveprofilecv" {
		return t.saveProfileCV(stub, args[1:])
	} else if function == "verifiersavereview" {
		return t.verifierSaveReview(stub, args[1:])
	} else if function == "publishreviews" {
		return t.publishReviews(stub, args[1:])
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

func (t *CVVerificationChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Register new user")

	actorType, found, err := cid.GetAttributeValue(stub, model.ActorAttribute)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the account type to register: %v", err))
	}
	if !found {
		return shim.Error("The account type to register could not be found")
	}

	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	actorID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the ID of the request owner: %v", err))
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
			return shim.Error(fmt.Sprintf("unable to register the new admin in the ledger: %v", err))
		}
		var newAdminAsByte []byte
		newAdminAsByte, err = convertObjectToByte(newAdmin)
		if err != nil {
			return shim.Error(fmt.Sprintf("unable convert the new admin to byte: %v", err))
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
			return shim.Error(fmt.Sprintf("unable to register the new applicant in the ledger: %v", err))
		}
		newApplicantAsByte, err := convertObjectToByte(newApplicant)
		if err != nil {
			return shim.Error(fmt.Sprintf("unable convert the new applicant to byte: %v", err))
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
			return shim.Error(fmt.Sprintf("unable to register the new verifier in the ledger: %v", err))
		}
		newVerifierAsByte, err := convertObjectToByte(newVerifier)
		if err != nil {
			return shim.Error(fmt.Sprintf("unable convert the new verifier to byte: %v", err))
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
			return shim.Error(fmt.Sprintf("unable to register the new employer in the ledger: %v", err))
		}
		newEmployerAsByte, err := convertObjectToByte(newEmployer)
		if err != nil {
			return shim.Error(fmt.Sprintf("unable convert the new employer to byte: %v", err))
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
		return shim.Error(fmt.Sprintf("only applicant users are allowed to add a CV: %v", err))
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
		return shim.Error(fmt.Sprintf("unable to convert cv byte to object: %v", err))
	}

	cvID := args[1]
	if cvID == "" {
		return shim.Error("The CV ID is empty.")
	}

	cv.Status = model.CVInDraft

	err = updateInLedger(stub, model.ObjectTypeCV, cvID, cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to create the CV in the ledger: %v", err))
	}

	cvAsByte, err := convertObjectToByte(cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable convert the CV to byte: %v", err))
	}

	fmt.Printf("CV created:\n  ID -> %s\n  Description -> %s\n", model.ObjectTypeCV, cvID)

	return shim.Success(cvAsByte)
}

func (t *CVVerificationChaincode) transitionCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Transition CV status")

	if len(args) != 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	cvIDToUpdate := args[0]
	if cvIDToUpdate  == "" {
		return shim.Error("The CV ID is empty.")
	}

	newStatus := args[1]
	if newStatus == "" {
		return shim.Error("The status value is empty.")
	}

	actorType, found, err := cid.GetAttributeValue(stub, model.ActorAttribute)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the type of the request owner: %v", err))
	}
	if !found {
		return shim.Error("The type of the request owner is not present")
	}

	var cv model.CVObject

	err = getFromLedger(stub, model.ObjectTypeCV, cvIDToUpdate, &cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve CV in the ledger: %v", err))
	}

	err = canCVBeTransitioned(actorType, newStatus, cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to perform request on ledger: %v", err))
	}

	cv.Status = newStatus

	err = updateInLedger(stub, model.ObjectTypeCV, cvIDToUpdate, cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to create the CV in the ledger: %v", err))
	}

	cvAsByte, err := convertObjectToByte(cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable convert the CV to byte: %v", err))
	}

	fmt.Printf("CV Status Updated:\n  ID -> %s\n  New Status -> %s\n", cvIDToUpdate, newStatus)

	return shim.Success(cvAsByte)
}

// Add CV Chaincode
// args: CV object
// CV ID is key, CVObject is the value
func (t *CVVerificationChaincode) saveProfileKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Save profile key")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorApplicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("only applicant users are allowed to update their key: %v", err))
	}

	// Check whether the number of arguments is sufficient
	if len(args) != 1 {
		return shim.Error("The number of arguments is invalid.")
	}

	publicKeyByte := args[0]

	if len(publicKeyByte) == 0 {
		return shim.Error("The publicKey is empty.")
	}

	var applicant model.Applicant

	applicantID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the ID of the request owner: %v", err))
	}

	err = getFromLedger(stub, model.ObjectTypeApplicant, applicantID, &applicant)

	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve applicant profile from the ledger: %v", err))
	}

	if applicant.ID != applicantID {
		return shim.Error("Unable to update profile as applicantID differs from profile ID")
	}

	applicant.Profile.PublicKey = string(publicKeyByte)

	if len(applicant.Profile.Reviews) > 0 {
		applicant.Profile.Reviews = *new(map[string]map[string][]byte)
	}

	err = updateInLedger(stub, model.ObjectTypeApplicant, applicantID, applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to create the CV in the ledger: %v", err))
	}

	applicantAsByte, err := convertObjectToByte(applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable convert the profile to byte: %v", err))
	}

	fmt.Printf("Resource updated:\n  ID -> %s\n  Description -> %s\n", model.ObjectTypeApplicant, applicantID)

	return shim.Success(applicantAsByte)
}


// Add CV Chaincode
// args: CV object
// CV ID is key, CVObject is the value
func (t *CVVerificationChaincode) saveProfileCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Save profile CV")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorApplicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("only applicant users are allowed to add a CV: %v", err))
	}

	// Check whether the number of arguments is sufficient
	if len(args) != 1 {
		return shim.Error("The number of arguments is invalid.")
	}

	cvID := args[0]
	if cvID == "" {
		return shim.Error("The CV ID is empty.")
	}

	var applicant model.Applicant

	applicantID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to identify the ID of the request owner: %v", err))
	}

	err = getFromLedger(stub, model.ObjectTypeApplicant, applicantID, &applicant)

	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve applicant profile from the ledger: %v", err))
	}

	if applicant.ID != applicantID {
		return shim.Error("Unable to update profile as applicantID differs from profile ID")
	}

	applicant.Profile.CVHistory = append(applicant.Profile.CVHistory, cvID)

	err = updateInLedger(stub, model.ObjectTypeApplicant, applicantID, applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to create the CV in the ledger: %v", err))
	}

	applicantAsByte, err := convertObjectToByte(applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable convert the profile to byte: %v", err))
	}

	fmt.Printf("Resource updated:\n  ID -> %s\n  Description -> %s\n", model.ObjectTypeApplicant, applicantID)

	return shim.Success(applicantAsByte)
}

// Add CV Chaincode
// args: CV object
// CV ID is key, CVObject is the value
func (t *CVVerificationChaincode) verifierSaveReview(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Save review")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorVerifier)
	if err != nil {
		return shim.Error(fmt.Sprintf("only verifier users are allowed to add a CV: %v", err))
	}

	// Check whether the number of arguments is sufficient
	if len(args) != 3 {
		return shim.Error("The number of arguments is invalid.")
	}

	applicantID := args[0]
	cvID := args[1]
	encryptedReviewString := args[2]
	var applicant model.Applicant

	if applicantID == "" {
		return shim.Error("The profile hash is empty.")
	}

	if cvID == "" {
		return shim.Error("The CV ID is empty.")
	}

	if len(encryptedReviewString) == 0 {
		return shim.Error("There is no rating to be saved.")
	}

	verifierID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to get invoking chaincode identity: %v", err))
	}

	err = getFromLedger(stub, model.ObjectTypeApplicant, applicantID, &applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve applicant profile in the ledger: %v", err))
	}

	var reviews = make(map[string]map[string][]byte)

	if applicant.Profile.Reviews != nil {
		reviews = applicant.Profile.Reviews
	}

	// CV currently hasn't been reviewed
	// Initialise the map
	if reviews[cvID] == nil {
		reviews[cvID] = make(map[string][]byte)
	}

	reviews[cvID][verifierID] = []byte(encryptedReviewString)
	applicant.Profile.Reviews = reviews

	// Put the updated profile back to the ledger
	err = updateInLedger(stub, model.ObjectTypeApplicant, applicantID, applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to save the review in the ledger: %v", err))
	}

	return shim.Success([]byte("Successfully saved the review"))
}

// Add CV Chaincode
// args: CV object
// CV ID is key, CVObject is the value
func (t *CVVerificationChaincode) publishReviews(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Save rating")

	err := cid.AssertAttributeValue(stub, model.ActorAttribute, model.ActorApplicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("only verifier users are allowed to add a CV: %v", err))
	}

	// Check whether the number of arguments is sufficient
	if len(args) != 2 {
		return shim.Error("The number of arguments is invalid.")
	}

	cvID := args[0]
	reviewsByte := args[1]
	var reviews []model.CVReview

	if cvID == "" {
		return shim.Error("The CV ID is empty.")
	}

	err = convertByteToObject([]byte(reviewsByte), &reviews)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to convert reviews byte to object: %v", err))
	}

	var applicant model.Applicant

	applicantID, err := cid.GetID(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to get invoking chaincode identity: %v", err))
	}

	err = getFromLedger(stub, model.ObjectTypeApplicant, applicantID, &applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to retrieve applicant profile in the ledger: %v", err))
	}

	var reviewsMap = make(map[string][]model.CVReview)

	if applicant.Profile.PublicReviews != nil {
		reviewsMap = applicant.Profile.PublicReviews
	}

	reviewsMap[cvID] = reviews

	applicant.Profile.PublicReviews = reviewsMap

	// Put the updated profile back to the ledger
	err = updateInLedger(stub, model.ObjectTypeApplicant, applicantID, applicant)
	if err != nil {
		return shim.Error(fmt.Sprintf("unable to save the review in the ledger: %v", err))
	}

	return shim.Success([]byte("Successfully saved the profile"))
}
