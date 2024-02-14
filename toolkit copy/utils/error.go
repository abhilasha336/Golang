package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/models"
)

// Getting error codes from context
func GetErrorCodes(ctx *gin.Context) map[string]any {
	contextError, _ := GetContext[map[string]any](ctx, consts.ContextErrorResponses)
	return contextError
}

func GenerateErrorResponse(errorDetails models.ErrorDetails, fieldErrors map[string]interface{}, serviceCode map[string]string, helpLink string) ([]map[string]interface{}, models.ErrorDetails) {
	var errorsData []map[string]interface{}
	response := make([]map[string]interface{}, 0)
	if fieldErrors != nil {
		for field, messages := range fieldErrors {
			errorData := map[string]interface{}{
				"field":      field,
				"message":    messages,
				"error_code": serviceCode[field], // Currently hard coded, will be changed after verfication of service based error codes
				"help":       helpLink,
			}
			errorsData = append(errorsData, errorData)
		}
		response = errorsData
	}

	// Set 'data' field only if there are errors

	if fieldErrors == nil {
		errorData := map[string]any{
			"message":    errorDetails.Message,
			"error_code": errorDetails.Code, // Currently hard coded, will be changed after verfication of service based error codes
		}
		errorsData = append(errorsData, errorData)
	}
	return response, errorDetails
}

// Function to retrieve error functions
func ParseFields(ctx *gin.Context, errorType string, fields string, errorsValue map[string]any, endpoint string, method string,
	serviceCode map[string]string, helpLink string) ([]map[string]interface{}, bool, float64, models.ErrorDetails) {
	var responseJSON []map[string]interface{}
	var errorDetails models.ErrorDetails
	// Access the errors value
	errorsMap, exists := errorsValue[consts.Errors].(map[string]any)
	if !exists {
		return nil, false, 0, errorDetails
	}

	// Access the endpoint-specific error map
	allErrorMap, exists := errorsMap[consts.AllError].(map[string]any)
	if !exists {
		return nil, false, 0, errorDetails
	}
	var ok bool
	message := allErrorMap[errorType].(map[string]any)
	typeErrMsg := message["message"]
	typeErrCode := message["errorCode"]
	errorDetails.Code, ok = typeErrCode.(float64)
	// x, ok := typeErrCode.(float64)

	if !ok {

		fmt.Println("code converting error")
	}

	errorDetails.Message, ok = typeErrMsg.(string)
	if !ok {
		fmt.Println("message converting error")
	}
	// other error msg handling.
	if errorType != consts.ValidationErr {
		errOthers := allErrorMap[errorType].(map[string]any)

		errCodeVal := errOthers[consts.ErrorCode]
		errVal := errCodeVal.(float64)
		responseJSON, _ = GenerateErrorResponse(errorDetails, nil, nil, helpLink)
		return responseJSON, true, errVal, errorDetails
	}

	validationErrorMap, exists := allErrorMap[consts.ValidationErr].(map[string]any)
	errorCodeVal := validationErrorMap[consts.ErrorCode]
	if !exists {
		return nil, false, errorCodeVal.(float64), errorDetails
	}

	insideErrors, _ := validationErrorMap[consts.Errors].(map[string]any)

	pairs := strings.Split(fields, ",")
	fieldErrors := make(map[string]any)

	for _, pair := range pairs {
		parts := strings.Split(pair, ":")

		fieldName := parts[0]
		errorKeys := strings.Split(parts[1], "|")

		fieldErrorMap, fieldExists := insideErrors[endpoint].(map[string]any)
		if !fieldExists {
			continue
		}

		methodErrorMap, methodExists := fieldErrorMap[method].(map[string]any)
		if !methodExists {
			continue
		}

		fieldErrorMessages := make([]string, 0)
		fieldSpecificErrors, fieldSpecificExists := methodErrorMap[fieldName].(map[string]any)
		if !fieldSpecificExists {
			continue
		}

		for _, value := range errorKeys {
			if errorValue, ok := fieldSpecificErrors[value].(string); ok {
				fieldErrorMessages = append(fieldErrorMessages, errorValue)
			}
		}

		if len(fieldErrorMessages) > 0 {
			fieldErrors[fieldName] = fieldErrorMessages
		}
	}
	responseJSON, errDet := GenerateErrorResponse(errorDetails, fieldErrors, serviceCode, helpLink)
	return responseJSON, true, errorCodeVal.(float64), errDet
}

// To retrieve errors
func GetErrorCode(data map[string]interface{}, endpoint string, method string, field string, key string) (string, error) {

	endpointData, ok := data[endpoint].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid endpoint data")
	}

	methodData, ok := endpointData[method].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid method data")
	}
	fieldData, ok := methodData[field].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid field data")
	}
	keyData, ok := fieldData[key].(string)
	if !ok {
		return "", errors.New("invalid key data")
	}

	return keyData, nil
}
