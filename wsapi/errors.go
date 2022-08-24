package wsapi

// Error is the object sent when an error occurs. It is sent when the message
// type is "Error" or "CloseError"
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// ErrorCode is an integer to detect the error returned by the websocket API.
type ErrorCode int

// All valid errors
const (
	ErrorCodeInvalidAuth         ErrorCode = 1
	ErrorCodeInternalError       ErrorCode = 500
	ErrorCodeInternalRedisClosed           = 600
)
