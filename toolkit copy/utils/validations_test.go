// nolint
package utils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsEmpty(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"Hello", false},
		{"123", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			result := IsEmpty(testCase.input)
			require.Equal(t, testCase.expected, result)
		})
	}
}

func TestIsValidIPAddress(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"192.168.1.1", true},
		{"::1", true},
		{"256.256.256.256", false},
		{"invalid", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			result := IsValidIPAddress(testCase.input)
			require.Equal(t, testCase.expected, result)
		})
	}
}
func TestIsValidEmail(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"", false},
		{"123", false},
		{"12gmail", false},
		{"123@fjhjfdfj...cdf", false},
		{"noreplyemail@gmail.com", true},
		{"logo.png", false},
		{"NOREPLYEMAIL@gmail.com", true},
		{"person@partner", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			result := IsValidEmail(testCase.input)
			require.Equal(t, testCase.expected, result)
		})
	}
}

func TestIsValidURL(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"", false},
		{"123", false},
		{"123@fjhjfdfj...cdf", false},
		{"https://www.mywebsite.com/profile", true},
		{"mywebsite", false},
		{".png", false},
		{"http://127.0.0.1/", true},
		{"http://localhost:3000/", true},
		{"http://abc-de-f-.example.com", false},
		{"https://127.0.0.1/a/b/c", true},
		{"http://www.-foobar.com/", false},
		{"http://www.demo---test.com/", false},
		{"http://.foo.com", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			result := IsValidURL(testCase.input)
			require.Equal(t, testCase.expected, result)
		})
	}
}
func TestIsValidName(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"demo user", true},
		{"123", false},
		{"name", true},
		{"AthiraAaathi", true},
		{"ng", false},
		{"n", false},
		{"!@#$$%%", false},
		{"nameffhgdsgfhgdhfgdgfdguhdusgh", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			_, result := IsValidName(testCase.input)
			require.Equal(t, testCase.expected, result)
		})
	}
}

func TestIsValidPhoneNumber(t *testing.T) {
	tests := []struct {
		phoneNumber   string
		defaultRegion string
		expected      bool
	}{
		{"+442079460958", "GB", true},    // Valid UK phone number
		{"123456789", "US", false},       // Invalid phone number without country code
		{"123", "US", false},             // Invalid partial phone number
		{"+1 123", "US", false},          // Invalid partial phone number
		{"", "US", false},                // Empty phone number
		{"+1 8005551212", "US", true},    // Toll-free US phone number
		{"+44 20 7946 0958", "GB", true}, // Valid UK phone number with spaces
		{"+1 123 456 7890", "US", false}, // Invalid US phone number with spaces
		{"123-456-7890", "US", false},    // Invalid phone number without country code
	}

	for _, test := range tests {
		result := IsValidPhoneNumber(test.phoneNumber, test.defaultRegion)
		if result != test.expected {
			t.Errorf("For phone number %s in region %s, expected %v, but got %v", test.phoneNumber, test.defaultRegion, test.expected, result)
		}
	}
}

func TestIsValidTitle(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		// Positive Test Cases
		{"ValidName123", true},     // Alphanumeric characters
		{"_Valid_Name_", true},     // Underscores allowed
		{"Hyphen-Name", true},      // Hyphens allowed
		{"Alphanumeric_123", true}, // Mixed alphanumeric and underscores
		{"ValidName123-", true},    // Hyphen at the end

		// Negative Test Cases
		{"Not Valid", false},       // Space not allowed
		{"Invalid@", false},        // Special character not allowed
		{"Name with space", false}, // Space not allowed
		{"@12345", false},          // Numeric characters only
		{"", false},                // Empty string
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := IsValidTitle(tc.input)
			if result != tc.expected {
				t.Errorf("Expected IsValidArtistName(%s) to return %v, but got %v", tc.input, tc.expected, result)
			}
		})
	}
}

func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{input: "", expected: true},
		{input: "hello", expected: false},
		{input: 0, expected: true},
		{input: 1, expected: false},
		{input: 0.0, expected: true},
		{input: 0.1, expected: false},
		{input: true, expected: false},
		{input: false, expected: true},
		{input: (*int)(nil), expected: true},
	}

	for _, test := range tests {
		actual := IsEmptyValue(test.input)
		require.Equal(t, test.expected, actual, "input: %v", test.input)
	}
}

func TestValidateStringLength(t *testing.T) {
	const (
		validMinLength = 3
		validMaxLength = 6
	)

	testCases := []struct {
		value     string
		minLength int
		maxLength int
		expectErr error
	}{
		{"abcdef", validMinLength, validMaxLength, nil}, // Valid length
		{"ab", validMinLength, validMaxLength, fmt.Errorf("must contain from %d-%d characters", validMinLength, validMaxLength)},      // Too short
		{"abcdefg", validMinLength, validMaxLength, fmt.Errorf("must contain from %d-%d characters", validMinLength, validMaxLength)}, // Too long
	}

	for _, tc := range testCases {
		t.Run(tc.value, func(t *testing.T) {
			err := ValidateStringLength(tc.value, tc.minLength, tc.maxLength)
			require.Equal(t, tc.expectErr, err)
		})
	}
}

