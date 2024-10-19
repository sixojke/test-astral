package v1

import (
	"regexp"
	"unicode"
)

func validateLogin(login string) bool {
	if len(login) < 8 || len(login) > 32 {
		return false
	}

	return regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(login)
}

func validatePassword(password string) bool {
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSymbol := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case !unicode.IsLetter(char) && !unicode.IsDigit(char):
			hasSymbol = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSymbol && len(password) > 8
}
