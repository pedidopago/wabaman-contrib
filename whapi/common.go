package whapi

import (
	"strconv"
	"time"
)

// The timestamp in the Graph APIs are unix timestamps (seconds since epoch)
// represented as strings.
type Timestamp string

func (t Timestamp) IsEmpty() bool {
	return t == ""
}

func (t Timestamp) ToSeconds() (int64, error) {
	return strconv.ParseInt(string(t), 10, 64)
}

func (t Timestamp) ToTime() (time.Time, error) {
	seconds, err := t.ToSeconds()
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(seconds, 0), nil
}

type DurationSeconds string

func (d DurationSeconds) ToDuration() (time.Duration, error) {
	v, err := strconv.ParseInt(string(d), 10, 64)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Duration(v) * time.Second, nil
}
