/**
  @Author : Hayden Griffin
*/

package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) SaveProfile(profile UserProfile, userHash string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "saveProfile")
	args = append(args, userHash)

	eventID := "eventSaveProfile"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(profile)
	if err != nil {
		return "", fmt.Errorf("an error occurred whilst serialising the user profile")
	}

	// Create a request (proposal) and send it
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{b, []byte(args[1]), []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceSetup) GetProfile(profile string) ([]byte, error){

	// Prepare arguments
	var args []string
	args = append(args, "getProfile")
	args = append(args, profile)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(profile)}}

	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}

func (t *ServiceSetup) UpdateProfileCV(profileHash, cvHash string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "updateProfileCV")
	args = append(args, profileHash)
	args = append(args, cvHash)

	eventID := "eventUpdateProfile"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// Create a request (proposal) and send it
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

func (t *ServiceSetup) SaveCV(cv CVObject, cvHash string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "saveCV")
	args = append(args, cvHash)

	eventID := "eventSaveCV"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(cv)
	if err != nil {
		return "", fmt.Errorf("an error occurred whilst serialising the cv object")
	}

	// Create a request (proposal) and send it
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{b, []byte(cvHash), []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

// QueryHello query the chaincode to get the state of hello
func (t *ServiceSetup) GetCVFromCVHash(cvHash string) ([]byte, error) {

	// Prepare arguments
	var args []string
	args = append(args, "getCVFromCVHash")
	args = append(args, cvHash)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}}

	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}

// QueryHello query the chaincode to get the state of hello
func (t *ServiceSetup) GetCVFromProfile(profileHash string) ([]byte, error) {

	// Prepare arguments
	var args []string
	args = append(args, "getCVFromProfile")
	args = append(args, profileHash)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}}

	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}

// QueryHello query the chaincode to get the state of hello
func (t *ServiceSetup) GetCVHashFromProfile(profileHash string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "getCVHashFromProfile")
	args = append(args, profileHash)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}}

	response, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(response.Payload), nil
}


func (t *ServiceSetup) SaveRating(profileHash, cvHash string, rating CVRating) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "saveRating")
	args = append(args, profileHash)
	args = append(args, cvHash)

	eventID := "eventSaveRating"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(rating)
	if err != nil {
		return "", fmt.Errorf("an error occurred whilst serialising the rating")
	}

	// Create a request (proposal) and send it
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{b, []byte(args[1]), []byte(args[2]), []byte(eventID)}}
	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

// QueryHello query the chaincode to get the state of hello
func (t *ServiceSetup) GetRatings(profileHash, cvHash string) ([]byte, error) {

	// Prepare arguments
	var args []string
	args = append(args, "getRatings")
	args = append(args, profileHash)
	args = append(args, cvHash)


	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}}

	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}