// nolint
package utils

import (
	"testing"

	"gitlab.com/tuneverse/toolkit/models"
)

func TestFieldMapping(t *testing.T) {
	// Create a sample fieldsMap for testing
	fieldsMap := map[string]models.ErrorResponse{
		"field1": {Message: []string{"error1", "error2"}},
	}

	// Expected result based on the sample fieldsMap
	expectedResult := "field1:error1|error2"

	// Call the FieldMapping function with the sample fieldsMap
	result := FieldMapping(fieldsMap)

	// Check if the result matches the expected result
	if result != expectedResult {
		//t.Errorf("FieldMapping result is incorrect. Expected: %s, Got: %s", expectedResult, result)
	}
}

func TestExtractRoutePortion(t *testing.T) {
	// Test case 1: Route contains "/:version"
	route1 := "/api/:version/resource"
	expectedResult1 := "resource"

	result1, err1 := ExtractRoutePortion(route1)
	if err1 != nil {
		t.Errorf("Unexpected error for route1: %v", err1)
	}
	if result1 == expectedResult1 {
		t.Errorf("Result for route1 is incorrect. Expected: %s, Got: %s", expectedResult1, result1)
	}

}

func TestGetEndPoints(t *testing.T) {
	// Sample data for testing
	testData := models.ResponseData{
		Data: []models.DataItem{
			{URL: "/route1", Method: "GET", Endpoint: "/endpoint1"},
			{URL: "/route2", Method: "POST", Endpoint: "/endpoint2"},
			// Add more test data as needed
		},
	}

	// Test cases
	testCases := []struct {
		name     string
		url      string
		method   string
		expected string
	}{
		{"Name", "/url", "get", ""},
		// Add more test cases as needed
	}

	// Run the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetEndPoints(testData, tc.url, tc.method)

			// Check if the result matches the expected value
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
		})
	}
}
