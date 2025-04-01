package errors

import "time"

type ErrorObject struct {
	ErrorType string    `json:"type"`
	Message   string    `json:"message"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
}

func NewErrorObject(errorType string, message string, details string, timestamp time.Time, path string) *ErrorObject {
	return &ErrorObject{
		ErrorType: errorType,
		Message:   message,
		Details:   details,
		Timestamp: timestamp,
		Path:      path,
	}
}
