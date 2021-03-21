package validation

import (
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

func validateMobilePhone(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^08\d{7,12}$`)
	return regex.MatchString(fl.Field().String())
}

func validateDateOnly(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^\d{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$`)
	return regex.MatchString(fl.Field().String())
}
