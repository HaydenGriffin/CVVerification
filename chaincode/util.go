package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// convertObjectToByte
func convertObjectToByte(object interface{}) ([]byte, error) {
	byteArray, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("error occurred whilst marshalling object to byte array: %v", err)
	}
	return byteArray, nil
}

// convertByteToObject
func convertByteToObject(byteArray []byte, result interface{}) error {
	err := json.Unmarshal(byteArray, result)
	if err != nil {
		return fmt.Errorf("error occurred whilst unmarshalling byte array to object: %v", err)
	}
	return nil
}

// getFromLedger retrieve an object from the ledger
func getFromLedger(stub shim.ChaincodeStubInterface, objectType string, id string, result interface{}) error {
	key, err := stub.CreateCompositeKey(objectType, []string{id})
	if err != nil {
		return fmt.Errorf("unable to create the object key for the ledger: %v", err)
	}
	resultAsByte, err := stub.GetState(key)
	if err != nil {
		return fmt.Errorf("unable to retrieve the object in the ledger: %v", err)
	}
	if resultAsByte == nil {
		return fmt.Errorf("the object doesn't exist in the ledger")
	}
	err = convertByteToObject(resultAsByte, result)
	if err != nil {
		return fmt.Errorf("unable to convert the result to object: %v", err)
	}
	return nil
}

// updateInLedger update an object in the ledger
func updateInLedger(stub shim.ChaincodeStubInterface, objectType string, id string, object interface{}) error {
	key, err := stub.CreateCompositeKey(objectType, []string{id})
	if err != nil {
		return fmt.Errorf("unable to create the object key for the ledger: %v", err)
	}

	objectAsByte, err := convertObjectToByte(object)
	if err != nil {
		return err
	}
	err = stub.PutState(key, objectAsByte)
	if err != nil {
		return fmt.Errorf("unable to put the object in the ledger: %v", err)
	}
	return nil
}
