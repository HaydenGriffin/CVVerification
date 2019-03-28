package main

import (
	"encoding/json"
	"fmt"
	"github.com/cvtracker/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// query function that handle every readonly in the ledger
func (t *CVTrackerChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	} else if function == "saveprofile" {
		return t.saveProfile(stub, args[1:])
	} else if function == "saveprofilecv" {
		return t.saveProfileCV(stub, args[1:])
	} else if function == "saverating" {
		return t.saveRating(stub, args[1:])
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

func (t *CVTrackerChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("Register new user")

	actorType, found, err := cid.GetAttributeValue(stub, model.ActorAttribute)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to identify the type of the request owner: %v", err))
	}
	if !found {
		return shim.Error("The type of the request owner is not present")
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
				Name: args[0],
			},
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
				Name: args[0],
			},
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
				Name: args[0],
			},
		}
		err = updateInLedger(stub, model.ObjectTypeVerifier, actorID, newVerifier)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable to register the new verifier in the ledger: %v", err))
		}
		newVerifierAsByte, err := convertObjectToByte(newVerifier)
		if err != nil {
			return shim.Error(fmt.Sprintf("Unable convert the new verifier to byte: %v", err))
		}

		fmt.Printf("Applicant:\n  ID -> %s\n  Name -> %s\n", actorID, args[0])

		return shim.Success(newVerifierAsByte)
	default:
		return shim.Error("The type of the request owner is unknown")
	}
}

func (t *CVTrackerChaincode) saveCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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
		return shim.Error("Unable to convert cv byte to object")
	}

	cvHash := args[1]
	if cvHash == "" {
		return shim.Error("The CV hash is empty.")
	}

	err = updateInLedger(stub, model.ObjectTypeCV, cvHash, cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to create the CV in the ledger: %v", err))
	}

	resourceAsByte, err := convertObjectToByte(cv)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable convert the resource to byte: %v", err))
	}

	fmt.Printf("Resource created:\n  ID -> %s\n  Description -> %s\n", model.ObjectTypeCV, cvHash)

	return shim.Success(resourceAsByte)
}

// Add CV Chaincode
// args: CV object
// CV Hash is key, CVObject is the value
func (t *CVTrackerChaincode) saveProfile(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check whether the number of arguments is sufficient
	if len(args) != 2 {
		return shim.Error("The number of arguments is invalid.")
	}

	var profile model.UserProfile

	err := json.Unmarshal([]byte(args[0]), &profile)
	if err != nil {
		return shim.Error("An error occurred whilst deserialising the object")
	}

	profileHash := args[1]

	if profileHash == "" {
		return shim.Error("The profile hash is empty.")
	}

	err = updateInLedger(stub, model.ObjectTypeProfile, profileHash, profile)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to create the CV in the ledger: %v", err))
	}

	resourceAsByte, err := convertObjectToByte(profile)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable convert the resource to byte: %v", err))
	}

	fmt.Printf("Resource created:\n  ID -> %s\n  Description -> %s\n", model.ObjectTypeCV, profileHash)

	return shim.Success(resourceAsByte)
}

// Add CV Chaincode
// args: CV object
// CV Hash is key, CVObject is the value
func (t *CVTrackerChaincode) saveProfileCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check whether the number of arguments is sufficient
	if len(args) != 2 {
		return shim.Error("The number of arguments is invalid.")
	}

	profileHash := args[0]

	if profileHash == "" {
		return shim.Error("The profile hash is empty.")
	}

	cvHash := args[1]

	if cvHash == "" {
		return shim.Error("The cv hash is empty.")
	}

	var profile model.UserProfile

	err := getFromLedger(stub, model.ObjectTypeProfile, profileHash, &profile)

	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve profile from the ledger: %v", err))
	}

	profile.CVHistory = append(profile.CVHistory, cvHash)

	err = updateInLedger(stub, model.ObjectTypeProfile, profileHash, profile)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to create the CV in the ledger: %v", err))
	}

	profileAsByte, err := convertObjectToByte(profile)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable convert the profile to byte: %v", err))
	}

	fmt.Printf("Resource updated:\n  ID -> %s\n  Description -> %s\n", model.ObjectTypeProfile, profileHash)

	return shim.Success(profileAsByte)
}

// Add CV Chaincode
// args: CV object
// CV Hash is key, CVObject is the value
func (t *CVTrackerChaincode) saveRating(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check whether the number of arguments is sufficient
	if len(args) != 3 {
		return shim.Error("The number of arguments is invalid.")
	}

	profileHash := args[0]
	cvHash := args[1]
	ratingString := args[2]
	var verifierRating model.CVReview
	var profile model.UserProfile

	if profileHash == "" {
		return shim.Error("The profile hash is empty.")
	}

	if cvHash == "" {
		return shim.Error("The cv hash is empty.")
	}

	if len(ratingString) == 0 {
		return shim.Error("There is no rating to be saved.")
	}

	err := convertByteToObject([]byte(ratingString), &verifierRating)
	if err != nil {
		fmt.Println(err)
		return shim.Error("Unable to convert rating byte to object")
	}

	verifierRating.Id, err = cid.GetID(stub)
	if err != nil {
		fmt.Println(err)
		return shim.Error("Unable to get invoking chaincode identity")
	}

	err = getFromLedger(stub, model.ObjectTypeProfile, profileHash, &profile)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve profile in the ledger: %v", err))
	}

	var reviews = make(map[string][]model.CVReview)
	existingRatingFound := false

	if profile.Reviews != nil {
		fmt.Println("profile.Reviews not nil")
		fmt.Println(profile.Reviews)
		reviews = profile.Reviews

		for i, rating := range profile.Reviews[cvHash] {

			// If the verifier has already rated the CV
			if rating.Id == verifierRating.Id {
				reviews[cvHash][i] = verifierRating
				existingRatingFound = true
				}
			}
		}

	if existingRatingFound == false {
		reviews[cvHash] = append(reviews[cvHash], verifierRating)
	}

	profile.Reviews = make(map[string][]model.CVReview)

	for cvHash, cvReview := range reviews {
		profile.Reviews[cvHash] = cvReview
	}

	fmt.Println(profile)

	// put the updated profile back to the ledger
	err = updateInLedger(stub, model.ObjectTypeProfile, profileHash, profile)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to create the CV in the ledger: %v", err))
	}

	return shim.Success([]byte("Successfully saved the rating"))
}
