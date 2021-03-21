package validation

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"

)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {

	if err := v.validator.RegisterValidation("mobile_phone", validateMobilePhone); err != nil {
		log.Error("failed register validation mobile_phone")
		return err
	}
	if err := v.validator.RegisterValidation("date_only", validateDateOnly); err != nil {
		log.Error("failed register validation date_only")
		return err
	}

	return v.validator.Struct(i)
}
