package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword ...
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// CheckPasswordHash ...
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateHash for general purpose
func GenerateHash(str string, ts string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str+ts), 5)
	return string(bytes), err
}
