package models

// MetaData represents metadata information for paginating a collection of data.
type MetaData struct {
	Total       int64 `json:"total"`
	PerPage     int32 `json:"per_page"`
	CurrentPage int32 `json:"current_page"`
	Next        int32 `json:"next"`
	Prev        int32 `json:"prev"`
}
