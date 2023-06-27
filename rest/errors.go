package rest

type ErrorCode int

// common error codes
const (
	ErrorCodeGenericInvalidParameter ErrorCode = 400
)

type RichError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}
