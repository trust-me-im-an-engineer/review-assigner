package validator

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	return &Validator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Struct validates a struct based on "validate" tags.
func (v *Validator) Struct(s any) error {
	return v.validator.Struct(s)
}
