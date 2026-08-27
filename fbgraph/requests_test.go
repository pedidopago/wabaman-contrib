package fbgraph

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	neturl "net/url"
	"strings"
	"testing"
	"time"
)

// The methods in templates.go, conversions.go and token.go had no tests at all,
// so the request each one builds was entirely unpinned: a wrong verb, a dropped
// query string, a nil body or a corrupted grant_type would break every call
// against Meta while the suite stayed green. These pin the request line, the
// credentials, the query and the decoded response for each.

type capturedRequest struct {
	method string
	path   string
	query  neturl.Values
	auth   string
	ctype  string
	host   string
	body   string
}

// stubGraph answers every request with respBody and records what it received.
func stubGraph(t *testing.T, respBody string) (*Client, <-chan capturedRequest) {
	t.Helper()

	got := make(chan capturedRequest, 4)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		got <- capturedRequest{
			method: r.Method,
			path:   r.URL.Path,
			query:  r.URL.Query(),
			auth:   r.Header.Get("Authorization"),
			ctype:  r.Header.Get("Content-Type"),
			host:   r.Header.Get("X-Orig-Host"),
			body:   string(body),
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	}))
	t.Cleanup(srv.Close)

	c := NewClient("token")
	c.GraphAPIVersion = "v99.0"
	c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

	return c, got
}

// await fails the test rather than hanging the package for the full go test
// timeout when a regression stops the request from reaching the stub.
func await(t *testing.T, got <-chan capturedRequest) capturedRequest {
	t.Helper()

	select {
	case r := <-got:
		return r
	case <-time.After(5 * time.Second):
		t.Fatal("the stub never received the request")
	}

	return capturedRequest{}
}

// checkAuth pins the bearer credential. Dropping it turns every call into an
// OAuthException 190 in production, which is exactly the class of regression
// these tests exist to catch.
func (r capturedRequest) checkAuth(t *testing.T) {
	t.Helper()

	if r.auth != "Bearer token" {
		t.Errorf("Authorization = %q, want Bearer token", r.auth)
	}
}

func (r capturedRequest) checkLine(t *testing.T, method, path string) {
	t.Helper()

	if r.method != method {
		t.Errorf("method = %s, want %s", r.method, method)
	}
	if r.path != path {
		t.Errorf("path = %q, want %q", r.path, path)
	}
}

func TestTemplateRequests(t *testing.T) {
	t.Run("GetMessageTemplate", func(t *testing.T) {
		c, got := stubGraph(t, `{"id":"TPL-1","name":"welcome","status":"APPROVED"}`)

		tpl, err := c.GetMessageTemplate(context.Background(), "TPL-1")
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if tpl.ID != "TPL-1" || tpl.Name != "welcome" {
			t.Errorf("response not decoded: %+v", tpl)
		}

		r := await(t, got)
		r.checkLine(t, http.MethodGet, "/v99.0/TPL-1")
		if r.auth != "Bearer token" {
			t.Errorf("Authorization = %q", r.auth)
		}
	})

	t.Run("GetMessageTemplates", func(t *testing.T) {
		c, got := stubGraph(t, `{"data":[{"id":"TPL-1","name":"welcome"}]}`)

		res, err := c.GetMessageTemplates(context.Background(), GetMessageTemplatesParameters{
			WhatsAppBusinessAccountID: "WABA-1",
			Limit:                     25,
		})
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if len(res.Data) != 1 || res.Data[0].ID != "TPL-1" {
			t.Errorf("response not decoded: %+v", res)
		}

		r := await(t, got)
		r.checkLine(t, http.MethodGet, "/v99.0/WABA-1/message_templates")
		r.checkAuth(t)
		// The query carries the paging and field selection; dropping it changes
		// which templates Meta returns.
		if len(r.query) == 0 {
			t.Error("the query string was dropped entirely")
		}
		if r.query.Get("limit") != "25" {
			t.Errorf("limit = %q, want 25", r.query.Get("limit"))
		}
		// Without fields Meta returns its default subset, and status,
		// components and quality_score come back empty -- a live failure the
		// caller sees as "the template has no content".
		if f := r.query.Get("fields"); !strings.Contains(f, "components") || !strings.Contains(f, "status") {
			t.Errorf("fields = %q, want components and status selected", f)
		}
	})

	t.Run("CreateMessageTemplate", func(t *testing.T) {
		c, got := stubGraph(t, `{"id":"TPL-NEW"}`)

		id, err := c.CreateMessageTemplate(context.Background(), "WABA-1", NewMessageTemplate{
			MessageTemplate: MessageTemplate{Name: "welcome", Language: "pt_BR"},
		})
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if id != "TPL-NEW" {
			t.Errorf("id = %q", id)
		}

		r := await(t, got)
		r.checkLine(t, http.MethodPost, "/v99.0/WABA-1/message_templates")
		r.checkAuth(t)
		if r.host != "graph.facebook.com" {
			t.Errorf("host = %q, want graph.facebook.com", r.host)
		}
		if r.ctype != "application/json" {
			t.Errorf("Content-Type = %q", r.ctype)
		}
		// A nil body creates nothing; the template must actually be sent.
		var sent map[string]any
		if err := json.Unmarshal([]byte(r.body), &sent); err != nil {
			t.Fatalf("body was not JSON (%q): %v", r.body, err)
		}
		if sent["name"] != "welcome" {
			t.Errorf("template never reached the body: %v", sent)
		}
	})

	t.Run("DeleteMessageTemplate", func(t *testing.T) {
		c, got := stubGraph(t, `{"success":true}`)

		if err := c.DeleteMessageTemplate(context.Background(), "WABA-1", "welcome"); err != nil {
			t.Fatalf("failed: %v", err)
		}

		r := await(t, got)
		// DELETE, not GET: the wrong verb silently deletes nothing.
		r.checkLine(t, http.MethodDelete, "/v99.0/WABA-1/message_templates")
		if r.query.Get("name") != "welcome" {
			t.Errorf("name = %q, want the template name", r.query.Get("name"))
		}
	})

	t.Run("UnarchiveMessageTemplates", func(t *testing.T) {
		c, got := stubGraph(t, `{"unarchived_templates":["TPL-1"],"failed_templates":{"TPL-2":"nope"}}`)

		un, failed, err := c.UnarchiveMessageTemplates(context.Background(), "WABA-1", []string{"TPL-1", "TPL-2"})
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if len(un) != 1 || un[0] != "TPL-1" {
			t.Errorf("unarchived = %v", un)
		}
		if failed["TPL-2"] != "nope" {
			t.Errorf("failed = %v", failed)
		}

		r := await(t, got)
		r.checkLine(t, http.MethodPost, "/WABA-1/message_templates/unarchive")
		r.checkAuth(t)
		// Deliberately api.facebook.com and unversioned: the versioned
		// graph.facebook.com endpoint answers 2500 for this call.
		if r.host != "api.facebook.com" {
			t.Errorf("host = %q, want api.facebook.com", r.host)
		}
	})
}

