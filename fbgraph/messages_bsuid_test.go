package fbgraph

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRouteBSUIDRecipient(t *testing.T) {
	t.Run("BSUID in To moves to recipient and to is omitted", func(t *testing.T) {
		m := &MessageObject{To: "US.13491208655302741918", Type: "text"}
		m.routeBSUIDRecipient()
		if m.Recipient != "US.13491208655302741918" {
			t.Fatalf("recipient = %q, want the BSUID", m.Recipient)
		}
		if m.To != "" {
			t.Fatalf("to = %q, want empty", m.To)
		}
		b, _ := json.Marshal(m)
		if strings.Contains(string(b), `"to"`) {
			t.Fatalf("marshaled JSON must omit `to` for a BSUID send: %s", b)
		}
		if !strings.Contains(string(b), `"recipient":"US.13491208655302741918"`) {
			t.Fatalf("marshaled JSON must carry recipient: %s", b)
		}
	})

	t.Run("parent BSUID also routes", func(t *testing.T) {
		m := &MessageObject{To: "US.ENT.11815799212886844830"}
		m.routeBSUIDRecipient()
		if m.Recipient != "US.ENT.11815799212886844830" || m.To != "" {
			t.Fatalf("parent BSUID not routed: to=%q recipient=%q", m.To, m.Recipient)
		}
	})

	t.Run("phone number stays in to", func(t *testing.T) {
		m := &MessageObject{To: "5511987654321"}
		m.routeBSUIDRecipient()
		if m.To != "5511987654321" || m.Recipient != "" {
			t.Fatalf("phone must stay in to: to=%q recipient=%q", m.To, m.Recipient)
		}
	})

	t.Run("explicit recipient is not overwritten", func(t *testing.T) {
		m := &MessageObject{To: "US.123", Recipient: "US.456"}
		m.routeBSUIDRecipient()
		if m.Recipient != "US.456" {
			t.Fatalf("preset recipient overwritten: %q", m.Recipient)
		}
	})
}
