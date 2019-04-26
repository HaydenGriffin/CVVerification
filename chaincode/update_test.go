package main

import (
	"bytes"
	"fmt"
	"github.com/cvverification/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func initApplicantProfile(stub shim.ChaincodeStubInterface, id string) error {

	applicant := model.Applicant{
		Actor: model.Actor{
			ID:       id,
			Username: "applicant" + id,
		},
		Profile: model.ApplicantProfile{},
	}

	err := updateInLedger(stub, model.ActorApplicant, id, applicant)
	if err != nil {
		return err
	}
	return nil
}

func initChaincode(test *testing.T) *shim.MockStub {
	stub := shim.NewMockStub("testingStub", new(CVVerificationChaincode))
	result := stub.MockInit("000", [][]byte{[]byte("init")})

	if result.Status != shim.OK {
		fmt.Println(fmt.Sprintf("init chaincode failed: %v", result.Message))
		test.FailNow()
	}
	return stub
}

func update(test *testing.T, stub *shim.MockStub, args [][]byte) {

	result := stub.MockInvoke("000", append([][]byte{[]byte("invoke"), []byte("update")}, args...))
	fmt.Println("Call:    ", string(args[0]), "(",string(bytes.Join(args,[]byte(", "))),")")
	fmt.Println("RetCode: ", result.Status)
	fmt.Println("RetMsg:  ", result.Message)
	fmt.Println("Payload: ", string(result.Payload))

	if result.Status != shim.OK {
		test.FailNow()
	}
}

func TestInit(test *testing.T) {
	_ = initChaincode(test)
}

func TestProfile(test *testing.T) {
	stub := initChaincode(test)
	stub.MockTransactionStart("1")

	err := initApplicantProfile(stub, "1")
	if err != nil {
		fmt.Println(fmt.Sprintf("save profile failed: %v", err))
		test.FailNow()
	}

	var applicant model.Applicant

	// Attempt to retrieve applicant from profile
	err = getFromLedger(stub, model.ActorApplicant, "1", &applicant)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
		test.FailNow()
	}

	// Attempt to retrieve object that does not exist
	err = getFromLedger(stub, model.ActorApplicant, "1", &applicant)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
	} else {
		test.FailNow()
	}

}