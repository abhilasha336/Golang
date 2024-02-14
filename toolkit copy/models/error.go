package models

// DataItem
type DataItem struct {
	URL      string `json:"URL"`
	Method   string `json:"Method"`
	Endpoint string `json:"Endpoint"`
}

// ResponseData
type ResponseData struct {
	Data []DataItem `json:"data"`
}
type ErrorData struct {
	Field     string   `json:"field"`
	Message   []string `json:"message"`
	ErrorCode string   `json:"error_code"`
	Help      string   `json:"help"`
}

type ErrorResponse struct {
	Code    string
	Message []string
}
type ErrorDetails struct {
	Code    float64
	Message string
}
