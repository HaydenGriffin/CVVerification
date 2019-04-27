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

func InitChaincode(test *testing.T) *shim.MockStub {
	cvvc := new(CVVerificationChaincode)
	cvvc.testing = true
	stub := shim.NewMockStub("testingStub", cvvc)
	result := stub.MockInit("000", [][]byte{[]byte("init")})

	if result.Status != shim.OK {
		fmt.Println(fmt.Sprintf("init chaincode failed: %v", result.Message))
		test.FailNow()
	}
	return stub
}

func InitApplicantProfile(test *testing.T, stub shim.ChaincodeStubInterface, id string) {

	applicant := model.Applicant{
		Actor: model.Actor{
			ID:       id,
			Username: "applicant" + id,
		},
		Profile: model.ApplicantProfile{},
	}

	err := updateInLedger(stub, model.ActorApplicant, id, applicant)
	if err != nil {
		fmt.Println(fmt.Sprintf("init applicant profile failed: %v", err))
		test.FailNow()
	}
}

func InitVerifierProfile(test *testing.T, stub shim.ChaincodeStubInterface, id string) {

	verifier := model.Verifier{
		Actor: model.Actor{
			ID:       id,
			Username: "verifier" + id,
		},
		Profile: model.VerifierProfile{},
	}

	err := updateInLedger(stub, model.ActorVerifier, id, verifier)
	if err != nil {
		fmt.Println(fmt.Sprintf("init verifier profile failed: %v", err))
		test.FailNow()
	}
}

func InitEmployerProfile(test *testing.T, stub shim.ChaincodeStubInterface, id string) {

	employer := model.Employer{
		Actor: model.Actor{
			ID:       id,
			Username: "employer" + id,
		},
		Profile: model.EmployerProfile{},
	}

	err := updateInLedger(stub, model.ActorEmployer, id, employer)
	if err != nil {
		fmt.Println(fmt.Sprintf("init employer profile failed: %v", err))
		test.FailNow()
	}
}

func UpdateValid(test *testing.T, stub *shim.MockStub, args [][]byte, responseObject interface{}) {

	result := stub.MockInvoke("000", append([][]byte{[]byte("invoke"), []byte("update")}, args...))
	fmt.Println("Calling function: ", string(args[0]), "(", string(bytes.Join(args[1:], []byte(", "))), ")")

	if result.Status != shim.OK {
		fmt.Println(fmt.Sprintf("update invoke function: %v failed: %v", string(args[0]), result.Message))
		test.FailNow()
	}

	if responseObject != nil {
		err := json.Unmarshal(result.Payload, responseObject)
		if err != nil {
			fmt.Println(fmt.Sprintf("unable to retrieve responseObject: %v", err))
			test.FailNow()
		}
	}
}

func UpdateInvalid(test *testing.T, stub *shim.MockStub, args [][]byte) {

	result := stub.MockInvoke("000", append([][]byte{[]byte("invoke"), []byte("update")}, args...))
	fmt.Println("Calling function: ", string(args[0]), "(", string(bytes.Join(args[1:], []byte(", "))), ")")

	if result.Status == shim.OK {
		fmt.Println(fmt.Sprintf("update invoke function unexpectedly succeeded: %v", string(args[0])))
		test.FailNow()
	}
}

func TestInit(test *testing.T) {
	_ = InitChaincode(test)
}

func TestProfileCreation(test *testing.T) {
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

	// Attempt to retrieve applicant that does not exist
	err = getFromLedger(stub, model.ActorApplicant, "2", &applicant)
	if err != nil {
	} else {
		fmt.Println(fmt.Sprintf("found profile that should not be there"))
		test.FailNow()
	}
}

