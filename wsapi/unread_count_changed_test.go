package wsapi

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestUnreadCountChangedWaitingSinceWire(t *testing.T) {
	t.Run("null when not waiting", func(t *testing.T) {
		b, err := json.Marshal(UnreadCountChanged{PhoneID: 1, ContactID: 2, Count: 0})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(b), `"waiting_since":null`) {
			t.Fatalf("waiting_since must be an explicit null, got %s", b)
		}
	})

	t.Run("round trip when waiting", func(t *testing.T) {
		at := time.Date(2026, 11, 5, 13, 4, 5, 123000000, time.UTC)
		b, err := json.Marshal(UnreadCountChanged{PhoneID: 1, ContactID: 2, Count: 3, WaitingSince: &at})
		if err != nil {
			t.Fatal(err)
		}
		var got UnreadCountChanged
		if err := json.Unmarshal(b, &got); err != nil {
			t.Fatal(err)
		}
		if got.WaitingSince == nil || !got.WaitingSince.Equal(at) {
			t.Fatalf("waiting_since = %v, want %v (wire %s)", got.WaitingSince, at, b)
		}
	})
}
