package validation

import "strings"

//EmailValidationInt interface
type EmailValidationInt interface {
	CheckEmailValid(email string) bool
}

//EmailValidation : validates email string. could be updated
type EmailValidation struct {
}

//CheckEmailValid : checkemail is over 10 chars and contains @ and . chars
func (e *EmailValidation) CheckEmailValid(email string) bool {

	if len(email) > 10 && strings.Contains(email, "@") && strings.Contains(email, ".") {
		return true
	}

	return false

}
