package utilities

import (
	"localization/internal/entities"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// parseFields parses the provided field string and extracts error messages based on the field structure.
// It returns a map of field errors and a boolean indicating if any errors were found.
// params
// @field		- Field is of the form fieldname:value (address:Valid)
// @errorsValue	- Error structure from database
// @endpoint 	- Endpoint names
// @method 		- Method of the endpoint (get,post..)
func ExtractFieldErrors(field string, errorsValue bson.M, endpoint string, method string) (map[string]any, bool) {
	fieldErrors := make(map[string]any)
	pairs := strings.Split(field, ",")

	for _, pair := range pairs {
		parts := strings.Split(pair, ":")

		fieldName := parts[0]
		errorKeys := strings.Split(parts[1], "|")

		// Check if the registration map exists and has the specified field
		endpointMap, endpointExists := errorsValue[endpoint].(bson.M)
		if !endpointExists {
			continue
		}

		methodMap, methodExists := endpointMap[method].(bson.M)
		if !methodExists {
			continue
		}

		fieldErrorMap, fieldExists := methodMap[fieldName].(bson.M)

		if !fieldExists {
			continue
		}

		fieldError := make(map[string]any)
		for _, key := range errorKeys {
			if errorMsg, ok := fieldErrorMap[key].(string); ok {
				fieldError[key] = errorMsg
			}
		}

		if len(fieldError) > 0 {
			fieldErrors[fieldName] = fieldError
		}
	}

	return fieldErrors, len(fieldErrors) > 0
}

// GenerateErrorResponse is used to generate error response in appropriate format
// params
// @message		- Message to be display
// @errorCode	- Errorcode to display
// @errorsValue	- Json error structure
func ErrorResponseFormat(message string, errorCode any, errorsValue map[string]any, key string) any {
	response := entities.ErrorResponse{
		Message:   message,
		ErrorCode: errorCode,
		Errors:    map[string]any{key: errorsValue},
	}

	return response
}
