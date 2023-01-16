package util

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func InitValidate() {
	Validate = validator.New()
}
