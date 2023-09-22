package infrastructure

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func GetValidator() *validator.Validate {
	if validate == nil {
		validate = validator.New()

		validate.RegisterValidation("institutional_email", func(fl validator.FieldLevel) bool {
			email := fl.Field().String()
			return strings.HasSuffix(email, "@upb.edu.co")
		})
	}

	return validate
}
