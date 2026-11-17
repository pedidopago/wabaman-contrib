package fbgraph

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRoutingTestClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c := NewClient("tok")
	c.HTTPClient = srv.Client()
	c.HTTPClient.Transport = rewriteHost{srv.URL, http.DefaultTransport}
	return c
}

func TestThreadControl(t *testing.T) {
	var gotPath, gotMethod string
	var got map[string]any
	c := newRoutingTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &got)
		_, _ = w.Write([]byte(`{"messaging_product":"whatsapp","request_id":"req1"}`))
	})
	resp, err := c.ThreadControl(context.Background(), "123", ThreadControlRequest{
		Action: ThreadControlPass, To: "5511999990000", Metadata: "m",
		ControlPass: &ThreadControlPassTarget{TargetRole: RoutingRoleAIAgent},
	})
	if err != nil || resp.RequestID != "req1" {
		t.Fatalf("resp=%+v err=%v", resp, err)
	}
	if gotMethod != http.MethodPost || gotPath != "/"+DefaultGraphAPIVersion+"/123/thread_control" {
		t.Errorf("%s %s", gotMethod, gotPath)
	}
	cp, _ := got["control_pass"].(map[string]any)
	if got["messaging_product"] != "whatsapp" || got["action"] != "pass" || got["to"] != "5511999990000" ||
		got["metadata"] != "m" || cp["target_role"] != "ai_agent" || got["recipient"] != nil {
		t.Errorf("body = %v", got)
	}

	// release carries no control_pass.
	got = nil
	if _, err := c.ThreadControl(context.Background(), "123", ThreadControlRequest{Action: ThreadControlRelease, Recipient: "BR.1"}); err != nil {
		t.Fatal(err)
	}
	if _, ok := got["control_pass"]; ok || got["recipient"] != "BR.1" || got["to"] != nil {
		t.Errorf("release body = %v", got)
	}
}

func TestThreadControlNeedsOneRecipient(t *testing.T) {
	c := NewClient("tok")
	for _, r := range []ThreadControlRequest{{Action: ThreadControlTake}, {Action: ThreadControlTake, To: "1", Recipient: "BR.1"}} {
		if _, err := c.ThreadControl(context.Background(), "123", r); err == nil {
			t.Errorf("%+v: want an error", r)
		}
	}
}

func TestGetRoutingRoles(t *testing.T) {
	c := newRoutingTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/"+DefaultGraphAPIVersion+"/123/routing_roles" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"roles":["customer_service","escalation"],"is_standby":true}`))
	})
	got, err := c.GetRoutingRoles(context.Background(), "123")
	if err != nil || len(got.Roles) != 2 || got.Roles[1] != RoutingRoleEscalation || !got.IsStandby {
		t.Fatalf("got=%+v err=%v", got, err)
	}
}

func TestGetRoutingRolesGraphError(t *testing.T) {
	c := newRoutingTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"message":"Conversation routing is not enabled for this phone number.","code":100}}`))
	})
	_, err := c.GetRoutingRoles(context.Background(), "123")
	var ge *GraphError
	if !errors.As(err, &ge) || ge.Code != 100 {
		t.Fatalf("err = %v, want a GraphError code 100", err)
	}
}
