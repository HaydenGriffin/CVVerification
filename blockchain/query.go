package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/cvverification/chaincode/model"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

// query internal method that allow to make query to the blockchain chaincode
func (u *User) query(args [][]byte, responseObject interface{}) error {

	response, err := u.ChannelClient.Query(
		channel.Request{ChaincodeID: u.Fabric.ChaincodeID, Fcn: "invoke", Args: append([][]byte{[]byte("query")}, args...)},
		channel.WithRetry(retry.DefaultChannelOpts),
	)

	if err != nil {
		return fmt.Errorf("unable to perform the query: %v", err)
	}

	if responseObject != nil {
		err = json.Unmarshal(response.Payload, responseObject)
		if err != nil {
			return fmt.Errorf("unable to convert response to the object given for the query: %v", err)
		}
	}

	return nil
}

// QueryAdmin query the blockchain chaincode to retrieve information about the current admin user connected
func (u *User) QueryID() (string, error) {
	var actor *model.Actor
	err := u.query([][]byte{[]byte("id")}, &actor)
	if err != nil {
		return "", err
	}
	return actor.ID, nil
}

// QueryAdmin query the blockchain chaincode to retrieve information about the current admin user connected
func (u *User) QueryAdmin() (*model.Admin, error) {
	var admin *model.Admin
	err := u.query([][]byte{[]byte("admin")}, &admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// QueryApplicant query the blockchain chaincode to retrieve information about the current applicant user connected
func (u *User) QueryApplicant() (*model.Applicant, error) {
	var applicant *model.Applicant
	err := u.query([][]byte{[]byte("applicant")}, &applicant)
	if err != nil {
		return nil, err
	}
	return applicant, nil
}

// QueryVerifier query the blockchain chaincode to retrieve information about the current applicant user connected
func (u *User) QueryVerifier() (*model.Verifier, error) {
	var verifier *model.Verifier
	err := u.query([][]byte{[]byte("verifier")}, &verifier)
	if err != nil {
		return nil, err
	}
	return verifier, nil
}

// QueryCV query the blockchain chaincode to retrieve information about the current applicant user connected
func (u *User) QueryCV(cvID string) (*model.CVObject, error) {
	var cv *model.CVObject
	err := u.query([][]byte{[]byte("cv"), []byte(cvID)}, &cv)
	if err != nil {
		return nil, err
	}
	return cv, nil
}

// QueryCV query the blockchain chaincode to retrieve information about the current applicant user connected
func (u *User) QueryCVs(status, filter string) (map[string]model.CVObject, error) {
	cvList := make(map[string]model.CVObject)
	err := u.query([][]byte{[]byte("cvs"), []byte(status), []byte(filter)}, &cvList)
	if err != nil {
		return nil, err
	}
	return cvList, nil
}

// QueryCV query the blockchain chaincode to retrieve information about the current applicant user connected
func (u *User) QueryCVReviews(ID, cvID string) ([]model.CVReview, error) {
	var ratings []model.CVReview
	err := u.query([][]byte{[]byte("cvreviews"), []byte(ID), []byte(cvID)}, &ratings)
	if err != nil {
		return nil, err
	}
	return ratings, nil
}

// QueryCV query the blockchain chaincode to retrieve information about the current applicant user connected
func (u *User) QueryVerifierCVReview(ID, cvID string) (model.CVReview, error) {
	var rating model.CVReview
	err := u.query([][]byte{[]byte("verifiercvreview"), []byte(ID), []byte(cvID)}, &rating)
	if err != nil {
		return model.CVReview{}, err
	}
	return rating, nil
}