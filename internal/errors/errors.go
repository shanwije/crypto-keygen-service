package errors

import "fmt"

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

var (
	ErrInvalidUserID       = &APIError{Code: 400, Message: "Invalid user ID"}
	ErrUnsupportedNetwork  = &APIError{Code: 400, Message: "Unsupported network"}
	ErrInternalServerError = &APIError{Code: 500, Message: "Internal server error"}
)

func NewAPIError(code int, message string) *APIError {
	return &APIError{Code: code, Message: message}
}
