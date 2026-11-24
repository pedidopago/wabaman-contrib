package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pedidopago/wabaman-contrib/rest"
)

// TestPPS10647_GetMessageByID covers the BogusOutput path of doRequest: the
// response body must reach GetMessageByID instead of being decoded into the
// bogus output struct and dropped.
func TestPPS10647_GetMessageByID(t *testing.T) {
	const hostBody = `{"message":{"id":270539328,"object_type":"host","waba_message_id":"wamid.HOST",` +
		`"metadata":{"ad_chatbot_actions":{"cards":[{"product_id":"P1"},{"product_id":"P2"}]},` +
		`"ad_info":{"id":"AD1","expires_at":"2026-11-30"}}}}`
	const clientBody = `{"message":{"id":163069928,"object_type":"client","waba_message_id":"wamid.CLIENT",` +
		`"metadata":{"source":"webhook"}}}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/message" {
			http.NotFound(w, r)
			return
		}
		switch r.URL.Query().Get("id") {
		case "wamid.HOST":
			_, _ = w.Write([]byte(hostBody))
		case "wamid.CLIENT":
			_, _ = w.Write([]byte(clientBody))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"code":"not_found","message":"message not found"}`))
		}
	}))
	defer srv.Close()

	cl := &Client{BaseURL: srv.URL}
	ctx := context.Background()

	t.Run("host message keeps its metadata", func(t *testing.T) {
		msg, err := cl.GetMessageByID(ctx, "wamid.HOST")
		if err != nil {
			t.Fatalf("GetMessageByID: %v", err)
		}
		sent, ok := msg.(*rest.SentMessage)
		if !ok {
			t.Fatalf("got %T, want *rest.SentMessage", msg)
		}
		if sent.GetObjectType() != "host" {
			t.Errorf("object_type = %q, want host", sent.GetObjectType())
		}
		actions, ok := sent.GetMetadata()["ad_chatbot_actions"].(map[string]any)
		if !ok {
			t.Fatalf("metadata.ad_chatbot_actions missing: %#v", sent.GetMetadata())
		}
		if cards, _ := actions["cards"].([]any); len(cards) != 2 {
			t.Errorf("cards = %#v, want 2 entries", actions["cards"])
		}
	})

	t.Run("client message", func(t *testing.T) {
		msg, err := cl.GetMessageByID(ctx, "wamid.CLIENT")
		if err != nil {
			t.Fatalf("GetMessageByID: %v", err)
		}
		recv, ok := msg.(*rest.ReceivedMessage)
		if !ok {
			t.Fatalf("got %T, want *rest.ReceivedMessage", msg)
		}
		if recv.GetMetadata()["source"] != "webhook" {
			t.Errorf("metadata = %#v", recv.GetMetadata())
		}
	})

	t.Run("not found returns an error", func(t *testing.T) {
		if _, err := cl.GetMessageByID(ctx, "wamid.MISSING"); err == nil {
			t.Fatal("expected an error for a 404 response")
		}
	})
}
