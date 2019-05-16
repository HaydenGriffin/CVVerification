package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cvverification/app/crypto"
	"github.com/cvverification/chaincode/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func QueryValid(test *testing.T, stub *shim.MockStub, args [][]byte, responseObject interface{}) {

	// Create a mock transaction proposal
	result := stub.MockInvoke("000", append([][]byte{[]byte("invoke"), []byte("query")}, args...))
	fmt.Println("Calling function: ", string(args[0]), "(", string(bytes.Join(args[1:], []byte(", "))), ")")

	// If the operation is unsuccessful, return an error
	if result.Status != shim.OK {
		fmt.Println(fmt.Sprintf("query invoke function: %v failed: %v", string(args[0]), result.Message))
		test.FailNow()
	}

	// If the responseObject cannot be converted to the specified Object type, return an error
	if responseObject != nil {
		err := json.Unmarshal(result.Payload, responseObject)
		if err != nil {
			fmt.Println(fmt.Sprintf("unable to retrieve responseObject: %v", err))
			test.FailNow()
		}
	}
}

func QueryInvalid(test *testing.T, stub *shim.MockStub, args [][]byte) {

	// Create a mock transaction proposal
	result := stub.MockInvoke("000", append([][]byte{[]byte("invoke"), []byte("query")}, args...))
	fmt.Println("Calling function: ", string(args[0]), "(", string(bytes.Join(args[1:], []byte(", "))), ")")

	// If the operation is successful, return an error
	if result.Status == shim.OK {
		fmt.Println(fmt.Sprintf("query invoke function unexpectedly succeeded: %v", string(args[0])))
		test.FailNow()
	}
}

func TestProfileKey(test *testing.T) {
	stub := InitChaincode(test)

	stub.MockTransactionStart("1")
	InitApplicantProfile(test, stub, "1")
	stub.MockTransactionEnd("1")

	var applicant model.Applicant
	// Attempt to retrieve applicant from profile
	err := getFromLedger(stub, model.ActorApplicant, "1", &applicant)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
		test.FailNow()
	}

	_, publicKey := crypto.GenerateKeyPair(2048)
	publicKeyBytes := crypto.PublicKeyToBytes(publicKey)

	stub.MockTransactionStart("2")
	UpdateValid(test, stub, [][]byte{[]byte("saveprofilekey"), publicKeyBytes, []byte("1")}, nil)
	stub.MockTransactionEnd("2")

	// Query profile key invalid operations
	QueryInvalid(test, stub, [][]byte{[]byte("applicantkey"), []byte("2")})
	QueryInvalid(test, stub, [][]byte{[]byte("applicantkey"), []byte("")})

	var applicantProfile model.ApplicantProfile

	// Query profile key valid operation
	QueryValid(test, stub, [][]byte{[]byte("applicantkey"), []byte("1")}, &applicantProfile)

	if applicantProfile.PublicKey != string(publicKeyBytes) {
		fmt.Println("query profile key incorrect")
	}
}

