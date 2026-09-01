package fbgraph

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestSendOutcomeClassification is the whole point of ErrSendOutcomeUnknown:
// the two halves must be told apart. A cancellation before the write is a
// message that provably never left, and tagging it would strand sends the
// caller should simply retry; a cancellation after the write may already be on
// a customer's phone, and not tagging it invites a duplicate. One direction
// costs deliveries, the other costs duplicates, so both are asserted here.
func TestSendOutcomeClassification(t *testing.T) {
	msg := &MessageObject{MessagingProduct: "whatsapp", To: "5511999999999", Type: "text"}

	t.Run("cancelled before the write is known not sent", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("server was reached: the request should never have left")
		}))
		defer srv.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := sendTo(ctx, srv.URL, msg)
		if err == nil {
			t.Fatal("want an error")
		}
		if errors.Is(err, ErrSendOutcomeUnknown) {
			t.Fatalf("tagged unknown, but nothing was written: %v", err)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("want context.Canceled, got %v", err)
		}
	})

	t.Run("cancelled after the write is unknown", func(t *testing.T) {
		// The handler stands in for Meta accepting the message and then being
		// slow to say so -- the window the sentinel exists for.
		got, release := make(chan struct{}), make(chan struct{})
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(got)
			<-release
		}))
		defer srv.Close()
		// Before srv.Close, which waits on this handler. The client's cancel
		// does not reliably reach r.Context() in time to do it for us.
		defer close(release)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			<-got
			cancel()
		}()

		_, err := sendTo(ctx, srv.URL, msg)
		if err == nil {
			t.Fatal("want an error")
		}
		if !errors.Is(err, ErrSendOutcomeUnknown) {
			t.Fatalf("want ErrSendOutcomeUnknown, got %v", err)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cause was dropped: %v", err)
		}
	})

	t.Run("Meta answering is never unknown", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"code":131026,"message":"message undeliverable"}}`))
		}))
		defer srv.Close()

		_, err := sendTo(context.Background(), srv.URL, msg)
		if err == nil {
			t.Fatal("want an error")
		}
		if errors.Is(err, ErrSendOutcomeUnknown) {
			t.Fatalf("a Graph error is a completed round trip, not an unknown one: %v", err)
		}
		if _, ok := AsGraphError(err); !ok {
			t.Fatalf("want a *GraphError, got %v", err)
		}
	})

	t.Run("a truncated 200 body is not a not-sent", func(t *testing.T) {
		// Meta accepted the message and the reply died on the way back: the
		// send happened, only its id was lost. Retrying it sends it twice.
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "64")
			_, _ = w.Write([]byte(`{"messages":`))
		}))
		defer srv.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := sendTo(ctx, srv.URL, msg)
		if err == nil {
			t.Fatal("want an error")
		}
		if !errors.Is(err, ErrSendOutcomeUnknown) {
			t.Fatalf("want ErrSendOutcomeUnknown, got %v", err)
		}
	})
}

// TestSendOutcomeSurvivesAFakeTransport pins the 200 case to the response
// rather than to the write trace. httptrace only fires through a real
// transport, so a substituted RoundTripper that fabricates a reply reports no
// write at all -- and a caller who swapped one in to test that they never
// resend is exactly who cannot afford to be told this send never happened.
func TestSendOutcomeSurvivesAFakeTransport(t *testing.T) {
	c := NewClient("test-token")
	c.HTTPClient = &http.Client{Transport: fakeOK{}}

	_, err := c.SendMessage(context.Background(), "1234567890",
		&MessageObject{MessagingProduct: "whatsapp", To: "5511999999999", Type: "text"})
	if err == nil {
		t.Fatal("want an error")
	}
	if !errors.Is(err, ErrSendOutcomeUnknown) {
		t.Fatalf("want ErrSendOutcomeUnknown, got %v", err)
	}
}

// fakeOK answers 200 with a body that cannot be decoded, without ever putting
// the request on a wire.
type fakeOK struct{}

func (fakeOK) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"messages":`)),
		Request:    req,
	}, nil
}

// sendTo points SendMessage at a test server. GraphAPIVersion cannot redirect
// the host, so the transport does it.
func sendTo(ctx context.Context, base string, msg *MessageObject) (*MessageObjectResult, error) {
	c := NewClient("test-token")
	c.HTTPClient = &http.Client{Transport: redirectTo(base)}

	return c.SendMessage(ctx, "1234567890", msg)
}

type redirectTransport struct{ base string }

func redirectTo(base string) http.RoundTripper { return &redirectTransport{base: base} }

func (t *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	u, err := req.URL.Parse(t.base)
	if err != nil {
		return nil, err
	}
	req = req.Clone(req.Context())
	req.URL.Scheme, req.URL.Host = u.Scheme, u.Host

	return http.DefaultTransport.RoundTrip(req)
}
