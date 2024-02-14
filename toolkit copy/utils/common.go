package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/tuneverse/toolkit/models"
)

var (
	BitSize64 = 64
)

func GetEndPoints(contextEndpoints models.ResponseData, url string, method string) string {

	// Access the data
	var endpoint string
	extractURL, _ := ExtractRoutePortion(url)
	for _, dataItem := range contextEndpoints.Data {
		if dataItem.URL == extractURL && dataItem.Method == method {
			endpoint = dataItem.Endpoint
			break
		}
	}
	return endpoint
}

func ExtractRoutePortion(route string) (string, error) {
	// Split the route using "/:version"
	routeParts := strings.Split(route, "/:version")
	if len(routeParts) > 1 {
		return routeParts[1], nil
	}
	return "", errors.New("route does not contain '/:version'")
}

// Function to generate fields of the form firstname:required|valid

func FieldMapping(fieldsMap map[string]models.ErrorResponse) string {
	var fields []string

	for fieldName, response := range fieldsMap {
		// Convert Message slice to a string
		errorKeyString := strings.Join(response.Message, "|")
		field := fmt.Sprintf("%s:%s", fieldName, errorKeyString)
		fields = append(fields, field)
	}

	return strings.Join(fields, ",")
}

// FormatToDecimalLimit
func FormatWithDecimalLimit(val float64, d string) (string, float64, error) {
	if d == "" {
		d = "2"
	}
	decimalLimit := "%." + d + "f"
	format := fmt.Sprintf(decimalLimit, val)

	val, err := strconv.ParseFloat(format, BitSize64)
	return format, val, err
}

// ConvertSliceToUUIDs converts a slice of interface{} to a slice of uuid.UUID.
func ConvertSliceToUUIDs(arg []interface{}) ([]string, error) {
	var result []string
	for _, item := range arg {
		id, ok := item.(string)
		if !ok {
			return nil, errors.New("inavlid uuid")
		}
		uuidValue, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		result = append(result, uuidValue.String())

	}
	return result, nil
}
