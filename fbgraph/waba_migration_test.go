package fbgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// rewriteTransport sends every request to the test server instead of
// graph.facebook.com, so the client's hardcoded URLs can be asserted on.
type rewriteTransport struct {
	target *url.URL
	gotReq *http.Request
	gotBdy []byte
}

func (rt *rewriteTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Body != nil {
		b, _ := io.ReadAll(r.Body)
		rt.gotBdy = b
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(b))
	}
	rt.gotReq = r.Clone(r.Context())

	r.URL.Scheme = rt.target.Scheme
	r.URL.Host = rt.target.Host
	return http.DefaultTransport.RoundTrip(r)
}

func newTestClient(t *testing.T, status int, body string) (*Client, *rewriteTransport) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parse test server url: %v", err)
	}

	rt := &rewriteTransport{target: u}
	return &Client{
		HTTPClient:      &http.Client{Transport: rt},
		AccessToken:     "test_token",
		GraphAPIVersion: "v23.0",
	}, rt
}

func TestSetPaymentMethodMigrationIntent(t *testing.T) {
	cl, rt := newTestClient(t, http.StatusOK,
		`{"migration_id":"mig_01HZYK3ABCDEF4567890","migration_status":"INITIATED"}`)

	out, err := cl.SetPaymentMethodMigrationIntent(context.Background(), "431859353351325",
		MigrationIntentRequest{Currency: "BRL"})
	if err != nil {
		t.Fatalf("intent: %v", err)
	}

	if got, want := rt.gotReq.URL.Path, "/v23.0/431859353351325/set_payment_method_migration_intent"; got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
	if rt.gotReq.Method != http.MethodPost {
		t.Errorf("method = %q, want POST", rt.gotReq.Method)
	}
	if got, want := rt.gotReq.Header.Get("Authorization"), "Bearer test_token"; got != want {
		t.Errorf("auth = %q, want %q", got, want)
	}

	// extended_credit_id must be absent for credit-card billing: sending an
	// empty one would have Meta look for a Line of Credit that isn't there.
	var sent map[string]any
	if err := json.Unmarshal(rt.gotBdy, &sent); err != nil {
		t.Fatalf("unmarshal sent body: %v", err)
	}
	if sent["currency"] != "BRL" {
		t.Errorf("currency = %v, want BRL", sent["currency"])
	}
	if _, present := sent["extended_credit_id"]; present {
		t.Error("extended_credit_id sent when empty")
	}

	if out.MigrationID != "mig_01HZYK3ABCDEF4567890" || out.MigrationStatus != MigrationStatusInitiated {
		t.Errorf("out = %+v", out)
	}
}

func TestSetPaymentMethodMigrationIntentWithCreditLine(t *testing.T) {
	cl, rt := newTestClient(t, http.StatusOK, `{"migration_id":"mig_x","migration_status":"INITIATED"}`)

	if _, err := cl.SetPaymentMethodMigrationIntent(context.Background(), "123",
		MigrationIntentRequest{Currency: "BRL", ExtendedCreditID: "987654321098765"}); err != nil {
		t.Fatalf("intent: %v", err)
	}

	var sent map[string]any
	if err := json.Unmarshal(rt.gotBdy, &sent); err != nil {
		t.Fatalf("unmarshal sent body: %v", err)
	}
	if sent["extended_credit_id"] != "987654321098765" {
		t.Errorf("extended_credit_id = %v", sent["extended_credit_id"])
	}
}

func TestGetMigrationStatusWithDestination(t *testing.T) {
	cl, rt := newTestClient(t, http.StatusOK, `{
		"status":"READY_TO_COMPLETE",
		"destination_waba":{"id":"998877665544332","name":"Acme WABA (BRL)","currency":"BRL",
			"timezone_id":"1","message_template_namespace":"1234abcd_namespace"},
		"id":"mig_01HZYK3ABCDEF4567890"}`)

	out, err := cl.GetMigrationStatus(context.Background(), "mig_01HZYK3ABCDEF4567890")
	if err != nil {
		t.Fatalf("status: %v", err)
	}

	if got, want := rt.gotReq.URL.Path, "/v23.0/mig_01HZYK3ABCDEF4567890"; got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
	if out.Status != MigrationStatusReadyToComplete {
		t.Errorf("status = %q", out.Status)
	}
	if out.DestinationWABA == nil {
		t.Fatal("destination_waba nil when present in payload")
	}
	if out.DestinationWABA.ID != "998877665544332" || out.DestinationWABA.Currency != "BRL" {
		t.Errorf("destination = %+v", out.DestinationWABA)
	}
	if out.DestinationWABA.TimezoneID != "1" {
		t.Errorf("timezone_id = %q, want 1 (needed for volume-tier month boundaries)", out.DestinationWABA.TimezoneID)
	}
}

// Early in a migration Meta has not created the clone yet; the absent
// destination must read as nil rather than a zero-valued WABA that callers
// might mistake for a real one.
func TestGetMigrationStatusWithoutDestination(t *testing.T) {
	cl, _ := newTestClient(t, http.StatusOK, `{"status":"IN_PROGRESS","id":"mig_x"}`)

	out, err := cl.GetMigrationStatus(context.Background(), "mig_x")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if out.DestinationWABA != nil {
		t.Errorf("destination_waba = %+v, want nil", out.DestinationWABA)
	}
	if out.Status != MigrationStatusInProgress {
		t.Errorf("status = %q", out.Status)
	}
}

