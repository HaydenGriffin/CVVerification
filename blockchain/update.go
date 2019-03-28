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
func (u *User) UpdateSaveCV(cvByte []byte, cvHash string) error {
	return u.update([][]byte{[]byte("savecv"), cvByte, []byte(cvHash)}, nil)
}

// UpdateSaveProfile allow to add a resource into the blockchain
func (u *User) UpdateSaveProfile(profileByte []byte, profileHash string) error {
	return u.update([][]byte{[]byte("saveprofile"), profileByte, []byte(profileHash)}, nil)
}

// UpdateSaveProfileCV allow to add a resource into the blockchain
func (u *User) UpdateSaveProfileCV(profileHash, cvHash string) error {
	return u.update([][]byte{[]byte("saveprofilecv"), []byte(profileHash), []byte(cvHash)}, nil)
}

// UpdateSaveProfileCV allow to add a resource into the blockchain
func (u *User) UpdateSaveRating(profileHash, cvHash string, ratingByte []byte) error {
	return u.update([][]byte{[]byte("saverating"), []byte(profileHash), []byte(cvHash), ratingByte}, nil)
}
/*
// UpdateDelete allow to delete a resource into the blockchain
func (u *User) UpdateDelete(resourceID string) error {
	return u.update([][]byte{[]byte("delete"), []byte(resourceID)}, nil)
}

// UpdateAcquire allow to acquire a resource into the blockchain
func (u *User) UpdateAcquire(resourceID string, mission string) error {
	return u.update([][]byte{[]byte("acquire"), []byte(resourceID), []byte(mission)}, nil)
}

// UpdateRelease allow to release a resource into the blockchain
func (u *User) UpdateRelease(resourceID string) error {
	return u.update([][]byte{[]byte("release"), []byte(resourceID)}, nil)
}
*/