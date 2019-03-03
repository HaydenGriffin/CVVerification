/**
  @Author : hanxiaodong
*/
package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

type UserProfile struct {
	Username	string	`json:"Name"`
	CVHistory []string `json:"CVHistory"`
	Ratings map[string] []CVRating
}

type CVObject struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`
	Speciality	string	`json:"Speciality"`
	CV	string	`json:"CV"`
	CVDate	string	`json:"CVDate"`
}

type CVRating struct {
	Name string `json:"Name"`
	Comment string `json:"Comment"`
	Rating int `json:"Rating"`
}

type ServiceSetup struct {
	ChaincodeID     string
	Client          *channel.Client
}

func registerEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("Register chaincode event failed: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Recieved a chaincode event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("Cannot map chaincode event to eventID(%s)", eventID)
	}
	return nil
}