func TestSaveCV(test *testing.T) {
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

	cvToAddByte := convertObjectToByteValid(test, cvToAdd)

	cvID := "applicant1CV1"

	// Attempt to save CV
	stub.MockTransactionStart("2")
	UpdateInvalid(test, stub, [][]byte{[]byte("savecv"), []byte(""), []byte(cvID)})
	UpdateValid(test, stub, [][]byte{[]byte("savecv"), cvToAddByte, []byte(cvID)}, nil)
	stub.MockTransactionEnd("2")

	var cvRetrieved model.CVObject
	// Attempt to retrieve cv from profile
	err = getFromLedger(stub, model.ObjectTypeCV, cvID, &cvRetrieved)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve cv failed: %v", err))
		test.FailNow()
	}

	// Ensure that the retrieved CV is in the correct status
	if cvRetrieved.Status != model.CVInDraft {
		fmt.Println(fmt.Sprintf("cv status incorrect, is: %v, should be: %v", cvRetrieved.Status, model.CVInDraft))
	}

	// Attempt to save CV to profile
	stub.MockTransactionStart("3")
	UpdateInvalid(test, stub, [][]byte{[]byte("saveprofilecv"), []byte(""), []byte("1")})
	UpdateInvalid(test, stub, [][]byte{[]byte("saveprofilecv"), []byte(""), []byte("2")})
	UpdateValid(test, stub, [][]byte{[]byte("saveprofilecv"), []byte(cvID), []byte("1")}, nil)
	stub.MockTransactionEnd("3")

	// Retrieve updated profile
	var updatedApplicant model.Applicant
	err = getFromLedger(stub, model.ActorApplicant, "1", &updatedApplicant)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
		test.FailNow()
	}

	if len(updatedApplicant.Profile.CVHistory) == 0 {
		fmt.Println("updated profile does not contain new CV application")
		test.FailNow()
	}

	if updatedApplicant.Profile.CVHistory[0] != cvID {
		fmt.Println("updated profile CV application is incorrect")
	}

}

func TestTransitionCV(test *testing.T) {
	stub := InitChaincode(test)

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

	cvToAddByte, err := convertObjectToByte(cvToAdd)
	if err != nil {
		fmt.Println(fmt.Sprintf("convert cv application to byte failed: %v", err))
		test.FailNow()
	}

	cvID := "applicant1CV1"

	// Attempt to save CV
	stub.MockTransactionStart("2")
	UpdateValid(test, stub, [][]byte{[]byte("savecv"), cvToAddByte, []byte(cvID)}, nil)
	stub.MockTransactionEnd("2")

	var cvRetrieved model.CVObject

	// Attempt to retrieve cv from profile
	err = getFromLedger(stub, model.ObjectTypeCV, cvID, &cvRetrieved)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve cv failed: %v", err))
		test.FailNow()
	}

	// Ensure that the retrieved CV is in the correct status
	if cvRetrieved.Status != model.CVInDraft {
		fmt.Println(fmt.Sprintf("cv status incorrect, is: %v, should be: %v", cvRetrieved.Status, model.CVInDraft))
		test.FailNow()
	}

	// Check transitioning a draft CV against invalid operations
	// Transition to same status (should fail)
	transitionCVInvalid(test, model.ActorApplicant, model.CVInDraft, cvRetrieved)

	// Transition as non-applicant user (should fail)
	transitionCVInvalid(test, model.ActorEmployer, model.CVInReview, cvRetrieved)
	transitionCVInvalid(test, model.ActorVerifier, model.CVInReview, cvRetrieved)
	transitionCVInvalid(test, model.ActorAdmin, model.CVInReview, cvRetrieved)

	//Check transitioning a draft CV against valid operations
	transitionCVValid(test, model.ActorApplicant, model.CVInReview, cvRetrieved)
	transitionCVValid(test, model.ActorApplicant, model.CVSubmitted, cvRetrieved)
	transitionCVValid(test, model.ActorApplicant, model.CVWithdrawn, cvRetrieved)

	cvRetrieved.Status = model.CVInReview

	// Check transitioning a CV in review against invalid operations
	// Transition to same status (should fail)
	transitionCVInvalid(test, model.ActorApplicant, model.CVInReview, cvRetrieved)

	// Transition as non-applicant user (should fail)
	transitionCVInvalid(test, model.ActorEmployer, model.CVSubmitted, cvRetrieved)
	transitionCVInvalid(test, model.ActorVerifier, model.CVSubmitted, cvRetrieved)
	transitionCVInvalid(test, model.ActorAdmin, model.CVSubmitted, cvRetrieved)

	//Check transitioning a CV in review against valid operations
	transitionCVValid(test, model.ActorApplicant, model.CVSubmitted, cvRetrieved)
	transitionCVValid(test, model.ActorApplicant, model.CVInDraft, cvRetrieved)
	transitionCVValid(test, model.ActorApplicant, model.CVWithdrawn, cvRetrieved)

	cvRetrieved.Status = model.CVSubmitted

	// Check transitioning a submitted CV against invalid operations
	// Transition to same status (should fail)
	transitionCVInvalid(test, model.ActorApplicant, model.CVSubmitted, cvRetrieved)

	// Transition to draft (should fail)
	transitionCVInvalid(test, model.ActorApplicant, model.CVInDraft, cvRetrieved)

	// Transition as non-applicant user (should fail)
	transitionCVInvalid(test, model.ActorEmployer, model.CVInReview, cvRetrieved)
	transitionCVInvalid(test, model.ActorVerifier, model.CVInReview, cvRetrieved)
	transitionCVInvalid(test, model.ActorAdmin, model.CVInReview, cvRetrieved)

	//Check transitioning a submitted CV against valid operations
	transitionCVValid(test, model.ActorApplicant, model.CVInReview, cvRetrieved)
	transitionCVValid(test, model.ActorApplicant, model.CVWithdrawn, cvRetrieved)

	cvRetrieved.Status = model.CVWithdrawn

	// Check transitioning a withdrawn CV against invalid operations
	// Transition to same status (should fail)
	transitionCVInvalid(test, model.ActorApplicant, model.CVWithdrawn, cvRetrieved)

	// Transition as non-applicant user (should fail)
	transitionCVInvalid(test, model.ActorEmployer, model.CVInReview, cvRetrieved)
	transitionCVInvalid(test, model.ActorVerifier, model.CVInReview, cvRetrieved)
	transitionCVInvalid(test, model.ActorAdmin, model.CVInReview, cvRetrieved)

	//Check transitioning a withdrawn CV against invalid operations
	transitionCVInvalid(test, model.ActorApplicant, model.CVInReview, cvRetrieved)
	transitionCVInvalid(test, model.ActorApplicant, model.CVInDraft, cvRetrieved)
	transitionCVInvalid(test, model.ActorApplicant, model.CVWithdrawn, cvRetrieved)
}

