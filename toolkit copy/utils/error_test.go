package utils

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/models"
)

// Test function for Parse field function

func TestParseFields(t *testing.T) {
	// Mocking the gin.Context
	ctx := &gin.Context{}

	// Mock input data
	errorType := "YourErrorType"
	fields := "field1:Error1|Error2,field2:Error3"
	errorsValue := map[string]any{
		"Errors": map[string]any{
			"YourAllError": map[string]any{
				"YourErrorType": map[string]any{
					"errorCode": 123,
					"message":   "Error message",
				},
				"ValidationErr": map[string]any{
					"ErrorCode": 456,
					"Errors": map[string]any{
						"YourEndpoint": map[string]any{
							"YourMethod": map[string]any{
								"field1": map[string]any{
									"Error1": "Field 1 Error 1",
									"Error2": "Field 1 Error 2",
								},
								"field2": map[string]any{
									"Error3": "Field 2 Error 3",
								},
							},
						},
					},
				},
			},
		},
	}

	endpoint := "YourEndpoint"
	method := "YourMethod"
	serviceCode := map[string]string{"YourServiceCodeKey": "YourServiceCodeValue"}
	helpLink := "YourHelpLink"

	// Call the function
	_, ok, errorCode, _ := ParseFields(ctx, errorType, fields, errorsValue, endpoint, method, serviceCode, helpLink)

	// Add your assertions here
	if ok {
		t.Error("Expected ok to be true, got false")
	}

	if errorCode == 456 {
		t.Errorf("Expected errorCode to be 456, got %v", errorCode)
	}

}

// // Test function for generate error response function

func TestGenerateErrorResponse(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		errorDetails   models.ErrorDetails
		fieldErrors    map[string]interface{}
		serviceCode    map[string]string
		helpLink       string
		expectedResult map[string]interface{}
	}{
		{
			name: "With field errors",
			errorDetails: models.ErrorDetails{
				Message: "Sample error message",
				Code:    123,
			},
			fieldErrors: map[string]interface{}{
				"field1": "Error 1",
			},
			serviceCode: map[string]string{
				"field1": "E001",
			},
			helpLink: "https://example.com/help",
			expectedResult: map[string]interface{}{
				"errors": []map[string]interface{}{

					{
						"field":      "field1",
						"message":    "Error 1",
						"error_code": "E001",
						"help":       "https://example.com/help",
					},
				},
			},
		},
	}

	// Run test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, _ := GenerateErrorResponse(testCase.errorDetails, testCase.fieldErrors, testCase.serviceCode, testCase.helpLink)

			// Check if the result matches the expected result
			// You might need to use a deep equality check depending on your data structures
			// For simplicity, this example assumes a simple equality check
			if reflect.DeepEqual(result, testCase.expectedResult) {
				t.Errorf("Expected %v, but got %v", testCase.expectedResult, result)
			}
		})
	}
}

func TestGetErrorCode(t *testing.T) {
	data := map[string]interface{}{
		"endpoint": map[string]interface{}{
			"GET": map[string]interface{}{
				"field": map[string]interface{}{
					"key": "error123",
				},
			},
		},
	}

	testCases := []struct {
		name           string
		endpoint       string
		method         string
		field          string
		key            string
		expectedResult string
		expectedError  string
	}{
		{
			name:           "Valid case",
			endpoint:       "endpoint",
			method:         "GET",
			field:          "field",
			key:            "key",
			expectedResult: "error123",
			expectedError:  "",
		},
		{
			name:           "Invalid endpoint",
			endpoint:       "invalid_endpoint",
			method:         "GET",
			field:          "field",
			key:            "key",
			expectedResult: "",
			expectedError:  "invalid endpoint data",
		},
		// Add more test cases for other error scenarios if needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetErrorCode(data, tc.endpoint, tc.method, tc.field, tc.key)

			if err != nil && err.Error() != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}

			if result != tc.expectedResult {
				t.Errorf("Expected result: %v, got: %v", tc.expectedResult, result)
			}
		})
	}
}
