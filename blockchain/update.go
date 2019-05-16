package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/cvverification/chaincode/model"
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

// UpdateRegister allow a user to register into the blockchain
func (u *User) UpdateRegister() error {
	return u.update([][]byte{[]byte("register"), []byte(u.Username)}, nil)
}

// UpdateSaveCV save a CV to the ledger
func (u *User) UpdateSaveCV(cvByte []byte, cvID string) error {
	return u.update([][]byte{[]byte("savecv"), cvByte, []byte(cvID)}, nil)
}

// UpdateTransitionCV transition the status of a CV application
func (u *User) UpdateTransitionCV(cvID, newStatus string) (*model.CVObject,error) {
	var cv *model.CVObject
	 err := u.update([][]byte{[]byte("transitioncv"), []byte(cvID), []byte(newStatus)}, &cv)
	 if err != nil {
	 	return nil, err
	 }
	 return cv, nil
}

// UpdateSaveProfileKey updates the applicant profile and saves the public key
func (u *User) UpdateSaveProfileKey(publicKey string) error {
	return u.update([][]byte{[]byte("saveprofilekey"), []byte(publicKey)}, nil)
}

// UpdateSaveProfileCV updates the applicant profile and saves the cvID
func (u *User) UpdateSaveProfileCV(cvID string) error {
	return u.update([][]byte{[]byte("saveprofilecv"), []byte(cvID)}, nil)
}

// UpdateVerifierSaveReview updates the applicant profile with the encrypted review (using applicants pub key)
func (u *User) UpdateVerifierSaveReview(applicantID, cvID string, reviewByte []byte) error {
	return u.update([][]byte{[]byte("verifiersavereview"), []byte(applicantID), []byte(cvID), reviewByte}, nil)
}

// UpdateVerifierSaveOrganisation updates the verifier profile with the new organisation name
func (u *User) UpdateVerifierSaveOrganisation(newOrganisation string) error {
	return u.update([][]byte{[]byte("verifiersaveorganisation"), []byte(newOrganisation)}, nil)
}

// UpdatePublishReviews saves the decrypted user reviews to the users profile
func (u *User) UpdatePublishReviews(cvID string, reviewsByte []byte) error {
	return u.update([][]byte{[]byte("publishreviews"), []byte(cvID), reviewsByte}, nil)
}

// UpdateEmployerSaveCV updates the employers profile with the CV application that they are interested in
func (u *User) UpdateEmployerSaveCV(cvID string) error {
	return u.update([][]byte{[]byte("employersavecv"), []byte(cvID)}, nil)
}