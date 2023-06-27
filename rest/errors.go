package rest

import (
	"fmt"
	"strings"

	"github.com/pedidopago/wabaman-contrib/fbgraph"
)

type ErrorCode int

// common error codes
const (
	ErrorCodeGenericInvalidParameter ErrorCode = 100
	ErrorCodeGenericBadRequest       ErrorCode = 400
)

type RichError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *RichError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

func NewRichErrorFromError(err error, fallbackCode ...int) *RichError {
	if err == nil {
		return nil
	}
	if e, ok := err.(*RichError); ok {
		return e
	}
	if e, ok := err.(*fbgraph.GraphError); ok {
		emsg := new(strings.Builder)
		emsg.WriteString(e.Message)
		if e.ErrorUserTitle != "" {
			emsg.WriteString("\n")
			emsg.WriteString(e.ErrorUserTitle)
		}
		if e.ErrorUserMsg != "" {
			emsg.WriteString("\n")
			emsg.WriteString(e.ErrorUserMsg)
		}
		return &RichError{
			Code:    ErrorCode(e.Code),
			Message: emsg.String(),
		}
	}
	code := ErrorCode(500)
	if len(fallbackCode) > 0 {
		code = ErrorCode(fallbackCode[0])
	}
	return &RichError{
		Code:    code,
		Message: err.Error(),
	}
}
