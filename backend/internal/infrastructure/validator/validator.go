package validator

import "github.com/go-playground/validator/v10"

type EchoValidator struct {
	validator *validator.Validate
}

func NewEchoValidator() *EchoValidator {
	return &EchoValidator{
		validator: validator.New(),
	}
}

func (v *EchoValidator) Validate(i any) error {
	return v.validator.Struct(i)
}
