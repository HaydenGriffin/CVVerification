package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

//Hash implements root.Hash
type Crypto struct{}

//Generate a salted hash for the input string
func Generate(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

//Compare string to generated hash
func Compare(hash string, s string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))

	return err == nil
}