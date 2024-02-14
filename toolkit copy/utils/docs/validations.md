## Overview

This provides function for common validations.

## Index

- [IsEmpty(str string) bool](#func-IsEmpty)
- [IsValidIPAddress(ipAddress string) bool](#func-IsValidIPAddress)
- [IsValidEmail(email string) bool](#func-IsValidEmail)
- [IsValidURL(url string) bool](#func-IsValidURL)
- [IsValidName(name string)( bool,error)](#func-IsValidName)
- [IsValidPhoneNumber(phoneNumber, defaultRegion string) bool](#func-IsValidPhoneNumber)
- [IsValidTitle(input string) bool](#func-IsValidTitle)
- [IsEmptyValue[T comparable](value T) bool](#func-IsEmptyValue)
- [ValidateStringLength(value string, minLength, maxLength int) error](#func-ValidateStringLength)
- [ValidateISRC(isrc string) bool](#func-ValidateISRC)
- [IsValidOrder(order string) bool](#func-IsValidOrder)
- [CountPlaceholders(input string) int](#func-CountPlaceholders)
- [ExtractPlaceholders(template string) []string](#func-ExtractPlaceholders)
- [GenerateRequiredPlaceholders(template string) map[string]interface{}](#func-GenerateRequiredPlaceholders)
- [GetKey(casted []string) []string](#func-GetKey)


### func IsEmpty

    IsEmpty(str string) bool

This function is used to check whether a given string is empty or not. It takes a single parameter, `str`, which is the string to be checked. The function returns `true` if the string is empty (consists only of whitespace characters) and `false` otherwise.

### func IsValidIPAddress

    IsValidIPAddress(ipAddress string) bool

This function is used to check if the given string is a valid IP address, either IPv4 or IPv6. It takes a single parameter, `ipAddress`, which is the string to be checked. The function attempts to parse the string as an IP address using the `net.ParseIP` function. If parsing succeeds, the function returns `true`, indicating that the string is a valid IP address. If parsing fails, it returns `false`, indicating that the string is not a valid IP address.

### func IsValidEmail

    IsValidEmail(email string) bool

This function is used to check if the given email is a valid email or not.It takes a single parameter ,
`email` . The function attempts to parse the string using `mail.ParseAddress` .If it succeeds ,the function returns `true` , which indicates the string is valid email .If it fails , it returns `false`,
indicating the string is not a valid email.

### func IsValidURL

    IsValidURL(url string) bool

This function is used to check if the given url is a valid or not.It takes a single parameter ,
`url` . The function attempts to parse the string using `url.ParseRequestURI` .If it succeeds ,the function returns `true` , which indicates the string is valid url .If it fails , it returns `false`,
indicating the string is not a valid url.

### func IsValidName

    IsValidName(name string)( bool,error)

This function is used to check if the given name is a valid or not. It takes a single parameter ,`name` . The function attempts to checks the length of the string is between minimum and maximum characters as we have defined in constants such as `MinLength` and `MaxLength` and it also checks whether it contains any special characters or symbols .If it succeeds ,the function returns `true` and error as `nil` , which indicates the string is valid name .If it fails , it returns `false`and particular error based on it ,which indicating the string is not a valid name.

### func IsValidPhoneNumber

    IsValidPhoneNumber(phoneNumber, defaultRegion string) bool

The IsValidPhoneNumber function validates a given `phoneNumber` by first trimming any whitespace from the input. It then attempts to parse the phone number using the libphonenumber library, providing a default region code. If the parsing is successful, indicating that the phone number is in a valid format for the specified region, it checks if the parsed number is valid using the libphonenumber library. If the parsed number is valid, the function returns `true`; otherwise, it returns `false`.

### func IsValidTitle

    IsValidTitle(input string) bool

The IsValidTitle function checks the validity of a given string input against a specific pattern. It uses a regular expression pattern `^[a-zA-Z0-9_-]+$` to match alphanumeric characters, underscores (\_), and hyphens (-). The function compiles this pattern into a regular expression and then uses the `MatchString` function from the `regexp` package to check if the input string matches the defined pattern. If the input string consists only of alphanumeric characters, underscores, and hyphens, the function returns `true`; otherwise, it returns `false`. This function is commonly used to validate titles in various contexts, ensuring they contain only allowed characters.

### func IsEmptyValue

    IsEmptyValue[T comparable](value T) bool

The function is a generic function that checks if a given `value` is empty or not. It uses a type switch to determine the type of the input value and then checks if the value is empty based on its type. The function can handle strings, integers, floats, booleans, and pointers. It returns `true` if the value is empty and `false` otherwise.


### func ValidateStringLength

func ValidateStringLength(value string, minLength, maxLength int) error

The `ValidateStringLength` function checks if the length of a given string falls within a specified range. It is useful for validating the length of strings. If the length is outside the specified range, it returns an error.

### func ValidateISRC

    func ValidateISRC(isrc string) bool

The `ValidateISRC` function checks if an ISRC (International Standard Recording Code) is valid. It verifies if the provided string adheres to the ISRC format.


### func ValidateISWC

    func ValidateISWC(iswc string) bool

The `ValidateISWC` function checks if an ISWC (International Standard Work Code) is valid. It validates if the provided string. conforms to the ISWC format.

### func IsValidOrder

This is used for checking whether the given order is correct ot not.

### func CountPlaceholders

This is used for counting the placeholders.

### func ExtractPlaceholders

This is used for extracting placeholders from a line of text.

### func GenerateRequiredPlaceholders

This is used for generte a function for setting the value required to the placeholders.

### func GetKey

This is used for obtaining key from the map.