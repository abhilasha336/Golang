package utils

import "github.com/gin-gonic/gin"

type ResponseData struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
}

// ResponseWriterWrapper is a custom response writer to capture the response
type ResponseWriterWrapper struct {
	gin.ResponseWriter
	responseData *ResponseData
}

// NewResponseWriterWrapper creates a new ResponseWriterWrapper
func NewResponseWriterWrapper(w gin.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		ResponseWriter: w,
		responseData:   &ResponseData{},
	}
}

// WriteHeader is called to write the response header
func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	rw.responseData.StatusCode = code
}

// Write is called to write the response body
func (rw *ResponseWriterWrapper) Write(b []byte) (int, error) {
	// Capture the response body
	rw.responseData.Body = string(b)

	// Capture the response headers
	rw.responseData.Headers = rw.ResponseWriter.Header()

	return rw.ResponseWriter.Write(b)
}

// GetResponseData returns the captured response data
func (rw *ResponseWriterWrapper) GetResponseData() ResponseData {
	return *rw.responseData
}
