package entities

import "gitlab.com/tuneverse/toolkit/models"

// Response represents a standard response format for API responses.
type Response struct {
	StatusCode int              `json:"status_code"`
	Message    string           `json:"message"`
	MetaData   *models.MetaData `json:"meta_data,omitempty"`
	Data       interface{}      `json:"data,omitempty"`
}
