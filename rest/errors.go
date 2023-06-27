package rest

import "fmt"

type ErrorCode int

// common error codes
const (
	ErrorCodeGenericInvalidParameter ErrorCode = 400
)

type RichError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *RichError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}
