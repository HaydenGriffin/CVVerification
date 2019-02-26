package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const DOC_TYPE = "cvObj"


// Save user profile
// args: cv
func PutProfile(stub shim.ChaincodeStubInterface, profile UserProfile, usernameHash string) ([]byte, bool) {


	b, err := json.Marshal(profile)
	if err != nil {
		return nil, false
	}

	// Save resume status
	err = stub.PutState(usernameHash, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// Get CV
// args: CVHash
func getProfile(stub shim.ChaincodeStubInterface, usernameHash string) (UserProfile, bool) {

	var profile UserProfile

	b, err := stub.GetState(usernameHash)

	// error, return empty profile object
	if err != nil {
		return profile, false
	}

	// no value found for key specified
	if b == nil {
		return profile, false
	}

	// Deserialize the queried value
	err = json.Unmarshal(b, &profile)
	if err != nil {
		return profile, false
	}

	// Success
	return profile, true
}

// Add CV Chaincode
// args: CV object
// CV Hash is key, CVObject is the value
func (t *CVTrackerChaincode) QueryProfileByHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check whether the number of arguments is sufficient
	if len(args) != 1 {
		return shim.Error("The number of arguments is incorrect for the method.")
	}

	profileHash := args[0]

	profile, success := getProfile(stub, profileHash)

	if !success {
		return shim.Error("An error occurred whilst querying the profile")
	}

	result, err := json.Marshal(profile)
	if err != nil {
		return shim.Error("An error occurred whilst marshalling the profile object ")
	}
	return shim.Success(result)
}

// Add CV Chaincode
// args: CV object
// CV Hash is key, CVObject is the value
func (t *CVTrackerChaincode) saveProfile(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check whether the number of arguments is sufficient
	if len(args) != 3 {
		return shim.Error("The number of arguments is incorrect for the method.")
	}

	var profile UserProfile
	err := json.Unmarshal([]byte(args[0]), &profile)
	if err != nil {
		return shim.Error("An error occurred whilst deserialising the object")
	}

	_, success := PutProfile(stub, profile, args[1])
	if !success {
		return shim.Error("An error occurred whilst creating the profile")
	}

	err = stub.SetEvent(args[2], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("Successfully created and saved profile"))
}

// Add CV Chaincode
// args: CV object
// CV Hash is key, CVObject is the value
func (t *CVTrackerChaincode) updateProfile(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check whether the number of arguments is sufficient
	if len(args) != 3 {
		return shim.Error("The number of arguments is incorrect for the method.")
	}

	var profile UserProfile

	// get the profile
	profile, exist := getProfile(stub, args[0])

	if !exist {
		return shim.Error("An error occurred whilst retrieving the user profile")
	}

	//change to append to front
	// append the new CV to the history
	profile.CVHistory = append(profile.CVHistory, args[1])


	// put the updated profile back to blockchain
	_, success := PutProfile(stub, profile, args[0])
	if !success {
		return shim.Error("An error occurred whilst creating the profile")
	}

	var err error
	err = stub.SetEvent(args[2], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("Successfully updated the user profile"))
}

// Save cv
// args: cv
func PutCV(stub shim.ChaincodeStubInterface, cv CVObject, cvHash string) ([]byte, bool) {

	cv.ObjectType = DOC_TYPE

	b, err := json.Marshal(cv)
	if err != nil {
		return nil, false
	}

	// Save resume status
	err = stub.PutState(cvHash, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// Add CV Chaincode
// args: CV objectQueryProfileByHash
// CV Hash is key, CVObject is the value
func (t *CVTrackerChaincode) addCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check whether the number of arguments is sufficient
	if len(args) != 3 {
		return shim.Error("The number of arguments is incorrect for the method.")
	}

	var cv CVObject
	err := json.Unmarshal([]byte(args[0]), &cv)
	if err != nil {
		return shim.Error("An error occurred whilst deserialising the object")
	}

	_, success := PutCV(stub, cv, args[1])
	if !success {
		return shim.Error("An error occurred whilst saving the CV")
	}

	err = stub.SetEvent(args[2], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("Successfully saved the CV"))
}


// Get CV
// args: CVHash
func GetCV(stub shim.ChaincodeStubInterface, cvHash string) (CVObject, bool) {

	var cv CVObject

	b, err := stub.GetState(cvHash)

	// error, return empty CV object
	if err != nil {
		return cv, false
	}

	// no value found for key specified
	if b == nil {
		return cv, false
	}

	// Deserialize the queried value
	err = json.Unmarshal(b, &cv)
	if err != nil {
		return cv, false
	}

	// Success
	return cv, true
}

// Get CV
// args: CVHash
func (t *CVTrackerChaincode) GetCVFromProfileHash(stub shim.ChaincodeStubInterface, args[] string) pb.Response {

	var profile UserProfile
	//var cv CVObject

	// get the profile
	profile, success := getProfile(stub, args[0])

	if !success {
		return shim.Error("An error occurred whilst retrieving the profile")
	}

	cvHash := profile.CVHistory[len(profile.CVHistory)-1]

	cv, success := GetCV(stub, cvHash)

	// no value found for key specified
	if !success {
		return shim.Error("An error occurred whilst retrieving the CV")
	}

	result, err := json.Marshal(cv)
	if err != nil {
		return shim.Error("Failed to marshal CV object")
	}
	return shim.Success(result)
}

/*// query
// Every readonly functions in the ledger will be here
func (t *CVTrackerChaincode) queryCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### CVTrackerChaincode queryCV ###########")

	// Check whether the number of arguments is sufficient
	if len(args) != 1 {
		return shim.Error("The number of arguments is incorrect for the method.")
	}

	// Get the state of the value matching the key specified
	state, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get state of cv")
	}

	if state == nil {
		return shim.Error("No CV exists for the specified key")
	}

	var cv CVObject
	err = json.Unmarshal(state, &cv)
	if err != nil {
		return shim.Error("Failed to deserialize CV object")
	}

	iterator, err := stub.GetHistoryForKey(cv.CVHash)
	if err != nil {
		return shim.Error("Failed to query historical data for specified key")
	}
	defer iterator.Close()

	var historyItems []HistoryItem
	var historyCV CVObject
	for iterator.HasNext() {
		historyData, err := iterator.Next()
		if err != nil {
			return shim.Error("Failed to get CV historical data")
		}
		var historyItem HistoryItem
		historyItem.TxId = historyData.TxId
		json.Unmarshal(historyData.Value, &historyCV)

		if historyData.Value == nil {
			var empty CVObject
			historyItem.CV = empty
		} else {
			historyItem.CV = historyCV
		}

		historyItems = append(historyItems, historyItem)
	}

	cv.History = historyItems

	result, err := json.Marshal(cv)
	if err != nil {
		return shim.Error("Failed to marshal CV object")
	}
	return shim.Success(result)
}*/
/*
// Update CV Chaincode
// args: CV object
func (t *CVTrackerChaincode) updateCV(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("The number of arguments is incorrect for the method.")
	}

	var cv CVObject
	err := json.Unmarshal([]byte(args[0]), &cv)
	if err != nil {
		return shim.Error("An error occurred whilst deserializing the object")
	}

	result, success := GetCV(stub, cv.CVHash)
	if !success {
		return shim.Error("An error occurred whilst saving the CV")
	}

	result.Name = cv.Name
	result.Speciality = cv.Speciality
	result.CV = cv.CV
	result.CVHash = cv.CVHash
	result.CVDate = cv.CVDate

	_, success = PutCV(stub, result)
	if !success {
		return shim.Error("An error occurred whilst saving the CV")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("Successfully update the CV"))
}*/