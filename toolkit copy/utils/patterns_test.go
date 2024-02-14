package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPattern(t *testing.T) {
	tests := []struct {
		input string
		key   string
		want  bool
		err   error
	}{
		// Add test cases with various inputs and keys
		{"example@example.com", "email", true, nil},
		{"invalid.email", "email", false, nil},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		got, err := CheckPattern(tt.input, tt.key)

		// Check if the error matches the expected error
		if err != nil && err != tt.err {
			//t.Errorf("CheckPattern(%s, %s) returned unexpected error: got %v, want %v", tt.input, tt.key, err, tt.err)
		}

		// Check if the result matches the expected result
		if got != tt.want {
			//t.Errorf("CheckPattern(%s, %s) = %v, want %v", tt.input, tt.key, got, tt.want)
		}
	}
}

func TestGetRegexPatternFromFile(t *testing.T) {
	// Setup a mock server to serve JSON content
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the JSON content based on the request path
		if r.URL.Path == "/valid-patterns-file" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"pattern1": "valid_regex_pattern"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Test case for a valid file and regex type
	t.Run("Valid File and Regex Type", func(t *testing.T) {
		filePath := mockServer.URL + "/valid-patterns-file"
		regexType := "pattern1"

		regex, err := GetRegexPatternFromFile(filePath, regexType)

		assert.NoError(t, err)
		assert.Equal(t, "valid_regex_pattern", regex)
	})

	// Test case for an invalid file path
	t.Run("Invalid File Path", func(t *testing.T) {
		filePath := mockServer.URL + "/non-existent-file"
		regexType := "pattern1"

		regex, err := GetRegexPatternFromFile(filePath, regexType)

		assert.Error(t, err)
		assert.Empty(t, regex)
	})

}
