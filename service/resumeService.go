/**
  @Author : hanxiaodong
*/

package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"encoding/json"
	"fmt"
)

func (t *ServiceSetup) SaveCV(cv CV) (string, error) {

	eventID := "eventAddCV"
	reg, notifier := registerEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// Serialize the cv object into a byte array
	b, err := json.Marshal(cv)
	if err != nil {
		return "", fmt.Errorf("An error occurred whilst serialising the object")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "addCV", Args: [][]byte{b, []byte(eventID)}}
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