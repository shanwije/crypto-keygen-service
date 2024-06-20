package errors

import "fmt"

type KeyGenError struct {
	Code    int
	Message string
}

func (e *KeyGenError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

var (
	ErrUnsupportedNetwork  = &KeyGenError{Code: 400, Message: "Unsupported network"}
	ErrInternalServerError = &KeyGenError{Code: 500, Message: "Internal server error"}
)

func NewKeyGenError(code int, message string) *KeyGenError {
	return &KeyGenError{Code: code, Message: message}
}
