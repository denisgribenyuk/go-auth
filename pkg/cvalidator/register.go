package cvalidator

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var V *validator.Validate

func Register() {
	fmt.Println("Registering custom validators")
	// register custom validators
	V = validator.New()

	if validatorEngine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for _, vv := range []struct {
			tag      string
			validate validator.Func
		}{
			{tag: "password", validate: validatePassword},
			{tag: "phone", validate: validatePhone},
		} {
			fmt.Printf("Registering %s validator\n", vv.tag)
			err := validatorEngine.RegisterValidation(vv.tag, vv.validate)
			if err != nil {
				fmt.Printf("Error registering %s validator", vv.tag)
			}
			err = V.RegisterValidation(vv.tag, vv.validate)
			if err != nil {
				fmt.Printf("Error registering %s validator", vv.tag)
			}
		}
	}
}
