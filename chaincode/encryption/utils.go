package encryption

import (
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
)
// getStateAndDecrypt retrieves the value associated to key,
// decrypts it with the supplied entity and returns the result
// of the decryption
func decrypt(ent entities.Encrypter, cipherText []byte) []byte {
	// at first we retrieve the ciphertext from the ledger
	if len(cipherText) == 0 {
		return []byte("")
	}

	result, err := ent.Decrypt(cipherText)

	if err != nil {
		return []byte("")
	}

	return result
}
// encryptAndPutState encrypts the supplied value using the
// supplied entity and puts it to the ledger associated to
// the supplied KVS key
func encrypt(ent entities.Encrypter, value []byte) []byte {
	// at first we use the supplied entity to encrypt the value
	cipherText, err := ent.Encrypt(value)
	if err != nil {
		return []byte("")
	}
	return cipherText
}