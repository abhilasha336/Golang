package utils

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ttacon/libphonenumber"
	"gitlab.com/tuneverse/toolkit/consts"
)

// IsEmpty Checks whether a string is empty or not
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsValidIPAddress checks if the given string is a valid IP address (IPv4 or IPv6).
func IsValidIPAddress(ipAddress string) bool {
	// Check if the address is IPv4 or IPv6 by attempting to parse it.
	ipAddress = strings.TrimSpace(ipAddress)
	if ip := net.ParseIP(ipAddress); ip == nil {
		return false
	}
	return true
}

// function to validate an url
func IsValidURL(str string) bool {
	var rxURL = regexp.MustCompile(consts.URL)

	if str == "" || utf8.RuneCountInString(str) >= consts.MaxURLRuneCount || len(str) <= consts.MinURLRuneCount || strings.HasPrefix(str, ".") {
		return false
	}
	strTemp := str
	if strings.Contains(str, ":") && !strings.Contains(str, "://") {
		strTemp = "http://" + str
	}
	u, err := url.Parse(strTemp)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return rxURL.MatchString(str)
}

// function to validate an email
func IsValidEmail(email string) bool {
	rxEmail := regexp.MustCompile(consts.Email)
	return rxEmail.MatchString(email)
}

// function to validate a name field
// It checks the length of the string is between minimum and maximum length
// also checks whether it contains any special characters or symbols
func IsValidName(name string) (error, bool) {
	// Check the length of the name
	if len(name) > consts.MaxLength || len(name) < consts.MinLength {
		return fmt.Errorf("the limit of Name field is between  %d and %d", consts.MinLength, consts.MaxLength), false
	}

	// Check if the name contains only alphabetic characters
	for _, char := range name {
		if unicode.IsDigit(char) || unicode.IsSymbol(char) {
			return errors.New("name must contain only characters"), false
		}
	}
	return nil, true
}

// IsValidPhoneNumber checks if a given string represents a valid phone number in the specified region.
func IsValidPhoneNumber(phoneNumber, defaultRegion string) bool {
	phoneNumber = strings.TrimSpace(phoneNumber)

	// Parse the phone number
	num, err := libphonenumber.Parse(phoneNumber, defaultRegion)
	if err != nil {
		return false
	}

	// Check if the parsed number is valid
	return libphonenumber.IsValidNumber(num)
}

// IsValidTitle checks if a given string consists of only alphanumeric characters, underscores (_), and hyphens (-).
func IsValidTitle(input string) bool {
	// Define a regular expression pattern to match alphanumeric characters, _, and -
	pattern := consts.ValidTitlePattern

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Use the MatchString function to check if the input matches the pattern
	return regex.MatchString(input)
}

// IsEmptyValue checks if the provided value is empty.
func IsEmptyValue[T comparable](value T) bool {
	switch v := any(value).(type) {
	case string:
		return v == ""
	case int, int16, int32, int64:
		return v == 0
	case float32, float64:
		return v == 0.0
	case bool:
		return !v
	default:
		// Use a type assertion to check if the value is a pointer
		if ptr, ok := any(value).(interface{ IsNil() bool }); ok {
			// If it's a pointer, return true if it's nil
			return ptr == nil || ptr.IsNil()
		}
		// If it's not a pointer, return true
		return true
	}
}

// ValidateStringLength checks if the string length is between minLength && maxLength
func ValidateStringLength(value string, minLength, maxLength int) error {
	n := len(strings.TrimSpace(value))
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

// ValidateISRC checks if an ISRC code is valid.
func ValidateISRC(isrc string) bool {
	// The pattern ensures that the ISRC code matches the format: CCRRRYYNNNNN
	// where C is a 2-letter country code, R is a 3-character registrant code,
	// Y is a 2-digit year code, and N is a 5-digit designation code.
	pattern := regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{3}\d{7}$`)

	// If the ISRC matches the pattern, it is considered valid.
	return pattern.MatchString(isrc)
}

// ValidateISWC checks if an ISWC code is valid.
func ValidateISWC(iswc string) bool {

	// The pattern ensures that the ISWC code matches the format: T1234567894
	pattern := regexp.MustCompile(`^[T]{1}[0-9]{9}\d{1}$`)

	// If the ISWC matches the pattern, it is considered valid.
	return pattern.MatchString(iswc)
}

// IsValidOrder checks if an order is valid.
func IsValidOrder(order string) bool {
	orderPattern := `^(asc|desc|ASC|DESC)$`
	return regexp.MustCompile(orderPattern).MatchString(order)
}

// CountPlaceholders for counting the placeholders
func CountPlaceholders(input string) int {
	count := 0
	openBrace := false

	for _, char := range input {
		if char == '{' {
			openBrace = true
		} else if char == '}' && openBrace {
			count++
			openBrace = false
		}
	}
	return count
}

// ExtractPlaceholders for extracting placeholders
func ExtractPlaceholders(template string) []string {
	re := regexp.MustCompile(`\{([^{}]+)\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	var placeholders []string
	for _, match := range matches {
		placeholders = append(placeholders, match[1])
	}

	return placeholders
}

// GenerateRequiredPlaceholders for setting the value required to the placeholders
func GenerateRequiredPlaceholders(template string) map[string]interface{} {
	placeholders := ExtractPlaceholders(template)

	placeholderRequirements := make(map[string]interface{})

	for _, placeholder := range placeholders {
		placeholderRequirements[placeholder] = "required"
	}

	return placeholderRequirements
}

// GetKey is used for getting the key of a map
func GetKey(casted []string) []string {
	var output []string
	for _, value := range casted {
		value = strings.Split(value, "=")[0]
		val, ok := consts.LengthConstraints[value]
		if ok {
			output = append(output, val)
			continue
		}
		output = append(output, value)
	}
	return output

}
