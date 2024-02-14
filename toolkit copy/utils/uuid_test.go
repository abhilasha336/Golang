package utils

import (
	"testing"
)

func TestIsValidUUID(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    string
		expected bool
	}{
		{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", true}, // Valid UUID
		{"invalid-uuid", false},                        // Invalid UUID
		{"", false},                                    // Empty input
	}

	// Run test cases
	for _, tc := range testCases {
		actual := IsValidUUID(tc.input)
		if actual != tc.expected {
			t.Errorf("For input %q, expected %t but got %t", tc.input, tc.expected, actual)
		}
	}
}
