package errors

import (
	"fmt"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

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
	ErrInvalidUserID       = &KeyGenError{Code: 400, Message: "userId must be a positive integer"}
	ErrNetworkRequired     = &KeyGenError{Code: 400, Message: "Network is required"}
)

func NewKeyGenError(code int, message string) *KeyGenError {
	return &KeyGenError{Code: code, Message: message}
}

func FormatValidationError(err error) string {
	var sb strings.Builder
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "UserID":
			switch err.Tag() {
			case "required":
				sb.WriteString("UserID is required. ")
			case "gt":
				sb.WriteString("UserID must be greater than 0. ")
			}
		case "Network":
			if err.Tag() == "required" {
				sb.WriteString("Network is required. ")
			}
		}
	}
	return strings.TrimSpace(sb.String())
}