func TestSaveProfileKey(test *testing.T) {
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
	UpdateInvalid(test, stub, [][]byte{[]byte("saveprofilekey"), nil, []byte("")})
	UpdateInvalid(test, stub, [][]byte{[]byte("saveprofilekey"), publicKeyBytes, []byte("2")})
	UpdateInvalid(test, stub, [][]byte{[]byte("saveprofilekey"), nil, []byte("1")})
	UpdateValid(test, stub, [][]byte{[]byte("saveprofilekey"), publicKeyBytes, []byte("1")}, nil)
	stub.MockTransactionEnd("2")

	err = getFromLedger(stub, model.ActorApplicant, "1", &applicant)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
		test.FailNow()
	}

	if applicant.Profile.PublicKey != string(publicKeyBytes) {
		fmt.Println("retrieved profile key incorrect")
	}
}

func TestSaveReview(test *testing.T) {
	stub := InitChaincode(test)

	// Save applicant to ledger
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

	privateKey, publicKey := crypto.GenerateKeyPair(2048)
	publicKeyBytes := crypto.PublicKeyToBytes(publicKey)

	stub.MockTransactionStart("2")
	UpdateValid(test, stub, [][]byte{[]byte("saveprofilekey"), publicKeyBytes, []byte("1")}, nil)
	stub.MockTransactionEnd("2")

	review := model.CVReview{
		Name:         "Reviewer One",
		Organisation: "Test Organisation",
		Type:         "Identification",
		Comment:      "Test Comment",
		Rating:       10,
	}

	reviewBytes := convertObjectToByteValid(test, review)

	encryptedReviewBytes := crypto.EncryptWithPublicKey(reviewBytes, publicKey)
	if len(encryptedReviewBytes) == 0 {
		fmt.Println("length of encrypted review is 0")
		test.FailNow()
	}

	stub.MockTransactionStart("3")
	UpdateValid(test, stub, [][]byte{[]byte("verifiersavereview"), []byte("1"), []byte("applicant1CV1"), encryptedReviewBytes, []byte("verifier1")}, nil)
	stub.MockTransactionEnd("3")

	err = getFromLedger(stub, model.ActorApplicant, "1", &applicant)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
		test.FailNow()
	}

	if applicant.Profile.Reviews["applicant1CV1"]["verifier1"] == nil {
		fmt.Println("could not find saved encrypted review")
		test.FailNow()
	}

	decryptedReviewBytes, err := crypto.DecryptWithPrivateKey(applicant.Profile.Reviews["applicant1CV1"]["verifier1"], privateKey)

	if string(decryptedReviewBytes) != string(reviewBytes) {
		fmt.Println("decrypted review bytes different to review bytes")
		test.FailNow()
	}

	var decryptedReview model.CVReview

	convertByteToObjectValid(test, decryptedReviewBytes, &decryptedReview)

	if decryptedReview.Name != review.Name {
		fmt.Println(fmt.Sprintf("decrypted review name: %v does not match review name: %v", decryptedReview.Name, review.Name))
	}

	if decryptedReview.Rating != review.Rating {
		fmt.Println(fmt.Sprintf("decrypted review rating: %v does not match review rating: %v", decryptedReview.Rating, review.Rating))
	}

	if decryptedReview.Type != review.Type {
		fmt.Println(fmt.Sprintf("decrypted review type: %v does not match review type: %v", decryptedReview.Type, review.Type))
	}

	if decryptedReview.Organisation != review.Organisation {
		fmt.Println(fmt.Sprintf("decrypted review organisation: %v does not match review organisation: %v", decryptedReview.Organisation, review.Organisation))
	}

	if decryptedReview.Comment != review.Comment {
		fmt.Println(fmt.Sprintf("decrypted review comment: %v does not match review comment: %v", decryptedReview.Comment, review.Comment))
	}
}

