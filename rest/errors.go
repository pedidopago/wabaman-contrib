package rest

import (
	"fmt"
	"strings"

	"github.com/pedidopago/wabaman-contrib/v2/fbgraph"
)

type ErrorCode int

// common error codes
const (
	ErrorCodeGenericInvalidParameter ErrorCode = 100
	ErrorCodeGenericBadRequest       ErrorCode = 400
	ErrCodeInternal                  ErrorCode = 500

	// ErrorCodeSendOutcomeUnknown is the REST-side form of
	// fbgraph.ErrSendOutcomeUnknown: the send reached Meta and the answer never
	// came back, so the message may already be on the customer's phone. A client
	// that retries on this code sends the customer the same message twice; it
	// must treat the send as terminal and leave the decision to a person.
	//
	// The value sits outside both the HTTP status range and Meta's Graph error
	// codes (the largest today is 135000), because Code also carries Meta's code
	// verbatim through NewRichErrorFromError and a collision there would make
	// the two indistinguishable. Matching by number rather than by message is
	// the point: messages are for people and change; this does not.
	ErrorCodeSendOutcomeUnknown ErrorCode = 900001
)

type RichError struct {
	HTTPStatus int       `json:"-"`
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
}

func (e *RichError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

func NewRichErrorFromError(err error, statusCode ...int) *RichError {
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
			HTTPStatus: e.HTTPStatusCode,
			Code:       ErrorCode(e.Code),
			Message:    emsg.String(),
		}
	}
	code := 500
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	return &RichError{
		HTTPStatus: code,
		Code:       ErrCodeInternal,
		Message:    err.Error(),
	}
}
