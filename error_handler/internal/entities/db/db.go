package db

// ErrorFieldAttributes represents the attributes of an error field.
type ErrorFieldAttributes map[string]any

// ErrorData represents error data structure.
type ErrorData struct {
	Field ErrorFieldAttributes `json:"field"`
}
