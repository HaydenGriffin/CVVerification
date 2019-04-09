package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

//Hash implements root.Hash
type Crypto struct{}

//Generate a salted hash for the input string
func GenerateFromByte(b []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

//Compare string to generated hash
func Compare(hash string, s string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))

	return err == nil
}