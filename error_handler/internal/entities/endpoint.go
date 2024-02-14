package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

// Endpoint represents an API endpoint.
type Endpoint struct {
	Url      string `json:"url" bson:"url"`
	Method   string `json:"method" bson:"method"`
	Endpoint string `json:"endpoint" bson:"endpoint"`
}

// GetEndpoint represents an API endpoint with an ID.
type GetEndpoint struct {
	ID       primitive.ObjectID
	Url      string
	Method   string
	Endpoint string
}

// DataItem represents an item of data.
type DataItem struct {
	URL      string `json:"URL"`
	Method   string `json:"Method"`
	Endpoint string `json:"Endpoint"`
}

// ResponseData represents a response containing data items.
type ResponseData struct {
	Data []DataItem `json:"data"`
}

// EndpointData represents data about an endpoint.
type EndpointData struct {
	URL      string `json:"url"`
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
}

// RequestData represents a request for endpoint data
type RequestData struct {
	Field EndpointData `json:"field"`
}

// MyDocument represents a document with endpoint data.
type MyDocument struct {
	ID   primitive.ObjectID `bson:"_id"`
	Data []EndpointData     `bson:"data"`
}

type ErrorResponse struct {
	Message   string         `json:"message"`
	ErrorCode any            `json:"errorCode"`
	Errors    map[string]any `json:"errors"`
}