func TestConversionRequests(t *testing.T) {
	t.Run("CreateDataset", func(t *testing.T) {
		c, got := stubGraph(t, `{"id":"DS-1"}`)

		id, err := c.CreateDataset(context.Background(), "WABA-1")
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if id != "DS-1" {
			t.Errorf("id = %q", id)
		}

		r := await(t, got)
		// POST creates; GET would silently return the existing list instead.
		r.checkLine(t, http.MethodPost, "/v99.0/WABA-1/dataset")
		if r.auth != "Bearer token" {
			t.Errorf("Authorization = %q", r.auth)
		}
	})

	t.Run("GetDatasetID", func(t *testing.T) {
		c, got := stubGraph(t, `{"data":[{"id":"DS-1"}]}`)

		id, err := c.GetDatasetID(context.Background(), "WABA-1")
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if id != "DS-1" {
			t.Errorf("id = %q", id)
		}

		r := await(t, got)
		r.checkLine(t, http.MethodGet, "/v99.0/WABA-1/dataset")
	})

	t.Run("SendConversionEvents", func(t *testing.T) {
		c, got := stubGraph(t, `{"events_received":2}`)

		n, err := c.SendConversionEvents(context.Background(), "DS-1", []ConversionEvent{
			{EventName: "Purchase"}, {EventName: "Lead"},
		})
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if n != 2 {
			t.Errorf("events_received = %d, want 2", n)
		}

		r := await(t, got)
		r.checkLine(t, http.MethodPost, "/v99.0/DS-1/events")
		r.checkAuth(t)
		var sent map[string]any
		if err := json.Unmarshal([]byte(r.body), &sent); err != nil {
			t.Fatalf("body was not JSON (%q): %v", r.body, err)
		}
		// Meta rejects a batch whose events are not under `data`.
		events, ok := sent["data"].([]any)
		if !ok || len(events) != 2 {
			t.Errorf("events not sent under data: %v", sent)
		}
	})
}

func TestTokenRequests(t *testing.T) {
	t.Run("DebugToken", func(t *testing.T) {
		c, got := stubGraph(t, `{"data":{"app_id":"APP","is_valid":true}}`)

		info, err := c.DebugToken(context.Background(), "input-token")
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if !info.IsValid || info.AppID != "APP" {
			t.Errorf("response not decoded: %+v", info)
		}

		r := await(t, got)
		r.checkLine(t, http.MethodGet, "/v99.0/debug_token")
		if r.query.Get("input_token") != "input-token" {
			t.Errorf("input_token = %q", r.query.Get("input_token"))
		}
		// Debug uses the app token in the query, not a bearer header.
		if r.query.Get("access_token") != "token" {
			t.Errorf("access_token = %q", r.query.Get("access_token"))
		}
	})

	t.Run("NewPermanentAccessToken", func(t *testing.T) {
		c, got := stubGraph(t, `{"access_token":"PERMANENT"}`)

		tok, err := c.NewPermanentAccessToken(context.Background(), "app-id", "app-secret", "temp")
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if tok != "PERMANENT" {
			t.Errorf("token = %q", tok)
		}

		r := await(t, got)
		// v13.0 is hardcoded here, unlike every other method, which honours
		// GraphAPIVersion. Pinned so the oddity is visible rather than implicit.
		r.checkLine(t, http.MethodGet, "/v13.0/oauth/access_token")
		// grant_type is what makes this an exchange; a corrupted value returns
		// a short-lived token that expires in hours.
		if r.query.Get("grant_type") != "fb_exchange_token" {
			t.Errorf("grant_type = %q, want fb_exchange_token", r.query.Get("grant_type"))
		}
		if r.query.Get("client_id") != "app-id" || r.query.Get("fb_exchange_token") != "temp" {
			t.Errorf("exchange parameters wrong: %v", r.query)
		}
	})
}