func TestSaveOrganisation(test *testing.T) {
	stub := InitChaincode(test)

	stub.MockTransactionStart("1")
	InitVerifierProfile(test, stub, "1")
	stub.MockTransactionEnd("1")

	stub.MockTransactionStart("2")
	UpdateInvalid(test, stub, [][]byte{[]byte("verifiersaveorganisation"), []byte("Test Organisation"), []byte("2")})
	UpdateInvalid(test, stub, [][]byte{[]byte("verifiersaveorganisation"), []byte(""), []byte("1")})
	UpdateValid(test, stub, [][]byte{[]byte("verifiersaveorganisation"), []byte("Test Organisation"), []byte("1")}, nil)
	stub.MockTransactionEnd("2")

	// Retrieve updated profile
	var updatedVerifier model.Verifier
	err := getFromLedger(stub, model.ActorVerifier, "1", &updatedVerifier)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
		test.FailNow()
	}

	if updatedVerifier.Profile.Organisation != "Test Organisation" {
		fmt.Println(fmt.Sprintf("updated verifier profiile organisation should be: %v, is: %v", "Test Organisation", updatedVerifier.Profile.Organisation))
	}
}

func TestPublishReviews(test *testing.T) {
	stub := InitChaincode(test)

	// Save applicant to ledger
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

	var reviews []model.CVReview

	review1 := model.CVReview{
		Name:         "Reviewer One",
		Organisation: "Test Organisation 1",
		Type:         "Identification",
		Comment:      "Test Comment",
		Rating:       10,
	}

	reviews = append(reviews, review1)

	review2 := model.CVReview{
		Name:         "Reviewer Two",
		Organisation: "Test Organisation 2",
		Type:         "Credit Check",
		Comment:      "Test Comment",
		Rating:       9,
	}

	reviews = append(reviews, review2)

	review3 := model.CVReview{
		Name:         "Reviewer Three",
		Organisation: "Test Organisation 3",
		Type:         "Credit Check",
		Comment:      "Test Comment",
		Rating:       8,
	}

	reviews = append(reviews, review3)

	reviewsByte := convertObjectToByteValid(test, reviews)

	stub.MockTransactionStart("2")
	UpdateInvalid(test, stub, [][]byte{[]byte("publishreviews"), []byte("applicant1CV1"), reviewsByte, []byte("2")})
	UpdateInvalid(test, stub, [][]byte{[]byte("publishreviews"), []byte("applicant1CV1"), nil, []byte("1")})
	UpdateValid(test, stub, [][]byte{[]byte("publishreviews"), []byte("applicant1CV1"), reviewsByte, []byte("1")}, nil)
	stub.MockTransactionEnd("2")

	err = getFromLedger(stub, model.ActorApplicant, "1", &applicant)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
		test.FailNow()
	}

	if len(applicant.Profile.PublicReviews["applicant1CV1"]) != 3 {
		fmt.Println("published reviews could not be found")
	}
}

func TestEmployerSaveCV(test *testing.T) {
	stub := InitChaincode(test)

	// Save employer to ledger
	stub.MockTransactionStart("1")
	InitEmployerProfile(test, stub, "1")
	stub.MockTransactionEnd("1")

	stub.MockTransactionStart("2")
	UpdateInvalid(test, stub, [][]byte{[]byte("employersavecv"), []byte("applicant1CV1"), []byte("2")})
	UpdateInvalid(test, stub, [][]byte{[]byte("employersavecv"), []byte("applicant1CV1"), []byte("")})
	UpdateValid(test, stub, [][]byte{[]byte("employersavecv"), []byte("applicant1CV1"), []byte("1")}, nil)
	stub.MockTransactionEnd("2")

	var employer model.Employer

	err := getFromLedger(stub, model.ActorEmployer, "1", &employer)
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve profile failed: %v", err))
		test.FailNow()
	}

	var found = false
	for _, cvID := range employer.Profile.ProspectiveCVs {
		if cvID == "applicant1CV1" {
			found = true
		}
	}

	if !found {
		fmt.Println("could not find cvID in the list of saved CVs")
		test.FailNow()
	}
}
