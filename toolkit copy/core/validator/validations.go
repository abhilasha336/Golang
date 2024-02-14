package validator

import "github.com/go-playground/validator/v10"

var validations = map[string]func(fl validator.FieldLevel) bool{
	"is-number": func(fl validator.FieldLevel) bool {
		return true
	},
}
