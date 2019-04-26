package main

import (
	"bytes"
	"fmt"
	"github.com/cvverification/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func TestInit(test *testing.T) {
	_ = InitChaincode(test)
}

func TestInitProfile(test *testing.T) {
	stub := InitChaincode(test)
	stub.MockTransactionStart("1")

	applicantIDs := [4]string{"1", "2", "3", "4"}

	for _, id := range applicantIDs {
		applicant := model.Applicant{
			Actor: model.Actor{
				ID:       id,
				Username: "applicant" + id,
			},
			Profile: model.ApplicantProfile{},
		}

		err := updateInLedger(stub, model.ActorApplicant, id, applicant)
		if err != nil {
			fmt.Println(fmt.Sprintf("save profile failed: %v", err))
			test.FailNow()
		}
	}

	applicant := model.Applicant{}

	err := getFromLedger(stub, model.ActorApplicant, "1", &applicant)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
		test.FailNow()
	}
}


func InitChaincode(test *testing.T) *shim.MockStub {
	stub := shim.NewMockStub("testingStub", new(CVVerificationChaincode))
	result := stub.MockInit("000", [][]byte{[]byte("init")})

	if result.Status != shim.OK {
		fmt.Println(fmt.Sprintf("init chaincode failed: %v", result.Message))
		test.FailNow()
	}
	return stub
}

func Update(test *testing.T, stub *shim.MockStub, args [][]byte) {

	result := stub.MockInvoke("000", append([][]byte{[]byte("invoke"), []byte("update")}, args...))
	fmt.Println("Call:    ", string(args[0]), "(",string(bytes.Join(args,[]byte(", "))),")")
	fmt.Println("RetCode: ", result.Status)
	fmt.Println("RetMsg:  ", result.Message)
	fmt.Println("Payload: ", string(result.Payload))

	if result.Status != shim.OK {
		test.FailNow()
	}
}