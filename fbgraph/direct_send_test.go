package fbgraph

import (
	"encoding/json"
	"testing"
)

func TestDirectSendJSON(t *testing.T) {
	m := MessageObject{
		MessagingProduct: "whatsapp",
		To:               "5511999999999",
		Type:             "text",
		Text:             &TextObject{Body: "Pedido enviado"},
		Category:         DirectSendUtility,
		DirectSendConfig: &DirectSendConfig{TemplateName: "order_shipped"},
		TTLSeconds:       600,
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got["category"] != "utility" || got["ttl_seconds"] != float64(600) ||
		got["direct_send_config"].(map[string]any)["template_name"] != "order_shipped" {
		t.Errorf("json = %s", b)
	}

	got = nil
	b, _ = json.Marshal(MessageObject{Type: "text", Text: &TextObject{Body: "x"}})
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"category", "direct_send_config", "ttl_seconds"} {
		if _, ok := got[k]; ok {
			t.Errorf("plain message carries %s: %s", k, b)
		}
	}
}
