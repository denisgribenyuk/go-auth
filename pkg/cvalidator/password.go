package cvalidator

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	minLength := 8
	maxLength := 255

	acceptableSymbols := make(map[rune]bool)
	for r := 'a'; r < 'z'; r++ {
		acceptableSymbols[r] = true
	}
	for r := 'A'; r < 'Z'; r++ {
		acceptableSymbols[r] = true
	}
	for r := '0'; r < '9'; r++ {
		acceptableSymbols[r] = true
	}
	for _, r := range "*.!@$%^&(){}[]:;<>,.?/~_+-=|\\/#`" {
		acceptableSymbols[r] = true
	}

	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, c := range password {
		if !acceptableSymbols[rune(c)] {
			return false
		}

		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true

		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	length := len(password) >= minLength && len(password) <= maxLength

	return hasNumber && hasUpperCase && hasLowercase && hasSpecial && length
}
