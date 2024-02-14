## Overview
This provides functions for handling date and time conversions, as well as adding hours to a given time.

## Index
- [ConvertDate(date string, format ...string) (string, time.Time, error)](#func-ConvertDate)
- [AddHours(d time.Time, hours int, format ...string) string](#func-AddHours)
- [FormatDateTime(inputTime string, inputFormat string, outputFormat string)](#func-FormatDateTime)

### func ConvertDate

    ConvertDate(date string, format ...string) (string, time.Time, error)

This function is used to parse a `date` string and convert it to `UTC` time. It takes the date string and an optional format string as parameters and returns a formatted date string in the provided format, the parsed `time.Time` object, and an `error`, if any. If the format is not provided, the default date format (dd-mm-yyyy) is used. If the date parameter is empty, an error is returned with the message "date is not provided."

### func AddHours

    AddHours(d time.Time, hours int, format ...string) string

This function is used to add a specified number of `hours` to a given `time.Time` object and return the resulting time as a formatted string. If the format is not provided, the default date format (dd-mm-yyyy) is used.

### func FormatDateTime

    FormatDateTime(inputTime string, inputFormat string, outputFormat string)

The FormatDateTime function is a Go function that converts a given date and time from one string format to another. It takes three inputs: the `inputTime` string representing the date and time, the `inputFormat` specifying the format of inputTime, and the `outputFormat` for the desired output format. It first attempts to parse the inputTime using the inputFormat. If successful, it returns the formatted date and time string in the specified outputFormat. If there's a parsing error, it returns an error to indicate that the input format does not match the provided date and time string.