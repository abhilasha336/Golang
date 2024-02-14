package validator

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationMap map[string]interface{}
type validationMapString map[string]string
type validatorFunc func(fl validator.FieldLevel) bool

var validate *validator.Validate

type FieldLevel validator.FieldLevel

type validationError struct {
	Namespace       string      `json:"namespace"` // can differ when a custom TagNameFunc is registered or
	Field           string      `json:"field"`     // by passing alt name to ReportError like below
	StructNamespace string      `json:"structNamespace"`
	StructField     string      `json:"structField"`
	Tag             string      `json:"tag"`
	ActualTag       string      `json:"actualTag"`
	Kind            interface{} `json:"kind"`
	Type            interface{} `json:"type"`
	Value           interface{} `json:"value"`
	Param           string      `json:"param"`
	FieldTag        string      `json:"field_tag"`
	Message         string      `json:"message"`
}

type Validate struct {
	validator      *validator.Validate
	EnableMessages bool
	OnlyTags       bool
	overrideRules  ValidationMap

	responseSlice map[string][]string
	reponseMap    map[string]map[string]string
}

type Options func(*validator.Validate)
type ValidationFunction func(FieldLevel) bool

// RegisterValidation
func RegisterValidation(key string, fn ValidationFunction) Options {
	return func(validate *validator.Validate) {
		validate.RegisterValidation(key, func(fl validator.FieldLevel) bool {
			return fn(fl)
		})
	}
}

// New
func New(options ...Options) *Validate {
	validate = validator.New()

	// register custom validators
	for key, fn := range bakedInValidators {
		validate.RegisterValidation(key, func(fl validator.FieldLevel) bool {
			return fn(fl)
		})
	}

	// register validation
	for _, opt := range options {
		opt(validate)
	}
	return &Validate{
		validator:      validate,
		EnableMessages: false,
		OnlyTags:       true,
		overrideRules:  nil,
	}
}

// WithRules
func (val *Validate) WithRules(rules ValidationMap) *Validate {
	val.overrideRules = rules
	return val
}

// validate the given Input
func (val *Validate) ValidateMap(data interface{}, responseBinder interface{}) (interface{}, error) {
	var errors = make(map[string]interface{})
	var (
		rules        = make(ValidationMap, 0)
		ruleMessages = make(ValidationMap, 0)
		input        = make(ValidationMap, 0)
	)
	fields, ok := data.(interface {
		Fields() ValidationMap
	})
	if !ok {
		return nil, fmt.Errorf("method `Fields` is not found on the input data")
	}

	v, ok := data.(interface {
		Rules() ValidationMap
	})
	if !ok {
		return nil, fmt.Errorf("method `Rules` is not found on the input data")
	}

	rulesMessages, ok := data.(interface {
		RulesMessages() ValidationMap
	})
	if ok {
		ruleMessages = rulesMessages.RulesMessages()
	}

	rules = v.Rules()
	// Check the user provieded rules that overrides
	if val.overrideRules != nil {
		rules = val.overrideRules
	}
	input = fields.Fields()

	errs := validate.ValidateMap(input, rules)

	if len(errs) > 0 {
		for parentKey, v := range errs {
			switch errs := v.(type) {
			case validator.ValidationErrors:
				val.prepareValidationErrors(rules, ruleMessages, parentKey, errs, errors)
			case map[string]interface{}:
				val.prepareNestedValidationErrors(rules, ruleMessages, parentKey, errs, errors)
			default:
				fmt.Printf("%T type\n", val)
			}
		}
	}

	if len(errors) > 0 {
		return errors, nil
	}

	// encode into string
	byt, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the response")
	}

	err = json.Unmarshal(byt, &responseBinder)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal the response")
	}

	return nil, nil
}

// prepareValidationErrors
// Responsible for handling validation errors within a single level of the map.
// Iterates over each validation error and delegates the preparation of error details to PrepareSlice or PrepareMap based on the OnlyTags flag
func (v *Validate) prepareValidationErrors(
	rules map[string]interface{},
	ruleMessages map[string]interface{},
	parentKey string,
	errs validator.ValidationErrors,
	errors map[string]interface{},
) {
	for _, err := range errs {
		if v.OnlyTags {
			v.PrepareSlice(rules, parentKey, err, errors)
		} else {
			v.PrepareMap(rules, ruleMessages, parentKey, err, errors)
		}
	}
}

// prepareNestedValidationErrors
// Manages nested validation errors within a map
// Iterates over each nested map, preparing error details recursively using the prepareValidationErrors method
func (v *Validate) prepareNestedValidationErrors(
	rules map[string]interface{},
	ruleMessages map[string]interface{},
	parentKey string,
	errs map[string]interface{},
	errors map[string]interface{},
) {
	for key, err := range errs {
		nestedErrs := make(map[string]interface{})

		v.prepareValidationErrors(
			rules[parentKey].(map[string]interface{}),
			ruleMessages,
			key,
			err.(validator.ValidationErrors),
			nestedErrs,
		)

		if len(nestedErrs) > 0 {
			if _, ok := errors[parentKey]; !ok {
				errors[parentKey] = make(map[string]interface{})
			}
			// copy map items
			maps.Copy(errors[parentKey].(map[string]interface{}), nestedErrs)
		}
	}
}

// PrepareSlice
// Placeholder method for handling validation errors in a slice or array
func (vl *Validate) PrepareSlice(
	rules ValidationMap, // rules map
	key string,
	err validator.FieldError,
	errors map[string]interface{},
) {

	if _, ok := errors[key]; !ok {
		errors[key] = []string{}
	}

	splits := strings.Split(rules[key].(string), ",")
	var canAppend bool
	for _, split := range splits {
		if split == key {
			canAppend = true
		}
		// check the string can be split based on the =
		nested := strings.Split(split, "=")
		if len(nested) > 0 && nested[0] == err.Tag() {
			canAppend = true
		}
		if canAppend {
			errors[key] = append(errors[key].([]string), split)
		}
	}
}

// PrepareMap
// Placeholder method for handling validation errors in a map
func (vl *Validate) PrepareMap(
	rules ValidationMap,
	ruleMessages map[string]interface{},
	key string,
	err validator.FieldError,
	errors map[string]interface{},
) {
	if _, ok := errors[key]; !ok {
		errors[key] = map[string]string{}
	}

	// check the field contains
	if message, ok := ruleMessages[err.Tag()].(string); ok {
		message = strings.ReplaceAll(message, "{field_name}", (key))
		message = strings.ReplaceAll(message, "{param}", err.Param())
		errors[key].(map[string]string)[err.Tag()] = message
	}

	// check the other validations that may fail
	splits := strings.Split(rules[key].(string), ",")
	var canAppend bool
	for _, split := range splits {
		// check the string can be split based on the =
		nested := strings.Split(split, "=")
		if len(nested) > 0 && nested[0] == err.Tag() {
			canAppend = true
		}

		if canAppend {
			if message, ok := ruleMessages[nested[0]].(string); ok {
				message = strings.ReplaceAll(message, "{field_name}", (key))
				if len(nested) > 1 {
					message = strings.ReplaceAll(message, "{param}", nested[1])
				}
				errors[key].(map[string]string)[nested[0]] = message
			}
		}
	}
}
