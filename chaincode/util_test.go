package main

import (
	"encoding/json"
	"fmt"
	"github.com/cvverification/app/crypto"
	"github.com/cvverification/chaincode/model"
	"testing"
)

func convertObjectToByteValid(test *testing.T, object interface{}) []byte {
	byteArray, err := json.Marshal(object)
	if err != nil {
		fmt.Println(fmt.Sprintf("convert object to byte valid failed, should succeed: %v", err))
		test.FailNow()
	}
	return byteArray
}

func convertByteToObjectValid(test *testing.T, byteArray []byte, result interface{}) {
	err := json.Unmarshal(byteArray, result)
	if err != nil {
		fmt.Println(fmt.Sprintf("convert byte to object failed, should succeed: %v", err))
		test.FailNow()
	}
}

func convertByteToObjectInvalid(test *testing.T, byteArray []byte, result interface{}) {
	err := json.Unmarshal(byteArray, result)
	if err == nil {
		fmt.Println(fmt.Sprintf("convert byte to object succeeded, should fail: %v", err))
		test.FailNow()
	}
}

func transitionCVValid(test *testing.T, actorType, transitionTo string, cv model.CVObject) {
	err := canCVBeTransitioned(actorType, transitionTo, cv)

	if err != nil {
		fmt.Println(fmt.Sprintf("transition applicant CV from: %v to: %v failed, should be successful", cv.Status, transitionTo))
		test.FailNow()
	}
}

func transitionCVInvalid(test *testing.T, actorType, transitionTo string, cv model.CVObject) {
	err := canCVBeTransitioned(actorType, transitionTo, cv)

	if err == nil {
		fmt.Println(fmt.Sprintf("transition applicant CV from: %v to: %v succeeded, should fail", cv.Status, transitionTo))
		test.FailNow()
	}
}

func returnCVValid(test *testing.T, actorType, filter string, cv model.CVObject) {
	returned := returnCV(actorType, filter, cv)
	if !returned {
		fmt.Println(fmt.Sprintf("cv status: %v not returned when it should have been for actor: %v, filter: %v", cv.Status, actorType, filter))
		test.FailNow()
	}
}

func returnCVInvalid(test *testing.T, actorType, filter string, cv model.CVObject) {
	returned := returnCV(actorType, filter, cv)
	if returned {
		fmt.Println(fmt.Sprintf("cv status: %v returned when it shouldn't have been for actor: %v, filter: %v", cv.Status, actorType, filter))
		test.FailNow()
	}
}

func TestObjectConversion(test *testing.T) {

	// Object conversion of CV object
	cv := model.CVObject{
		Name:       "Applicant One",
		Date:       "2019-04-20",
		Industry:   "Test Industry",
		Level:      "Junior",
		Status:     model.CVInDraft,
		CV:         "Test CV Personal Statement",
		CVSections: make(map[string]string),
	}
	cv.CVSections["Skills"] = "Test Skills"
	cv.CVSections["Experience"] = "Test Experience"
	cv.CVSections["Education"] = "Test Education"

	// Convert the CV
	cvByte := convertObjectToByteValid(test, cv)

	if len(cvByte) == 0 {
		fmt.Println("cvByte is length 0")
		test.FailNow()
	}

	var convertedCV model.CVObject

	// Convert byte to object invalid
	convertByteToObjectInvalid(test, nil, convertedCV)
	convertByteToObjectInvalid(test, cvByte, nil)

	// Using non-pointer
	convertByteToObjectInvalid(test, cvByte, convertedCV)

	// Convert byte to object valid
	convertByteToObjectValid(test, cvByte, &convertedCV)

	if convertedCV.Name != "Applicant One" {
		fmt.Println("Converted CV name incorrect")
		test.FailNow()
	}

	if convertedCV.Date != "2019-04-20" {
		fmt.Println("Converted CV date incorrect")
		test.FailNow()
	}

	if convertedCV.Industry != "Test Industry" {
		fmt.Println("Converted CV industry incorrect")
		test.FailNow()
	}

	if convertedCV.Level != "Junior" {
		fmt.Println("Converted CV level incorrect")
		test.FailNow()
	}

	if convertedCV.Status != model.CVInDraft {
		fmt.Println("Converted CV status incorrect")
		test.FailNow()
	}

	if convertedCV.CV != "Test CV Personal Statement" {
		fmt.Println("Converted CV CV incorrect")
		test.FailNow()
	}

	if convertedCV.CVSections["Skills"] != "Test Skills" {
		fmt.Println("Converted CVSection skills incorrect")
		test.FailNow()
	}

	if convertedCV.CVSections["Experience"] != "Test Experience" {
		fmt.Println("Converted CVSection experience incorrect")
		test.FailNow()
	}

	if convertedCV.CVSections["Education"] != "Test Education" {
		fmt.Println("Converted CVSection education incorrect")
		test.FailNow()
	}
}

