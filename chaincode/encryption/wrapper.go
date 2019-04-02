package encryption

import (
	// "fmt"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
)
// EncCC example simple Chaincode implementation of a chaincode that uses encryption/signatures
type EncCC struct {
	bccspInst bccsp.BCCSP
}
// Encrypter encrypts the data and writes to the ledger
func Encrypter(stub shim.ChaincodeStubInterface, valueAsBytes []byte) []byte {
	factory.InitFactories(nil)
	encCC := EncCC{factory.GetDefault()}
	encKey := make([]byte, 32)
	iv := make([]byte, 16)
	ent, _ := entities.NewAES256EncrypterEntity("ID", encCC.bccspInst, encKey, iv)

	encCC.bccspInst.GetKey()

	/*	if err != nil {
		return
	}*/
	return encrypt(ent, valueAsBytes)
}
// Decrypter decrypts the data and writes to the ledger
func Decrypter(stub shim.ChaincodeStubInterface, cipherText []byte) []byte {
	factory.InitFactories(nil)
	encCC := EncCC{factory.GetDefault()}
	// fmt.Println(encCC)
	decKey := make([]byte, 32)
	iv := make([]byte, 16)
	ent, _ := entities.NewAES256EncrypterEntity("ID", encCC.bccspInst, decKey, iv)
	/*if err != nil {
		return nil, err
	}*/
	return decrypt(ent, cipherText)
}