func TestValidateISRC(t *testing.T) {
	testCases := []struct {
		isrc   string
		expect bool
	}{
		{"USAAA1234567", true},   // Valid ISRC
		{"usaaa1234567", false},  // Lowercase country code
		{"USAAA12345678", false}, // Too long
		{"USAAA123456", false},   // Too short
		{"USAAA12345G", false},   // Non-numeric last digit
	}

	for _, tc := range testCases {
		t.Run(tc.isrc, func(t *testing.T) {
			result := ValidateISRC(tc.isrc)
			if result != tc.expect {
				t.Errorf("Expected %v, but got %v", tc.expect, result)
			}
		})
	}
}

func TestValidateISWC(t *testing.T) {
	testCases := []struct {
		iswc   string
		expect bool
	}{
		{"T1234567894", true},   // Valid ISWC
		{"T12345678945", false}, // Too long
		{"T123456789", false},   // Too short
		{"T123456789G", false},  // Non-numeric last digit
	}

	for _, tc := range testCases {
		t.Run(tc.iswc, func(t *testing.T) {
			result := ValidateISWC(tc.iswc)
			if result != tc.expect {
				t.Errorf("Expected %v, but got %v", tc.expect, result)
			}
		})
	}
}

func TestGetKey(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"key1=value1", "key2=value2"}, []string{"key1", "key2"}},
		{[]string{"key3=value3", "key4=value4"}, []string{"key3", "key4"}},
		{[]string{"key5=value5", "key6=value6"}, []string{"key5", "key6"}},
		{[]string{"key7=value7", "key8=value8"}, []string{"key7", "key8"}},
		{[]string{"key9=value9", "key10=value10"}, []string{"key9", "key10"}},
	}

	for _, test := range tests {
		result := GetKey(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For input %v, expected %v but got %v", test.input, test.expected, result)
		}
	}
}

func TestGenerateRequiredPlaceholders(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		{"Hello, {name}! How are you, {name}?", map[string]interface{}{"name": "required"}},
		{"{city} is a beautiful place.", map[string]interface{}{"city": "required"}},
		{"There are {count} items in stock.", map[string]interface{}{"count": "required"}},
		{"", map[string]interface{}{}}, // Empty input should return an empty map
	}

	for _, test := range tests {
		result := GenerateRequiredPlaceholders(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For input %q, expected %v but got %v", test.input, test.expected, result)
		}
	}
}

func TestCountPlaceholders(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"Hello, {name}! How are you, {name}?", 2},
		{"{city} is a beautiful place.", 1},
		{"There are {count} items in stock.", 1},
		{"{one} {two} {three}", 3},
		{"No placeholders here", 0}, // No placeholders, should return 0
		{"{}", 1},                   // One placeholder, should return 1
		{"{{}}", 1},                 // Escaped placeholder, should return 1
		{"{{{}}}", 1},               // Nested escaped placeholder, should return 1
		{"{{", 0},                   // Unmatched opening brace, should return 0
		{"}}", 0},                   // Unmatched closing brace, should return 0
	}

	for _, test := range tests {
		result := CountPlaceholders(test.input)
		if result != test.expected {
			t.Errorf("For input %q, expected %d but got %d", test.input, test.expected, result)
		}
	}
}
func TestIsValidOrder(t *testing.T) {
	tests := []struct {
		order    string
		expected bool
	}{
		{"asc", true},
		{"desc", true},
		{"ASC", true},
		{"DESC", true},
		{"Asc", false},  // Incorrect casing
		{"Desc", false}, // Incorrect casing
		{"invalid", false},
		{"", false}, // Empty string should return false
	}

	for _, test := range tests {
		result := IsValidOrder(test.order)
		if result != test.expected {
			t.Errorf("For order %q, expected %v but got %v", test.order, test.expected, result)
		}
	}
}

func TestExtractPlaceholders(t *testing.T) {
	tests := []struct {
		name       string
		template   string
		wantResult []string
	}{
		{
			name:       "NoPlaceholders",
			template:   "This is a test string with no placeholders",
			wantResult: nil,
		},
		{
			name:       "SinglePlaceholder",
			template:   "Hello, {name}!",
			wantResult: []string{"name"},
		},
		{
			name:       "MultiplePlaceholders",
			template:   "Welcome, {firstName} {lastName} to our {event}!",
			wantResult: []string{"firstName", "lastName", "event"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := ExtractPlaceholders(tt.template)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ExtractPlaceholders(%q) = %v, want %v", tt.template, gotResult, tt.wantResult)
			}
		})
	}
}