func TestResumeMigration(t *testing.T) {
	cl, rt := newTestClient(t, http.StatusOK, `{
		"status":"COMPLETED",
		"destination_waba":{"id":"998877665544332","currency":"BRL"},
		"id":"mig_x"}`)

	out, err := cl.ResumeMigration(context.Background(), "mig_x")
	if err != nil {
		t.Fatalf("resume: %v", err)
	}

	if got, want := rt.gotReq.URL.Path, "/v23.0/mig_x/resume_migration"; got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
	if rt.gotReq.Method != http.MethodPost {
		t.Errorf("method = %q, want POST", rt.gotReq.Method)
	}
	if out.Status != MigrationStatusCompleted {
		t.Errorf("status = %q", out.Status)
	}
}

func TestGetWABAPhoneNumbers(t *testing.T) {
	cl, rt := newTestClient(t, http.StatusOK, `{"data":[
		{"id":"813500388509960","display_phone_number":"+55 11 3271-0305","verified_name":"Suprema Farma","quality_rating":"GREEN"},
		{"id":"700388133151639","display_phone_number":"+55 44 9857-0329","verified_name":"Manipulacao SP","quality_rating":"GREEN"}]}`)

	out, err := cl.GetWABAPhoneNumbers(context.Background(), "629535923565046")
	if err != nil {
		t.Fatalf("phone numbers: %v", err)
	}

	if got, want := rt.gotReq.URL.Path, "/v23.0/629535923565046/phone_numbers"; got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
	if len(out) != 2 {
		t.Fatalf("len = %d, want 2", len(out))
	}
	if out[0].ID != "813500388509960" || out[0].DisplayPhoneNumber != "+55 11 3271-0305" {
		t.Errorf("first = %+v", out[0])
	}
}

func TestGetWABAInfo(t *testing.T) {
	cl, rt := newTestClient(t, http.StatusOK,
		`{"id":"629535923565046","currency":"USD","name":"Suprema","timezone_id":"1"}`)

	out, err := cl.GetWABAInfo(context.Background(), "629535923565046")
	if err != nil {
		t.Fatalf("waba info: %v", err)
	}
	if got, want := rt.gotReq.URL.Path, "/v23.0/629535923565046"; got != want {
		t.Errorf("path = %q, want %q", got, want)
	}
	if got := rt.gotReq.URL.Query().Get("fields"); got != "id,currency,name,timezone_id" {
		t.Errorf("fields = %q", got)
	}
	if out.Currency != "USD" {
		t.Errorf("currency = %q", out.Currency)
	}
}

// A Meta error must surface as a GraphError carrying the code, so callers can
// tell an ineligible WABA from a transient failure.
func TestMigrationGraphError(t *testing.T) {
	cl, _ := newTestClient(t, http.StatusBadRequest,
		`{"error":{"message":"(#100) Coexistence WABAs are not eligible","type":"OAuthException","code":100}}`)

	_, err := cl.SetPaymentMethodMigrationIntent(context.Background(), "123",
		MigrationIntentRequest{Currency: "BRL"})
	if err == nil {
		t.Fatal("expected error")
	}

	ge, ok := AsGraphError(err)
	if !ok {
		t.Fatalf("err = %v, want a *GraphError", err)
	}
	if ge.Code != 100 {
		t.Errorf("code = %d, want 100", ge.Code)
	}
}

// A bare 5xx with no Meta error body must still classify, preserving the status
// so a retry decision is possible.
func TestMigrationServerErrorPreservesStatus(t *testing.T) {
	cl, _ := newTestClient(t, http.StatusBadGateway, `<html>gateway error</html>`)

	_, err := cl.ResumeMigration(context.Background(), "mig_x")
	if err == nil {
		t.Fatal("expected error")
	}

	ge, ok := AsGraphError(err)
	if !ok {
		t.Fatalf("err = %v, want a *GraphError", err)
	}
	if ge.HTTPStatusCode != http.StatusBadGateway {
		t.Errorf("HTTPStatusCode = %d, want 502", ge.HTTPStatusCode)
	}
}

// A truncated phone-number list is indistinguishable from "these numbers did
// not migrate" -- a conclusion the caller acts on -- so paging must complete.
func TestGetWABAPhoneNumbersPaginates(t *testing.T) {
	page := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		page++
		if r.URL.Query().Get("after") == "" {
			_, _ = io.WriteString(w, `{"data":[{"id":"1","display_phone_number":"+55 11 1111-1111"}],
				"paging":{"cursors":{"after":"CUR1"},"next":"https://graph.facebook.com/next"}}`)
			return
		}
		// Last page still carries a cursor but no `next`.
		_, _ = io.WriteString(w, `{"data":[{"id":"2","display_phone_number":"+55 11 2222-2222"}],
			"paging":{"cursors":{"after":"CUR2"}}}`)
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cl := &Client{
		HTTPClient:      &http.Client{Transport: &rewriteTransport{target: u}},
		AccessToken:     "t",
		GraphAPIVersion: "v23.0",
	}

	out, err := cl.GetWABAPhoneNumbers(context.Background(), "waba")
	if err != nil {
		t.Fatalf("phone numbers: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("len = %d, want 2 (both pages)", len(out))
	}
	if out[0].ID != "1" || out[1].ID != "2" {
		t.Errorf("ids = %q,%q", out[0].ID, out[1].ID)
	}
	if page != 2 {
		t.Errorf("requested %d pages, want exactly 2 -- a cursor on the last page must not cause another round trip", page)
	}
}
