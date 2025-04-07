package errors

import "time"

type CustomError struct {
	ErrorType string    `json:"type"`
	Message   string    `json:"message"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

func NewCustomError(errorType string, message string, details string, timestamp time.Time) *CustomError {
	return &CustomError{
		ErrorType: errorType,
		Message:   message,
		Details:   details,
		Timestamp: timestamp,
	}
}
