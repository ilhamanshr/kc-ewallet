package strhelper

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash hash a string
func Hash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckHash check hashed string with original string
func CheckHash(hashedStr string, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
	return err == nil
}
