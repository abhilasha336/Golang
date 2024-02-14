package models

// Struct for storing activity log details
type ActivityLog struct {
	MemberID string                 `json:"member_id"`
	Action   string                 `json:"activity_type"`
	Data     map[string]interface{} `json:"data"`
}
type ActivityLogResponse struct {
	Body       string
	StatusCode int
}
