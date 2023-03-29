package cvalidator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	regex, _ := regexp.Compile(`^\+\d{7,11}$`)
	result := regex.MatchString(phone)
	return result
}
