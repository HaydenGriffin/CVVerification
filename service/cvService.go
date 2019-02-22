/**
  @Author : Hayden Griffin
*/

package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// QueryHello query the chaincode to get the state of hello
func (t *ServiceSetup) QueryCVByHash(cvHash string) ([]byte, error) {

	// Prepare arguments
	var args []string
	args = append(args, "queryCV")
	args = append(args, cvHash)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}}

	response, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return response.Payload, nil
}

func (t *ServiceSetup) SaveCV(cv CVObject) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "addCV")

	eventID := "eventAddCV"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(cv)
	if err != nil {
		return "", fmt.Errorf("an error occurred whilst serialising the cv object")
	}

	// Create a request (proposal) and send it
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{b, []byte(eventID)}}
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

func (t *ServiceSetup) ModifyCV(cv CVObject) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "updateCV")

	eventID := "eventModifyEdu"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(cv)
	if err != nil {
		return "", fmt.Errorf("an error occurred whilst serialising the cv object")
	}

	// Create a request (proposal) and send it
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{b, []byte(eventID)}}
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

func (t *ServiceSetup) RateCV(cvHash string, rating CVRating) (string, error) {

	cv, err := t.QueryCVByHash(cvHash)

	if err != nil {
		fmt.Printf(err.Error())
	} else {

	}

	// Prepare arguments
	var args []string
	args = append(args, "addCV")

	eventID := "eventAddCV"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(cv)
	if err != nil {
		return "", fmt.Errorf("an error occurred whilst serialising the cv object")
	}

	// Create a request (proposal) and send it
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{b, []byte(eventID)}}
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