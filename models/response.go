package models

type SuccessResponse struct {
	Success bool   `json:"success"`           // true
	Message string `json:"message,omitempty"` // optional message
	Data    any    `json:"data,omitempty"`    // main payload
}

type ErrorResponse struct {
	Success bool       `json:"success"`         // false
	Message string     `json:"message"`         // human-friendly message
	Error   *ErrorInfo `json:"error,omitempty"` // optional structured error
}

type ErrorInfo struct {
	Code    string `json:"code"`              // e.g., "USER_NOT_FOUND"
	Details string `json:"details,omitempty"` // optional: stack trace, etc.
}

// PaginatedResponse represents a paginated response from the API
type PaginatedResponse struct {
	Success    bool   `json:"success"`           // true
	Message    string `json:"message,omitempty"` // optional message
	Data       any    `json:"data"`
	Page       int    `json:"page"`
	TotalPages int    `json:"total_pages"`
	TotalItems int    `json:"total_items"`
	Limit      int    `json:"limit"`
}
