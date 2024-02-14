package utilities

import (
	"encoding/json"
	"errors"
	"member/internal/consts"
	"member/internal/entities"
	"regexp"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/ttacon/libphonenumber"
)

// ValidateEmail checks if the provided email address is valid.
func ValidateEmail(email string) bool {
	if strings.TrimSpace(email) == "" {
		return false // Email is empty
	}

	err := checkmail.ValidateFormat(email)
	return err == nil
}

// ValidateName checks if the provided name is valid (letters, numbers, spaces, periods, apostrophes, hyphens/dashes).
func ValidateName(name string) bool {
	regex := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\s.'-]*$`)
	return regex.MatchString(name)
}

// IsValidUUID checks if the input string is a valid UUID.
func IsValidUUID(id string) bool {
	if len(id) != 36 {
		return false
	}
	segmentLengths := []int{8, 4, 4, 4, 12}
	index := 0
	for _, length := range segmentLengths {
		if index > 0 && id[index-1] != '-' {
			return false
		}
		for _, c := range id[index : index+length] {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
		index += length + 1
	}
	return true
}

// ValidatePassword checks the validity of a password.
func ValidatePassword(password string) error {
	// Check minimum length
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	hasUppercase := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUppercase = true
			break
		}
	}
	if !hasUppercase {
		return errors.New("Password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	hasLowercase := false
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLowercase = true
			break
		}
	}
	if !hasLowercase {
		return errors.New("Password must contain at least one lowercase letter")
	}

	// Check for at least one special character
	specialCharPattern := regexp.MustCompile(`[!@#$%^&*()_+{}\[\]:;<>,.?~]`)
	if !specialCharPattern.MatchString(password) {
		return errors.New("Password must contain at least one special character")
	}

	// Check for spaces
	if regexp.MustCompile(`\s`).MatchString(password) {
		return errors.New("Password cannot contain spaces")
	}

	return nil
}

// IsValidPhoneNumber checks if a given string represents a valid phone number in the specified region.
//
// This function validates whether the provided string resembles a valid phone number in the specified region.
// It first trims any leading or trailing white spaces from the input string and then parses it as a phone number
// using the default region provided. If the parsing is successful, it checks if the parsed number is valid.
//
// Parameters:
//
//	@phoneNumber (string): The string to be validated as a phone number.
//	@defaultRegion (string): The default region code (e.g., "US", "GB") to use for parsing the phone number.
//
// Returns:
//
//	bool: True if the input string resembles a valid phone number in the specified region, false otherwise.
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

// IsValidZIPCode checks if the provided string is a valid ZIP code.
// It allows ZIP codes with lengths between 4-10 characters, including alphanumeric characters and hyphens.
func IsValidZIPCode(zipcode string) bool {
	// Define a regular expression pattern for ZIP codes.
	// This pattern allows for 4 to 10 characters, including alphanumeric characters and hyphens.
	pattern := `^[a-zA-Z0-9\-]{4,10}$`

	// Compile the regular expression pattern.
	regex := regexp.MustCompile(pattern)

	// Use the regex.MatchString() function to check if the ZIP code matches the pattern.
	return regex.MatchString(zipcode)
}

// ValidateNameLength is for validating length of the name
func ValidateMaximumNameLength(name string) bool {
	if len(name) > consts.MaximumNameLength {
		return false
	}
	return true
}

func ValidateJwtToken(token string, jwtKey string) (response entities.JwtValidateResponse) {

	jwtKeyVal := []byte(jwtKey)

	defer func() {
		if rec := recover(); rec != nil {
			response.ErrorMsg = "Token seems to be invalid!"
		}
	}()

	response = entities.JwtValidateResponse{
		Valid:       false,
		MemberID:    "0",
		PartnerID:   "",
		PartnerName: "",
		MemberName:  "",
		Roles:       []string{},
		MemberEmail: "",
		ErrorMsg:    "",
	}

	claims := &entities.Claims{}
	if len(token) == 0 {
		response.ErrorMsg = "No authorization token passed"
	} else {
		tokenParsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKeyVal, nil
		})
		if claims, ok := tokenParsed.Claims.(*entities.Claims); ok && tokenParsed.Valid {
			response.Valid = true
			response.MemberID = claims.MemberID
			response.PartnerID = claims.PartnerID
			response.MemberType = claims.MemberType
			response.Roles = claims.Roles
			response.MemberEmail = claims.MemberEmail
			response.MemberName = claims.MemberName
			response.PartnerName = claims.PartnerName
		} else {
			response.ErrorMsg = err.Error()
		}
	}
	return response
}

// IsEmpty checks whether the given string is empty
func IsEmpty(str string) bool {

	return strings.TrimSpace(str) == ""
}

func MarshalNullableString(ns entities.NullableString) ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return []byte(`""`), nil
}
