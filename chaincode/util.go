package main

import (
	"encoding/json"
	"fmt"
	"github.com/cvverification/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strings"
)

// convertObjectToByte
func convertObjectToByte(object interface{}) ([]byte, error) {
	byteArray, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("error occurred whilst marshalling object to byte array: %v", err)
	}
	return byteArray, nil
}

// convertByteToObject
func convertByteToObject(byteArray []byte, result interface{}) error {
	err := json.Unmarshal(byteArray, result)
	if err != nil {
		return fmt.Errorf("error occurred whilst unmarshalling byte array to object: %v", err)
	}
	return nil
}

// getFromLedger retrieve an object from the ledger
func getFromLedger(stub shim.ChaincodeStubInterface, objectType string, id string, result interface{}) error {
	key, err := stub.CreateCompositeKey(objectType, []string{id})
	if err != nil {
		return fmt.Errorf("unable to create the object key for the ledger: %v", err)
	}
	resultAsByte, err := stub.GetState(key)
	if err != nil {
		return fmt.Errorf("unable to retrieve the object in the ledger: %v", err)
	}
	if resultAsByte == nil {
		return fmt.Errorf("the object doesn't exist in the ledger")
	}
	err = convertByteToObject(resultAsByte, result)
	if err != nil {
		return fmt.Errorf("unable to convert the result to object: %v", err)
	}
	return nil
}

// updateInLedger update an object in the ledger
func updateInLedger(stub shim.ChaincodeStubInterface, objectType string, id string, object interface{}) error {
	key, err := stub.CreateCompositeKey(objectType, []string{id})
	if err != nil {
		return fmt.Errorf("unable to create the object key for the ledger: %v", err)
	}

	objectAsByte, err := convertObjectToByte(object)
	if err != nil {
		return err
	}
	err = stub.PutState(key, objectAsByte)
	if err != nil {
		return fmt.Errorf("unable to put the object in the ledger: %v", err)
	}
	return nil
}

// canCVBeTransitioned contains logic to ensure that a CV can be transitioned
func canCVBeTransitioned(actorType, transitionTo string, cv model.CVObject) error {
	// Can't transition case to status it currently is
	if transitionTo == cv.Status {
		return fmt.Errorf("unable to transition to empty status")
	}

	switch transitionTo {
	case model.CVInDraft:
		if cv.Status == model.CVInReview && actorType == model.ActorApplicant {
			return nil
		}
	case model.CVInReview:
		if (cv.Status == model.CVInDraft || cv.Status == model.CVSubmitted) && actorType == model.ActorApplicant {
			return nil
		}
	case model.CVSubmitted:
		if (cv.Status == model.CVInDraft || cv.Status == model.CVInReview) && actorType == model.ActorApplicant {
			return nil
		}
	default:
		return fmt.Errorf("unable to transition CV object from: %v to: %v", cv.Status, transitionTo)
	}

	return fmt.Errorf("unable to transition CV object from: %v to: %v", cv.Status, transitionTo)
}

// returnCV checks whether the CV is in the correct state for the user calling the function
func returnCV(actorType, filter string, cv model.CVObject) bool {

	// If the user specifies a filter - check this first
	if filter != "" {
		filter = strings.ToLower(filter)
		// Check to see if the industry contains the filter
		if !strings.Contains(strings.ToLower(cv.Industry), filter) {
			return false
		}
	}

	switch actorType {
	// Verifier users are able to view all CVs that are in review
	case model.ActorVerifier:
		if cv.Status == model.CVInReview {
			return true
		}
	case model.ActorEmployer:
		if cv.Status == model.CVSubmitted {
			return true
		}
	default:
		return false
	}
	return false
}
