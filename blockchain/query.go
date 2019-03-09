// Copyright 2018 Antoine CHABERT, toHero.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/cvtracker/chaincode/model"
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
func (u *User) QueryAdmin() (*model.Admin, error) {
	var admin *model.Admin
	err := u.query([][]byte{[]byte("admin")}, &admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// QueryConsumer query the blockchain chaincode to retrieve information about the current consumer user connected
func (u *User) QueryConsumer() (*model.Applicant, error) {
	var applicant *model.Applicant
	err := u.query([][]byte{[]byte("consumer")}, &applicant)
	if err != nil {
		return nil, err
	}
	return applicant, nil
}


