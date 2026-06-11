package whapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// The timestamp in the Graph APIs are unix timestamps (seconds since epoch)
// represented as strings.
type Timestamp string

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*t = Timestamp(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("Timestamp: cannot unmarshal %s", string(data))
	}
	*t = Timestamp(n.String())
	return nil
}

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

func (d *DurationSeconds) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*d = DurationSeconds(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("DurationSeconds: cannot unmarshal %s", string(data))
	}
	*d = DurationSeconds(n.String())
	return nil
}

func (d DurationSeconds) ToDuration() (time.Duration, error) {
	v, err := strconv.ParseInt(string(d), 10, 64)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Duration(v) * time.Second, nil
}
