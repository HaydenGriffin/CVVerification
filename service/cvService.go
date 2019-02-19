/**
  @Author : Hayden Griffin
*/

package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// QueryHello query the chaincode to get the state of hello
func (t *ServiceSetup) QueryHello() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "hello")

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}}

	response, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(response.Payload), nil
}


// InvokeHello
func (t *ServiceSetup) InvokeHello(value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "invoke")
	args = append(args, "hello")
	args = append(args, value)

	eventID := "eventInvoke"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")


	// Create a request (proposal) and send it
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}}
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