func TestReviewEncryption(test *testing.T) {
	privateKey1, publicKey1 := crypto.GenerateKeyPair(2048)
	privateKey2, _ := crypto.GenerateKeyPair(2048)

	review := model.CVReview{
		Name:         "Reviewer One",
		Organisation: "Test Organisation",
		Type:         "Identification",
		Comment:      "Test Comment",
		Rating:       10,
	}

	reviewBytes := convertObjectToByteValid(test, review)

	encryptedReviewBytes := crypto.EncryptWithPublicKey(reviewBytes, publicKey1)
	if len(encryptedReviewBytes) == 0 {
		fmt.Println("length of encrypted review is 0")
		test.FailNow()
	}

	// Decrypt with the wrong private key
	decryptedReviewBytes, err := crypto.DecryptWithPrivateKey(encryptedReviewBytes, privateKey2)
	if err == nil {
		fmt.Println("decryption succeeded when it should have failed")
		test.FailNow()
	}

	// Decrypt with the correct private key
	decryptedReviewBytes, err = crypto.DecryptWithPrivateKey(encryptedReviewBytes, privateKey1)
	if err != nil {
		fmt.Println("decryption succeeded when it should have failed")
		test.FailNow()
	}
	if len(decryptedReviewBytes) != len(reviewBytes) {
		fmt.Println(fmt.Sprintf("decrypted review bytes incorrect, is: %v, should be: %v", len(encryptedReviewBytes), len(decryptedReviewBytes)))
	}

	var convertedReview model.CVReview

	convertByteToObjectValid(test, decryptedReviewBytes, &convertedReview)

	if convertedReview != review {
		fmt.Println("converted review is not the same as retrieved review")
		test.FailNow()
	}
}

func TestReturnCV(test *testing.T) {

	cv := model.CVObject{
		Name:       "Applicant One",
		Date:       "2019-04-20",
		Industry:   "Test Industry",
		Level:      "Junior",
		Status:     model.CVInDraft,
		CV:         "Test CV Personal Statement",
		CVSections: make(map[string]string),
	}
	cv.CVSections["Skills"] = "Test Skills"
	cv.CVSections["Experience"] = "Test Experience"
	cv.CVSections["Education"] = "Test Education"

	// Check returning a draft CV against invalid operations
	returnCVInvalid(test, model.ActorAdmin, "", cv)
	returnCVInvalid(test, model.ActorVerifier, "", cv)
	returnCVInvalid(test, model.ActorEmployer, "", cv)
	returnCVInvalid(test, model.ActorApplicant, "", cv)

	cv.Status = model.CVWithdrawn

	// Check returning a withdrawn CV against invalid operations
	returnCVInvalid(test, model.ActorAdmin, "", cv)
	returnCVInvalid(test, model.ActorVerifier, "", cv)
	returnCVInvalid(test, model.ActorEmployer, "", cv)
	returnCVInvalid(test, model.ActorApplicant, "", cv)

	cv.Status = model.CVInReview

	// Check returning a CV in review against invalid operations
	returnCVInvalid(test, model.ActorAdmin, "", cv)
	returnCVInvalid(test, model.ActorEmployer, "", cv)
	returnCVInvalid(test, model.ActorApplicant, "", cv)
	returnCVInvalid(test, model.ActorVerifier, "Test Industry1", cv)
	returnCVInvalid(test, model.ActorVerifier, "a", cv)

	// Check returning a CV in review against valid operations
	returnCVValid(test, model.ActorVerifier, "", cv)
	returnCVValid(test, model.ActorVerifier, "T", cv)
	returnCVValid(test, model.ActorVerifier, "Test", cv)
	returnCVValid(test, model.ActorVerifier, "Test Industry", cv)

	cv.Status = model.CVSubmitted

	// Check returning a submitted CV against invalid operations
	returnCVInvalid(test, model.ActorAdmin, "", cv)
	returnCVInvalid(test, model.ActorApplicant, "", cv)
	returnCVInvalid(test, model.ActorVerifier, "", cv)
	returnCVInvalid(test, model.ActorEmployer, "Test Industry1", cv)
	returnCVInvalid(test, model.ActorEmployer, "a", cv)

	// Check returning a submitted CV against valid operations
	returnCVValid(test, model.ActorEmployer, "", cv)
	returnCVValid(test, model.ActorEmployer, "T", cv)
	returnCVValid(test, model.ActorEmployer, "Test", cv)
	returnCVValid(test, model.ActorEmployer, "Test Industry", cv)
}
