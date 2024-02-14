## Overview

This package contains utility functions for working with database queries.

## Index

- [BuildSearchQuery(value string, fields ...string) string](#func-BuildSearchQuery)
- [CalculateOffset(page, limit int32) string](#func-CalculateOffset)
-[PreparePlaceholders(n int)string](#func-PreparePlaceholders)

### func BuildSearchQuery

    BuildSearchQuery(value string, fields ...string) string

BuildSearchQuery constructs an SQL query that performs a `case-insensitive` partial search for the given value in the specified fields. It combines multiple field searches using the `ILIKE` operator and `OR` conditions. The resulting query is suitable for use in database searches.

### func CalculateOffset

    CalculateOffset(page, limit int32) string

CalculateOffset computes the `OFFSET` and `LIMIT` SQL clause based on the `current page` and the `number of items per page`. It is useful for paginating database queries by specifying the range of results to fetch.

### func-PreparePlaceholders

    PreparePlaceholders(n int)string
    
PreparePlaceholders generates and returns a comma-separated string of SQL placeholders for 'n' parameters.It is useful for dynamically creating placeholders when building SQL queries with a variable number of parameters.