func TestCV(test *testing.T) {
	stub := InitChaincode(test)

	// Create a test CV object
	cvToAdd := model.CVObject{
		Name:       "Applicant One",
		Date:       "2019-04-20",
		Industry:   "Test Industry",
		Level:      "Junior",
		CV:         "Test CV Personal Statement",
		CVSections: make(map[string]string),
	}
	cvToAdd.CVSections["Skills"] = "Test Skills"
	cvToAdd.CVSections["Experience"] = "Test Experience"
	cvToAdd.CVSections["Education"] = "Test Education"

	// Convert object to byte
	cvToAddByte := convertObjectToByteValid(test, cvToAdd)

	cvID := "applicant1CV1"

	// SaveCV
	stub.MockTransactionStart("1")
	UpdateValid(test, stub, [][]byte{[]byte("savecv"), cvToAddByte, []byte(cvID)}, nil)
	stub.MockTransactionEnd("1")

	// Query invalid CVs
	QueryInvalid(test, stub, [][]byte{[]byte("cv"), []byte("")})
	QueryInvalid(test, stub, [][]byte{[]byte("cv"), []byte("applicant1")})

	var cvRetrieved model.CVObject

	// Query valid CV
	QueryValid(test, stub, [][]byte{[]byte("cv"), []byte(cvID)}, &cvRetrieved)

	// Ensure the retrieved name is the same
	if cvRetrieved.Name != cvToAdd.Name {
		fmt.Println(fmt.Sprintf("retrieved cv name: %v does not match added cv name: %v", cvRetrieved.Name, cvToAdd.Name))
		test.FailNow()
	}

	if cvRetrieved.Date != cvToAdd.Date {
		fmt.Println(fmt.Sprintf("retrieved cv date: %v does not match added cv date: %v", cvRetrieved.Date, cvToAdd.Date))
		test.FailNow()
	}

	if cvRetrieved.Industry != cvToAdd.Industry {
		fmt.Println(fmt.Sprintf("retrieved cv industry: %v does not match added cv industry: %v", cvRetrieved.Industry, cvToAdd.Industry))
		test.FailNow()
	}

	if cvRetrieved.Level != cvToAdd.Level {
		fmt.Println(fmt.Sprintf("retrieved cv level: %v does not match added cv level: %v", cvRetrieved.Level, cvToAdd.Level))
		test.FailNow()
	}

	if cvRetrieved.CV != cvToAdd.CV {
		fmt.Println(fmt.Sprintf("retrieved cv cv: %v does not match added cv cv: %v", cvRetrieved.CV, cvToAdd.CV))
		test.FailNow()
	}

	if cvRetrieved.CVSections["Skills"] != cvToAdd.CVSections["Skills"] {
		fmt.Println(fmt.Sprintf("retrieved cv skills: %v does not match added cv skills: %v", cvRetrieved.CVSections["Skills"], cvToAdd.CVSections["Skills"]))
		test.FailNow()
	}

	if cvRetrieved.CVSections["Experience"] != cvToAdd.CVSections["Experience"] {
		fmt.Println(fmt.Sprintf("retrieved cv experience: %v does not match added cv experience: %v", cvRetrieved.CVSections["Experience"], cvToAdd.CVSections["Experience"]))
		test.FailNow()
	}

	if cvRetrieved.CVSections["Education"] != cvToAdd.CVSections["Education"] {
		fmt.Println(fmt.Sprintf("retrieved cv education: %v does not match added cv education: %v", cvRetrieved.CVSections["Education"], cvToAdd.CVSections["Education"]))
		test.FailNow()
	}
}

