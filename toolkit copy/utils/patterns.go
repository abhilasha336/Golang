package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.com/tuneverse/toolkit/consts"
)

// GetRegexPatternFromFile retrieves a regex pattern from a JSON file based on the provided type (key).
// Parameters:
//   - filePath: The path to the JSON file containing regex patterns.
//   - regexType: The type (key) of the regex pattern to retrieve from the JSON file.
//
// Returns:
//   - string: The regex pattern corresponding to the provided type.
//   - error: Any error that occurred during the process, or nil if successful.
func GetRegexPatternFromFile(filePath, regexType string) (string, error) {
	var regex string

	resp, err := http.Get(filePath)
	if err != nil {
		fmt.Println("Error fetching file:", err)
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}

	// Print the JSON content as a reconstructed JSON string
	reconstructedJSON, _ := json.Marshal(data)

	// Create a map to hold the JSON data
	patterns := make(map[string]string)

	// Unmarshal JSON data into the map
	err = json.Unmarshal(reconstructedJSON, &patterns)
	fmt.Println("")
	if err != nil {
		log.Errorf("Error while reading the file %v", err)
		return "", err
	}

	if value, ok := patterns[regexType]; ok {
		regex = value
	}
	return regex, err
}

// CheckPattern checks if the given input string matches a regex pattern identified by the provided key.
// It reads the regex pattern from a JSON file using GetRegexPatternFromFile.
// Parameters:
//   - input: The string to be checked against the regex pattern.
//   - key: The identifier for the desired regex pattern (e.g., "email", "empty", etc.).
//
// Returns:
//   - bool: true if the input matches the regex pattern, false otherwise.
//   - error: Any error that occurred during the process, or nil if successful.
func CheckPattern(input string, key string) (bool, error) {
	// Get the regex pattern based on the provided key
	pattern, err := GetRegexPatternFromFile(consts.Validations.JsonPathKey, key)
	if err != nil {
		return false, err
	}

	// If the key is for "trackextension," modify the input to remove the dot from the extension
	if key == consts.Validations.TrackExtensionKey {
		lastDotIndex := strings.LastIndex(input, ".")

		// Check if a dot is found and it's not the last character
		if lastDotIndex != -1 && lastDotIndex < len(input)-1 {
			// Extract the extension (including the dot)
			extension := input[lastDotIndex:]
			input = extension[1:]
		}
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	// Test whether the given string matches the compiled regex
	matcher := regex.MatchString(input)

	if key == consts.Validations.SpecialCharKey || key == consts.Validations.URLKey || key == consts.Validations.DateKey {
		matcher = !matcher
	}
	return matcher, nil
}
