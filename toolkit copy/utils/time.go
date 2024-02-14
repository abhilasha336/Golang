package utils

import (
	"fmt"
	"time"

	"gitlab.com/tuneverse/toolkit/consts"
)

// ConvertDate parses the date to UTC.
//
// params:
//
//	@date: required.
//	@format: optional.
//
// if format is not provided the date should be of this format dd-mm-yyyy
func ConvertDate(date string, format ...string) (string, time.Time, error) {
	if IsEmpty(date) {
		return "", time.Time{}, fmt.Errorf("date is not provided")
	}
	if len(format) == 0 {
		format = append(format, consts.DefaultDateFormat)
	}
	parsedDate, err := time.Parse(format[0], date)
	if err == nil {
		return parsedDate.Format(format[0]), parsedDate, nil
	}
	return "", time.Time{}, err
}

// AddHours adds a specified number of hours to a given time and returns
// the resulting time as a string in the RFC3339 format.
//
// params:
//
//	@time: required.
//	@hours: required.
//	@format: optional.
func AddHours(d time.Time, hours int, format ...string) string {
	if len(format) == 0 {
		format = append(format, consts.DefaultDateFormat)
	}
	return d.Add(time.Duration(hours) * time.Hour).Format(format[0])
}

func DateTime(t time.Time) string {
	return t.UTC().Format(time.DateTime)
}

// FormatDateTime will returns the time in required format
// params
// @ inputTime - Input time
// @ outputFormat - Output format
// Returns the formatted time
func FormatDateTime(inputTime string, outputFormat ...string) (string, error) {
	var format string
	if len(outputFormat) > 0 {
		format = outputFormat[0]
	} else {
		format = consts.OutputFormat
	}

	parsedTime, err := time.Parse(consts.InputFormat, inputTime)
	if err != nil {
		return "", err
	}

	formattedTime := parsedTime.Format(format)

	return formattedTime, nil
}
