package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var (
	bakedInValidators = map[string]validatorFunc{
		"isString": isString,
	}
)

// checkt the given input is a string
func isString(fl validator.FieldLevel) bool {
	return fl.Field().Kind() == reflect.String
}
