package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConvertDate(t *testing.T) {
	testCases := []struct {
		date            string
		format          string
		expectedRFC3339 string
		expectedError   string
	}{
		{"2023-09-19", "2006-01-02", "2023-09-19", ""},
		{"", "2006-01-02", "2023-09-19", "date is not provided"},
		{"26-10-1999", "02/01/2006", "", "parsing time"},
		{"26/10/1999", "02/01/2006", "26/10/1999", ""},
		{"2023-09-19T12:34:56Z", time.RFC3339, "2023-09-19T12:34:56Z", ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.date, func(t *testing.T) {
			result, _, err := ConvertDate(testCase.date, testCase.format)
			if testCase.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), testCase.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedRFC3339, result)
			}
		})
	}
}

func TestAddHours(t *testing.T) {
	baseTime, _ := time.Parse(time.RFC3339, "2023-09-19T12:00:00Z")

	testCases := []struct {
		hours          int
		format         string
		expectedResult string
	}{
		{1, time.RFC3339, "2023-09-19T13:00:00Z"},
		{-1, time.RFC3339, "2023-09-19T11:00:00Z"},
		{0, time.RFC3339, "2023-09-19T12:00:00Z"},
		{2, "2006-01-02 15:04:05", "2023-09-19 14:00:00"},
		{-2, "2006-01-02 15:04:05", "2023-09-19 10:00:00"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedResult, func(t *testing.T) {
			result := AddHours(baseTime, testCase.hours, testCase.format)
			require.Equal(t, testCase.expectedResult, result)
		})
	}
}

func TestFormatDateTime(t *testing.T) {
	// Define test cases
	testCases := []struct {
		inputTime      string
		expectedOutput string
		outputFormat   string
		expectedError  bool
	}{

		{
			inputTime:      "invalidTime",
			expectedOutput: "",
			outputFormat:   "02-Jan-2006 03:04 PM",
			expectedError:  true, // This test case is expected to return an error
		},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tc := range testCases {
		result, err := FormatDateTime(tc.inputTime, tc.outputFormat)

		// Check if the error matches the expected result
		if (err != nil) != tc.expectedError {
			t.Errorf("Test failed for inputTime: %s, outputFormat: %s. Expected error: %v, got error: %v", tc.inputTime, tc.outputFormat, tc.expectedError, err)
			continue
		}

		// Check if the result matches the expected output
		if result != tc.expectedOutput {
			t.Errorf("Test failed for inputTime: %s, outputFormat: %s. Expected: %s, got: %s", tc.inputTime, tc.outputFormat, tc.expectedOutput, result)
		}
	}
}