func TestCVs(test *testing.T) {
	stub := InitChaincode(test)

	cv1 := model.CVObject{
		Name:       "Applicant One",
		Date:       "2019-04-20",
		Industry:   "Test Industry",
		Level:      "Junior",
		Status: model.CVSubmitted,
		CV:         "Test CV Personal Statement",
		CVSections: make(map[string]string),
	}
	cv1.CVSections["Skills"] = "Test Skills"
	cv1.CVSections["Experience"] = "Test Experience"
	cv1.CVSections["Education"] = "Test Education"

	cvByte := convertObjectToByteValid(test, cv1)

	UpdateValid(test, stub, [][]byte{[]byte("savecv"), cvByte, []byte("cv1")}, nil)


	// SaveCV
	stub.MockTransactionStart("1")
	err := updateInLedger(stub, model.ObjectTypeCV, "applicant1CV1", cv1)
	if err != nil {
		fmt.Println(fmt.Sprintf("save cv in ledger error: %v", err))
	}
	stub.MockTransactionEnd("1")

	cv2 := model.CVObject{
		Name:       "Applicant Two",
		Date:       "2019-04-20",
		Status: model.CVInReview,
		Industry:   "Computer Science",
		Level:      "Intermediate",
		CV:         "Test CV Personal Statement",
		CVSections: make(map[string]string),
	}

	// SaveCV
	stub.MockTransactionStart("2")
	err = updateInLedger(stub, model.ObjectTypeCV, "applicant2CV1", cv2)
	if err != nil {
		fmt.Println(fmt.Sprintf("save cv in ledger error: %v", err))
	}
	stub.MockTransactionEnd("2")

	cv3 := model.CVObject{
		Name:       "Applicant Three",
		Date:       "2019-04-20",
		Status: model.CVSubmitted,
		Industry:   "Paramedic",
		Level:      "Senior",
		CV:         "Test CV Personal Statement",
		CVSections: make(map[string]string),
	}

	// SaveCV
	stub.MockTransactionStart("3")
	err = updateInLedger(stub, model.ObjectTypeCV, "applicant3CV1", cv3)
	if err != nil {
		fmt.Println(fmt.Sprintf("save cv in ledger error: %v", err))
	}
	stub.MockTransactionEnd("3")

	cv4 := model.CVObject{
		Name:       "Applicant Three",
		Date:       "2019-04-20",
		Status: model.CVInDraft,
		Industry:   "Paramedic",
		Level:      "Junior",
		CV:         "Test CV Personal Statement",
		CVSections: make(map[string]string),
	}

	// SaveCV
	stub.MockTransactionStart("4")
	err = updateInLedger(stub, model.ObjectTypeCV, "applicant3CV2", cv4)
	if err != nil {
		fmt.Println(fmt.Sprintf("save cv in ledger error: %v", err))
	}
	stub.MockTransactionEnd("4")

	cv5 := model.CVObject{
		Name:       "Applicant Four",
		Date:       "2019-04-20",
		Status: model.CVInReview,
		Industry:   "Marketing Director",
		Level:      "Administrator",
		CV:         "Test CV Personal Statement",
		CVSections: make(map[string]string),
	}

	// SaveCV
	stub.MockTransactionStart("5")
	err = updateInLedger(stub, model.ObjectTypeCV, "applicant4CV1", cv5)
	if err != nil {
		fmt.Println(fmt.Sprintf("save cv in ledger error: %v", err))
	}
	stub.MockTransactionEnd("5")


	// Perform invalid queries for CVs and ensure they fail
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVInReview), []byte(""), []byte(model.ActorApplicant)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVInDraft), []byte(""), []byte(model.ActorApplicant)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVSubmitted), []byte(""), []byte(model.ActorApplicant)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVWithdrawn), []byte(""), []byte(model.ActorApplicant)})

	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVInReview), []byte(""), []byte(model.ActorAdmin)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVInDraft), []byte(""), []byte(model.ActorAdmin)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVSubmitted), []byte(""), []byte(model.ActorAdmin)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVWithdrawn), []byte(""), []byte(model.ActorAdmin)})

	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVInDraft), []byte(""), []byte(model.ActorVerifier)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVSubmitted), []byte(""), []byte(model.ActorVerifier)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVWithdrawn), []byte(""), []byte(model.ActorVerifier)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVInReview), []byte("test12345"), []byte(model.ActorVerifier)})


	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVInDraft), []byte(""), []byte(model.ActorEmployer)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVInReview), []byte(""), []byte(model.ActorEmployer)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVWithdrawn), []byte(""), []byte(model.ActorEmployer)})
	QueryInvalid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVSubmitted), []byte("test12345"), []byte(model.ActorEmployer)})

	var cvList []model.CVObject

	QueryValid(test, stub, [][]byte{[]byte("cvs"), []byte(model.CVInReview), []byte(""), []byte(model.ActorVerifier)}, cvList)

	if len(cvList) != 2 {
		fmt.Println(fmt.Sprintf("incorrect number of cvs returned, expected 2, got: %v", len(cvList)))
	}


}

func TestCVReviews(test *testing.T) {

}