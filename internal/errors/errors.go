package errors

import "fmt"

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{Code: code, Message: message}
}
