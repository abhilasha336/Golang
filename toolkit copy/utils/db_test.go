// nolint
package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildSearchQuery(t *testing.T) {
	// Test cases for BuildSearchQuery
	tests := []struct {
		value    string
		fields   []string
		expected string
	}{
		{"keyword", []string{"field1", "field2"}, " field1 ILIKE '%keyword%' OR field2 ILIKE '%keyword%' "},
		{"hello", []string{"title", "author"}, " title ILIKE '%hello%' OR author ILIKE '%hello%' "},
		{"john", []string{"name"}, " name ILIKE '%john%' "},
	}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			result := BuildSearchQuery(test.value, test.fields...)
			require.Equal(t, test.expected, result)
		})
	}
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		page     int32
		limit    int32
		expected string
	}{
		{1, 10, "LIMIT 10 OFFSET 0"},
		{2, 10, "LIMIT 10 OFFSET 10"},
		{3, 20, "LIMIT 20 OFFSET 40"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Page%dLimit%d", test.page, test.limit), func(t *testing.T) {
			result := CalculateOffset(test.page, test.limit)
			require.Equal(t, test.expected, result)
		})
	}
}

// TestPreparePlaceholders is a unit test for the PreparePlaceholders function.
// It tests the generation of SQL placeholders for a given number and compares the result with the expected placeholders string.
func TestPreparePlaceholders(t *testing.T) {
	// Test cases with input n and expected placeholders string
	testCases := []struct {
		n              int
		expectedResult string
	}{
		{0, ""},
		{1, "$1"},
		{3, "$1,$2,$3"},
		{5, "$1,$2,$3,$4,$5"},
	}

	// Iterate through test cases and run the test for each case
	for _, tc := range testCases {
		result := PreparePlaceholders(tc.n)

		// Check if the result matches the expected placeholders string
		if result != tc.expectedResult {
			t.Errorf("For n=%d, expected placeholders '%s', but got '%s'", tc.n, tc.expectedResult, result)
		}
	}
}
