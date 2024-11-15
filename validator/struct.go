package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	service *validator.Validate
}

func NewStructValidator() *Validator {
	return &Validator{service: validator.New()}
}

func (v *Validator) Validate(req interface{}) error {
	return v.service.Struct(req)
}
