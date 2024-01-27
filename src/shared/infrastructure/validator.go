package infrastructure

import (
	"regexp"
	"strings"
	"time"

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

		validate.RegisterValidation("secure_password", func(fl validator.FieldLevel) bool {
			password := fl.Field().String()
			hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
			hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
			hasSpecialCharacter := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
			return hasLetter && hasNumber && hasSpecialCharacter
		})

		validate.RegisterValidation("RFC3339_date", func(fl validator.FieldLevel) bool {
			date := fl.Field().String()
			_, err := time.Parse(time.RFC3339, date)
			return err == nil
		})
	}

	return validate
}
