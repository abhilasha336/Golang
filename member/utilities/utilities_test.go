package utilities_test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"member/utilities"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	cryptoHash "gitlab.com/tuneverse/toolkit/utils/crypto"
)

// TestValidateEmail is a unit test for the ValidateEmail function in the 'utilities' package.
// It tests the validation of email addresses, including valid, empty, and invalid formats.
func TestValidateEmail(t *testing.T) {
	// Test case for a valid email
	validEmail := "test@example.com"
	if !utilities.ValidateEmail(validEmail) {
		t.Errorf("Expected true for a valid email, but got false")
	}

	// Test case for an empty email
	emptyEmail := ""
	if utilities.ValidateEmail(emptyEmail) {
		t.Errorf("Expected false for an empty email, but got true")
	}

	// Test case for an invalid email format
	invalidEmail := "invalid-email"
	if utilities.ValidateEmail(invalidEmail) {
		t.Errorf("Expected false for an invalid email format, but got true")
	}
}

// TestValidateName is a unit test for the ValidateName function in the 'utilities' package.
// It tests the validation of names, checking for the presence of special characters.
func TestValidateName(t *testing.T) {
	// Test case for a valid name
	validName := "John Doe"
	if !utilities.ValidateName(validName) {
		t.Errorf("Expected true for a valid name, but got false")
	}

	// Test case for a name with special characters
	invalidName := "John-Doe"
	if utilities.ValidateName(invalidName) {
		t.Errorf("Expected false for a name with special characters, but got true")
	}
}

// TestIsValidZIPCode is a unit test for the IsValidZIPCode function in the 'utilities' package.
// It tests the validation of ZIP codes, checking for valid formats.
func TestIsValidZIPCode(t *testing.T) {
	// Test case for a valid ZIP code
	validZIPCode := "12345"
	if !utilities.IsValidZIPCode(validZIPCode) {
		t.Errorf("Expected true for a valid ZIP code, but got false")
	}

	// Test case for an invalid ZIP code format
	invalidZIPCode := "ABCD"
	if utilities.IsValidZIPCode(invalidZIPCode) {
		t.Errorf("Expected false for an invalid ZIP code format, but got true")
	}
}

// TestIsValidPhoneNumber is a unit test for the IsValidPhoneNumber function in the 'utilities' package.
// It tests the validation of phone numbers, considering the phone number format and region.
func TestIsValidPhoneNumber(t *testing.T) {
	// Test case for a valid phone number in Indian format
	validPhoneNumber := "+911234567890" // +91 is the country code for India
	defaultRegion := "IN"
	if !utilities.IsValidPhoneNumber(validPhoneNumber, defaultRegion) {
		t.Errorf("Expected true for a valid Indian phone number, but got false")
	}

	// Test case for an invalid phone number
	invalidPhoneNumber := "invalid"
	if utilities.IsValidPhoneNumber(invalidPhoneNumber, defaultRegion) {
		t.Errorf("Expected false for an invalid phone number, but got true")
	}
}

// TestIsValidUUID is a unit test for the IsValidUUID function in the 'utilities' package.
// It tests the validation of UUID strings, checking if they are in the correct format.
func TestIsValidUUID(t *testing.T) {
	// Test cases with input UUID strings and expected results
	testCases := []struct {
		id            string
		expectedValid bool
	}{
		{"123e4567-e89b-12d3-a456-426655440000", true},
		{"invalid-uuid", false},
		{"123e4567-e89b-12d3-a456-42665544", false},
	}

	// Iterate through test cases and run the test for each case
	for _, tc := range testCases {
		isValid := utilities.IsValidUUID(tc.id)

		// Check the result against the expected value
		assert.Equal(t, tc.expectedValid, isValid)
	}
}

// TestHashPassword is a unit test for the HashPassword function in the 'utilities' package.
// It tests the hashing of passwords and compares the result with the expected hash values.
func TestHashPassword(t *testing.T) {
	// Test cases with input passwords and expected hash values
	testCases := []struct {
		password     string
		expectedHash string
	}{
		{"Password1@", "e0b30e8c74f9f46e442dbd4521422f92"},
		{"AnotherPassword1@", "20880c25864da6b90feda23d4d080965"},
	}

	// Iterate through test cases and run the test for each case
	for _, tc := range testCases {
		hash, err := cryptoHash.Hash(tc.password)
		if err != nil {
			fmt.Printf("Error while hashing the password: %v\n", err)
		}
		// Compute the expected hash using MD5
		expectedMD5Hash := md5.Sum([]byte(tc.password))
		expectedHexHash := hex.EncodeToString(expectedMD5Hash[:])

		// Check if the result matches the expected hash value
		if strings.Compare(expectedHexHash, hash) != 0 {
			t.Errorf("For password '%s', expected hash '%s', but got '%s'", tc.password, expectedHexHash, hash)
		}
	}
}

// TestValidateNameLength is a unit test for the ValidateNameLength function in the 'utilities' package.
// It tests the validation of name lengths and compares the result with the expected validation results.
func TestValidateMaximumNameLength(t *testing.T) {
	// Test cases with input names and expected validation results
	testCases := []struct {
		name            string
		expectedIsValid bool
	}{
		{"JohnDoe", true}, // Valid name within the length limit
		{"ThisIsAReallyLongNameThatExceedsTheMaximumLengthLimiThisShouldFail", false}, // Name exceeds the length limit
		{"", true}, // Empty name is considered valid
	}

	// Iterate through test cases and run the test for each case
	for _, tc := range testCases {
		isValid := utilities.ValidateMaximumNameLength(tc.name)

		// Check if the result matches the expected validation result
		if isValid != tc.expectedIsValid {
			t.Errorf("For name '%s', expected validation result '%v', but got '%v'", tc.name, tc.expectedIsValid, isValid)
		}
	}
}
