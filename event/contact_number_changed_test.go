package event

import (
	"testing"
	"time"
)

func TestContactNumberChangedEventJSONRoundTrip(t *testing.T) {
	in := ContactNumberChangedEvent{
		ContactID:  42,
		PhoneID:    7,
		BranchID:   "01BRANCH",
		CustomerID: "01CUSTOMER",
		OldNumber:  "5511900000001",
		NewNumber:  "5511900000002",
		UserID:     "BR.1000000000000001",
		OccurredAt: time.Date(2026, 10, 5, 17, 24, 0, 0, time.UTC),
	}

	var out ContactNumberChangedEvent
	if err := out.FromJSON(in.ToJSON()); err != nil {
		t.Fatal(err)
	}
	if out != in {
		t.Errorf("round trip = %+v, want %+v", out, in)
	}
}
