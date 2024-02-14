package utils

import (
	"fmt"
	"strings"
)

// BuildSearchQuery generates an SQL query for a case-insensitive partial search on one or more fields.
func BuildSearchQuery(value string, fields ...string) string {
	var searchQ strings.Builder
	for index, field := range fields {
		_, err := searchQ.WriteString(fmt.Sprintf(" %s %s '%%%s%%' ", field, "ILIKE", value))
		_ = err
		if index != len(fields)-1 {
			_, err := searchQ.WriteString("OR")
			_ = err
		}
	}
	return searchQ.String()
}

// CalculateOffset calculates the offset value from page and limit
func CalculateOffset(page, limit int32) string {
	offset := fmt.Sprintf("%s %d %s %d", "LIMIT", limit, "OFFSET", (page-1)*limit)
	return offset
}

// PreparePlaceholders generates and returns a comma-separated string of SQL placeholders for 'n' parameters.
func PreparePlaceholders(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(fmt.Sprintf("$%d,", i+1))
	}
	return strings.TrimSuffix(b.String(), ",")
}
