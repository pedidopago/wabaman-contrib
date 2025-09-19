package fbgraph

import "errors"

var (
	// ErrApplicationRateLimitReached is returned when the Facebook application rate limit is reached (code 4).
	ErrApplicationRateLimitReached = errors.New("application rate limit reached")
)
