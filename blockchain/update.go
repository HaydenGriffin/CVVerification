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

// UpdateAdd allow to add a resource into the blockchain
func (u *User) UpdateAddCV(resourceID, resourceDescription string) error {
	return u.update([][]byte{[]byte("add"), []byte(resourceID), []byte(resourceDescription)}, nil)
}

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
