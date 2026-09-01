package fbgraph

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptrace"
	"sync/atomic"
)

var (
	// ErrApplicationRateLimitReached is returned when the Facebook application rate limit is reached (code 4).
	ErrApplicationRateLimitReached = errors.New("application rate limit reached")

	// ErrSendOutcomeUnknown marks a failed send that must not be resent.
	//
	// Sending is not idempotent and has three outcomes, not two: sent, not
	// sent, and unknown. Meta accepts a message the moment the request reaches
	// it, so any failure after that point -- a cancelled context, a timeout, a
	// dropped connection, a body we could not read -- leaves a message that may
	// well have been delivered behind an error that reads like a refusal.
	// Resending on those duplicates a real message to a real customer.
	//
	// Errors from SendMessage and SendMarketingMessage carry this sentinel
	// whenever the request reached the wire, and deliberately do not carry it
	// otherwise: a send that failed before the write is known not to have
	// happened, and is the ordinary case worth retrying. Callers that retry
	// must check for it, because "not sent" is no longer the safe default
	// reading of a non-nil error:
	//
	//	if errors.Is(err, fbgraph.ErrSendOutcomeUnknown) {
	//		// reconcile, dead-letter, or give up -- never resend
	//	}
	//
	// A *GraphError never carries it. Meta answering, however it answered,
	// means the round trip completed and the outcome is known.
	//
	// Accurate classification requires a Client.HTTPClient whose transport
	// honours httptrace; see that field.
	ErrSendOutcomeUnknown = errors.New("send outcome unknown, message may have been delivered")
)

// watchRequestWritten instruments ctx to record whether the request reached the
// wire, and returns the flag to read once the round trip has failed.
//
// The flag is the sent/not-sent boundary. It is set once net/http has finished
// writing the request, which is short of proof that Meta received it -- but
// everything before it is proof that Meta did not, and that is the half worth
// being certain about. Erring the other way costs a duplicate message; this way
// costs a retry we could have skipped.
func watchRequestWritten(ctx context.Context) (context.Context, *atomic.Bool) {
	wrote := new(atomic.Bool)

	// A nil context is rejected with an error further down, by
	// NewRequestWithContext. WithClientTrace would panic on it here instead.
	if ctx == nil {
		return ctx, wrote
	}

	return httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			// A write that failed did not deliver the request, and net/http
			// reports that here rather than by not calling this at all.
			if info.Err == nil {
				wrote.Store(true)
			}
		},
	}), wrote
}

// sendErr wraps err for a send that failed, tagging it ErrSendOutcomeUnknown
// when the request had already reached the wire.
func sendErr(wrote *atomic.Bool, stage string, err error) error {
	if wrote.Load() {
		return fmt.Errorf("%s: %w: %w", stage, ErrSendOutcomeUnknown, err)
	}

	return fmt.Errorf("%s: %w", stage, err)
}
