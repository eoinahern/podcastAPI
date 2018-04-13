package util

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

//PasswordEncryptUtilInt interface
type PasswordEncryptUtilInt interface {
	Encrypt(password string) string
	CheckSame(DBpassword string, sentPassword string) bool
}

//PasswordEncryptUtil : used to encrypt password and check if its valid
type PasswordEncryptUtil struct {
}

//Encrypt : encrypts my password
func (p *PasswordEncryptUtil) Encrypt(password string) string {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}

//CheckSame : checks if encypted passwords are the same
func (p *PasswordEncryptUtil) CheckSame(DBpassword string, sentPassword string) bool {

	if err := bcrypt.CompareHashAndPassword([]byte(DBpassword), []byte(sentPassword)); err != nil {
		return false
	}

	return true
}
