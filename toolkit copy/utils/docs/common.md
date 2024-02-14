## Overview

This module deals with common functions necessary for the proper working of localization module. It deals with GetEndPoints , ExtractRoutePortion, FieldMapping and AppendValuesToMap

## Index

- [ GetEndPoints(contextEndpoints models.ResponseData, url string, method string) string](#func-GetEndPoints)
- [ ExtractRoutePortion(route string) (string, error)](#func-ExtractRoutePortion)
- [ FieldMapping(fieldsMap map[string][]string) string](#func-FieldMapping)
- [ AppendValuesToMap(m map[string][]string, key string, values ...string)](#func-AppendValuesToMap)
- [ func ConvertSliceToUUIDs(arg []interface{}) ([]string, error)](#func-ConvertSliceToUUIDs)



### func GetEndPoints

    GetEndPoints(contextEndpoints models.ResponseData, url string, method string) string

The function retrieves an endpoint based on the provided URL and HTTP method from the given contextEndpoints. It first extracts a portion of the provided URL using the ExtractRoutePortion function. Then, it iterates through the Data array within contextEndpoints and compares each data item's URL and method with the extracted URL and provided HTTP method. If a match is found, the corresponding endpoint is assigned to the endpoint variable. The function will returns the retrieved endpoint as a string. If no matching endpoint is found, an empty string is returned.

### func ExtractRoutePortion

    ExtractRoutePortion(route string) (string, error)

The function extracts a portion of the provided route based on a specific delimiter, which is "/:version". It splits the input route string using "/:version" as the delimiter. If the route contains the "/:version" delimiter and the split results in more than one part, it returns the portion following the delimiter. If the route does not contain the "/:version" delimiter, it returns an empty string and an error indicating that the route does not contain the expected delimiter.

### func FieldMapping

    FieldMapping(fieldsMap map[string][]string) string

The function iterates over the provided fieldsMap, where each entry consists of a field name (as the key) and a list of associated error keys (as the value). For each field in the map, it creates a formatted string that combines the field name and the associated error keys in the format "field_name:error_key1|error_key2|...". These formatted strings for each field are then appended to a slice called fields. Finally, the function joins all the formatted strings with commas and returns the resulting concatenated string. The output string represents the mapping of field names to their respective error keys in a specified format.

### func AppendValuesToMap

    AppendValuesToMap(m map[string][]string, key string, values ...string)

The function appends the provided values to the slice associated with the given key in the map. If the key already exists in the map, the values are appended to the existing slice of values for that key. If the key does not exist in the map, a new entry is created with the given key and the provided values as the associated slice of strings. The function does not have a return value, as it modifies the map passed as a parameter directly.


### func ConvertSliceToUUIDs

    func ConvertSliceToUUIDs(arg []interface{}) ([]string, error)

The `ConvertSliceToUUIDs` function takes a slice of interface{} as input and aims to convert it into a slice of strings representing UUIDs.  If any of the element cannot be converted to a UUID, it returns an error. If the conversion is successful, the UUID is transformed into its string representation and added to the result slice.