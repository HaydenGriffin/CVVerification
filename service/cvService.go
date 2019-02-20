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
func (t *ServiceSetup) QueryCVByHash(cvHash string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "queryCV")
	args = append(args, cvHash)

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}}

	response, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(response.Payload), nil
}

// InvokeHello
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