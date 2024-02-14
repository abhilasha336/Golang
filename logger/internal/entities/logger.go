package entities

import (
	"time"

	"gitlab.com/tuneverse/toolkit/utils"
)

// Log represents a log entry in the application.
type Log struct {
	RequestID    string                `json:"req_id,omitempty" bson:"req_id,omitempty"`
	Method       string                `json:"method,omitempty" bson:"method,omitempty"`
	Endpoint     string                `json:"endpoint,omitempty" bson:"endpoint,omitempty"`
	Service      string                `json:"service,omitempty" bson:"service,omitempty"`
	UserIP       string                `json:"user_ip,omitempty" bson:"user_ip,omitempty"`
	Filename     string                `json:"filename,omitempty" bson:"filename,omitempty"`
	Message      string                `json:"message,omitempty" bson:"message,omitempty"`
	ResponseCode int                   `json:"response_code,omitempty" bson:"response_code,omitempty"`
	DurationMS   string                `json:"duration_ms,omitempty" bson:"duration_ms,omitempty"`
	Day          int                   `json:"day,omitempty" bson:"day,omitempty"`
	Month        int                   `json:"month,omitempty" bson:"month,omitempty"`
	Year         int                   `json:"year,omitempty" bson:"year,omitempty"`
	Timestamp    time.Time             `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	LogLevel     string                `json:"log_level,omitempty" bson:"log_level,omitempty"`
	URL          string                `json:"uri,omitempty" bson:"uri,omitempty"`
	RequestDump  utils.RequestDataDump `json:"request_dump,omitempty" bson:"request_dump,omitempty"`
	ResponseDump utils.ResponseData    `json:"response_dump,omitempty" bson:"response_dump,omitempty"`
	Data         Data                  `json:"data,omitempty" bson:"data,omitempty"`
}

// Data represents additional data associated with a log entry.
type Data struct {
	UserID    string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	PartnerID string `json:"partner_id,omitempty" bson:"partner_id,omitempty"`
}

// LogParams represents the parameters for filtering and querying logs.
type LogParams struct {
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	Service    string `form:"service"`
	HTTPMethod string `form:"method"`
	LogLevel   string `form:"log_level"`
	Endpoint   string `form:"endpoint"`
	UserIP     string `form:"user_ip"`
	UserID     string `form:"user_id"`
	PartnerID  string `form:"partner_id"`
	Page       int32  `form:"page"`
	Limit      int32  `form:"limit"`
}
