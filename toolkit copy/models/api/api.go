package api

// Response represents the structure of the API response from a service.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}
