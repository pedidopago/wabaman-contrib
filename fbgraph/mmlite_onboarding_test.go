package fbgraph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMMLiteOnboardingStatus(t *testing.T) {
	tests := []struct {
		name string
		body string
		want MMLiteOnboardingStatus
	}{
		{"waba onboarded", `{"marketing_messages_onboarding_status":"ONBOARDED","id":"1"}`, MMLiteOnboarded},
		{"portfolio signed", `{"marketing_messages_onboarding_status":"TERM_OF_SERVICE_SIGNED","id":"1"}`, MMLiteTermsOfServiceSigned},
		// Absent field must NOT become NOT_STARTED: "Meta will not tell us" is
		// not the same as "nobody asked them", and one of those is a claim
		// about a customer's legal acceptance.
		{"field absent", `{"id":"1"}`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotFields string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotFields = r.URL.Query().Get("fields")
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			c := NewClient("tok")
			c.HTTPClient = srv.Client()
			// Point the hardcoded graph host at the stub.
			c.HTTPClient.Transport = rewriteHost{srv.URL, http.DefaultTransport}

			got, err := c.GetMMLiteOnboardingStatus(context.Background(), "123")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("status = %q, want %q", got, tt.want)
			}
			if gotFields != mmLiteOnboardingStatusField {
				t.Fatalf("fields = %q, want exactly %q (requesting more risks losing the whole call on a token that cannot see one of them)",
					gotFields, mmLiteOnboardingStatusField)
			}
		})
	}
}

func TestGetMMLiteOnboardingStatusHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"message":"nope","code":100}}`))
	}))
	defer srv.Close()

	c := NewClient("tok")
	c.HTTPClient = srv.Client()
	c.HTTPClient.Transport = rewriteHost{srv.URL, http.DefaultTransport}

	if _, err := c.GetMMLiteOnboardingStatus(context.Background(), "123"); err == nil {
		t.Fatal("expected an error on a non-200 response")
	}
}

// rewriteHost sends every request to the stub server instead of graph.facebook.com,
// preserving path and query so assertions on them stay meaningful.
type rewriteHost struct {
	base string
	next http.RoundTripper
}

func (r rewriteHost) RoundTrip(req *http.Request) (*http.Response, error) {
	u := *req.URL
	stub, err := http.NewRequest(req.Method, r.base+u.Path+"?"+u.RawQuery, req.Body)
	if err != nil {
		return nil, err
	}
	stub.Header = req.Header
	return r.next.RoundTrip(stub.WithContext(req.Context()))
}
