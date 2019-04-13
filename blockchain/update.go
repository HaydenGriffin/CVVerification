package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

// update internal method that allow a user to invoke on the blockchain chaincode
func (u *User) update(args [][]byte, responseObject interface{}) error {

	response, err := u.ChannelClient.Execute(
		channel.Request{ChaincodeID: u.Fabric.ChaincodeID, Fcn: "invoke", Args: append([][]byte{[]byte("update")}, args...)},
		channel.WithRetry(retry.DefaultChannelOpts),
	)
	if err != nil {
		return fmt.Errorf("unable to perform the update: %v", err)
	}

	if responseObject != nil {
		err = json.Unmarshal(response.Payload, responseObject)
		if err != nil {
			return fmt.Errorf("unable to convert response to the object given for the update: %v", err)
		}
	}

	return nil
}

// UpdateRegister allow to register a user into the blockchain
func (u *User) UpdateRegister() error {
	return u.update([][]byte{[]byte("register"), []byte(u.Username)}, nil)
}

// UpdateSaveCV allow to add a resource into the blockchain
func (u *User) UpdateSaveCV(cvByte []byte, cvID string) error {
	return u.update([][]byte{[]byte("savecv"), cvByte, []byte(cvID)}, nil)
}

// UpdateSaveCV allow to add a resource into the blockchain
func (u *User) UpdateTransitionCV(cvID, newStatus string) error {
	return u.update([][]byte{[]byte("transitioncv"), []byte(cvID), []byte(newStatus)}, nil)
}

// UpdateSaveProfileCV allow to add a resource into the blockchain
func (u *User) UpdateSaveProfileKey(publicKey string) error {
	return u.update([][]byte{[]byte("saveprofilekey"), []byte(publicKey)}, nil)
}

// UpdateSaveProfileCV allow to add a resource into the blockchain
func (u *User) UpdateSaveProfileCV(cvID string) error {
	return u.update([][]byte{[]byte("saveprofilecv"), []byte(cvID)}, nil)
}

// UpdateSaveProfileCV allow to add a resource into the blockchain
func (u *User) UpdateSaveRating(ID, cvID string, ratingByte []byte) error {
	return u.update([][]byte{[]byte("saverating"), []byte(ID), []byte(cvID), ratingByte}, nil)
}
