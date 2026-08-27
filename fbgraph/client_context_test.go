package fbgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// spyTransport stands in for the network so this test never issues a real
// request, and records whether the request it was handed carried the caller's
// context. That flag is the contract: http.Client hands the request to its
// transport even when the context is already dead, so being reached proves
// nothing, and the returned error proves nothing either -- a method that builds
// its own context completes a real round trip and returns Meta's OAuth error,
// which reads like a failure just as a cancellation does. Whether req.Context()
// is the dead one is the only thing that actually separates the two.
type spyTransport struct {
	sawLiveContext atomic.Bool
}

func (t *spyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	err := req.Context().Err()
	if err == nil {
		t.sawLiveContext.Store(true)
	}

	// net/http closes the request body when a round trip fails. UploadMedia
	// streams its body through an io.Pipe, so without this the writer blocks
	// forever on a pipe nobody is reading.
	if req.Body != nil {
		_ = req.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	return nil, errors.New("spyTransport: request did not carry the caller's context")
}

// TestClientMethodsHonorCanceledContext pins the contract these five methods
// gained: the request carries the caller's context, so a dead context fails the
// call before any network I/O instead of running to the client timeout. The
// refusing transport is what makes that an assertion rather than an assumption
// -- without it a method that ignores the context quietly completes a real
// round trip to Meta and returns a non-nil error, which every check here would
// otherwise accept as success.
func TestClientMethodsHonorCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	spy := &spyTransport{}
	c := NewClient("token")
	c.HTTPClient = &http.Client{Transport: spy}
	msg := &MessageObject{MessagingProduct: "whatsapp", To: "5511999999999", Type: "text"}

	tests := []struct {
		name string
		call func() error
	}{
		{"SendMessage", func() error {
			_, err := c.SendMessage(ctx, "phone-id", msg)
			return err
		}},
		{"SendMarketingMessage", func() error {
			_, err := c.SendMarketingMessage(ctx, "phone-id", msg)
			return err
		}},
		{"GetMedia", func() error {
			_, err := c.GetMedia(ctx, "media-id")
			return err
		}},
		{"DownloadMedia", func() error {
			_, err := c.DownloadMedia(ctx, &GetMediaResult{URL: "https://graph.facebook.com/media"}, &bytes.Buffer{})
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy.sawLiveContext.Store(false)

			err := tt.call()
			if err == nil {
				t.Fatal("expected an error on a canceled context, got nil")
			}
			if spy.sawLiveContext.Load() {
				t.Fatalf("%s built its own context: the request reached the transport without the caller's cancellation (err=%v)", tt.name, err)
			}
			if !errors.Is(err, context.Canceled) {
				t.Errorf("expected context.Canceled, got %v", err)
			}
		})
	}

	// UploadMedia streams through an io.Pipe, so a canceled context can surface
	// either as the failed request or as the write end noticing the closed
	// pipe. Both are correct. What is not correct is leaving the goroutine that
	// issues the request parked on its hand-off channel: that path fires on
	// every cancellation, so an unbuffered channel would leak one goroutine per
	// canceled upload. Asserting only err != nil walks straight past that.
	t.Run("UploadMedia", func(t *testing.T) {
		spy.sawLiveContext.Store(false)
		before := runtime.NumGoroutine()

		_, err := c.UploadMedia(ctx, "phone-id", "image/png", strings.NewReader("x"), 1, "f.png")
		if err == nil {
			t.Fatal("expected an error on a canceled context, got nil")
		}
		// A dead context must reach the caller as such. net/http closes the
		// pipe reader plainly when the request fails, so without help the
		// writer's io.ErrClosedPipe is all the caller sees -- and a caller that
		// branches on context.Canceled to decide "retry" versus "this media is
		// broken" gets it wrong in the one case this change exists for.
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled to survive the multipart write, got %v", err)
		}
		if spy.sawLiveContext.Load() {
			t.Fatalf("UploadMedia built its own context: the request reached the transport without the caller's cancellation (err=%v)", err)
		}

		// The request goroutine unwinds independently of the return, so give it
		// a moment before counting.
		for range 20 {
			if runtime.NumGoroutine() <= before {
				return
			}
			time.Sleep(50 * time.Millisecond)
		}

		buf := make([]byte, 1<<16)
		t.Errorf("UploadMedia leaked a goroutine on a canceled context:\n%s", buf[:runtime.Stack(buf, true)])
	})
}

// TestDownloadMediaAbortsInFlightOnCancel is the reason these signatures
// changed. A context-free request ignores cancellation and holds the caller
// until the client timeout (120s on DefaultHTTPClient), which is what kept a
// draining pod waiting on work already doomed to fail at its next step.
// DownloadMedia is the one method whose URL the caller supplies, so it is the
// one that can be pointed at a test server.
func TestDownloadMediaAbortsInFlightOnCancel(t *testing.T) {
	// Buffered rather than closed: a transport-level retry can enter the
	// handler twice, and closing an already-closed channel would panic the
	// whole test binary.
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case started <- struct{}{}:
		default:
		}
		select {
		case <-r.Context().Done(): // never answers
		case <-release:
		}
	}))
	// Order matters: deferred calls run LIFO, so close(release) frees the
	// handler before srv.Close waits on it. Reversed, a regression that leaves
	// the handler parked would block Close and kill the run on the package
	// timeout instead of failing at the deadline below.
	defer srv.Close()
	defer close(release)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewClient("token")

	done := make(chan error, 1)
	go func() {
		_, err := c.DownloadMedia(ctx, &GetMediaResult{URL: srv.URL, MimeType: "image/png"}, &bytes.Buffer{})
		done <- err
	}()

	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("the handler was never reached: DownloadMedia returned without issuing the request")
	}
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("DownloadMedia ignored the cancellation and kept waiting on the server")
	}
}

// TestUploadMediaFailsFastOnUnusableContext covers the branch that only became
// reachable when UploadMedia started building its request with the context:
// NewRequestWithContext rejects a nil context, and the request goroutine bails
// before anything reads the pipe. Without closing the read end, the multipart
// writer parks on it and UploadMedia never returns at all -- a caller wedged
// forever, which is strictly worse than the stall this change set out to fix.
func TestUploadMediaFailsFastOnUnusableContext(t *testing.T) {
	var ctx context.Context // nil: what a half-migrated call site looks like

	c := NewClient("token")
	c.HTTPClient = &http.Client{Transport: &spyTransport{}}

	done := make(chan error, 1)
	go func() {
		_, err := c.UploadMedia(ctx, "phone-id", "image/png", strings.NewReader("payload"), 7, "f.png")
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected an error for a nil context, got nil")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("UploadMedia never returned: the request goroutine left the multipart writer parked on the pipe")
	}
}

// errFaulty is what a source stream returns when it dies partway through --
// a truncated S3 object, a dropped disk read.
var errFaulty = errors.New("source stream failed midway")

type faultyReader struct{ left int }

func (r *faultyReader) Read(p []byte) (int, error) {
	if r.left <= 0 {
		return 0, errFaulty
	}

	n := min(len(p), r.left)
	r.left -= n

	return n, nil
}

// bodyReadingTransport consumes the request body the way a real transport does
// and reports how the read ended: cleanly, or with an error.
type bodyReadingTransport struct{ readErr chan error }

func (t *bodyReadingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	_, err := io.Copy(io.Discard, req.Body)
	t.readErr <- err

	return nil, errors.New("bodyReadingTransport: no response")
}

// TestUploadMediaDoesNotSendATruncatedBody pins what happens when the source
// reader fails partway through. Closing the pipe cleanly would end the body with
// a well-formed boundary, so Meta would accept a silently truncated upload,
// return a media id nobody reads, and leave a response nobody closes. The write
// end must fail the reader instead, so the request dies with it.
func TestUploadMediaDoesNotSendATruncatedBody(t *testing.T) {
	tr := &bodyReadingTransport{readErr: make(chan error, 1)}
	c := NewClient("token")
	c.HTTPClient = &http.Client{Transport: tr}

	if _, err := c.UploadMedia(context.Background(), "phone-id", "image/png", &faultyReader{left: 1 << 12}, 1<<20, "f.png"); err == nil {
		t.Fatal("expected an error when the source reader fails, got nil")
	}

	select {
	case err := <-tr.readErr:
		if err == nil {
			t.Error("the transport read the body to a clean EOF: a truncated upload was sent as if it were complete")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the transport never finished reading the request body")
	}
}

type trackedBody struct {
	io.Reader
	closed atomic.Bool
}

func (b *trackedBody) Close() error {
	b.closed.Store(true)

	return nil
}

// answeringTransport replies before the body is written, the way a server does
// when it rejects the request outright.
type answeringTransport struct{ body *trackedBody }

func (t *answeringTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// The RoundTripper contract requires closing the request body; doing it is
	// also what fails the multipart write upstream.
	if req.Body != nil {
		_ = req.Body.Close()
	}

	return &http.Response{StatusCode: http.StatusBadRequest, Body: t.body, Header: make(http.Header)}, nil
}

// TestUploadMediaClosesResponseOnWriteFailure covers the path where the request
// succeeds and the body write does not. net/http prefers a received response
// over a request-body error, so UploadMedia returns the write error while a live
// response sits behind it -- and every one of those paths returns before the
// <-done that owns it. Unclaimed, the connection and its read loop stay pinned
// until the peer gives up.
func TestUploadMediaClosesResponseOnWriteFailure(t *testing.T) {
	body := &trackedBody{Reader: strings.NewReader(
		`{"error":{"code":100,"message":"Invalid phone id","type":"OAuthException"}}`)}
	c := NewClient("token")
	c.HTTPClient = &http.Client{Transport: &answeringTransport{body: body}}

	_, err := c.UploadMedia(context.Background(), "phone-id", "image/png",
		strings.NewReader(strings.Repeat("x", 1<<20)), 1<<20, "f.png")
	if err == nil {
		t.Fatal("expected an error when the server answers before the body is written, got nil")
	}

	// The server said why, and permanently. Returning the local pipe error
	// instead reads as transient I/O, so the caller retries a request that can
	// never succeed.
	if !strings.Contains(err.Error(), "Invalid phone id") {
		t.Errorf("expected the server's reason, got %v", err)
	}
	// And the local diagnostic must survive alongside it: which stage failed is
	// the only thing that says what went wrong on this side.
	if !strings.Contains(err.Error(), "create form file") {
		t.Errorf("the write-side error was discarded: %v", err)
	}
	// The Graph error has to stay reachable through the wrapping, or every
	// caller switching on it breaks.
	if _, ok := AsGraphError(err); !ok {
		t.Errorf("AsGraphError could not find the Graph error in %v", err)
	}
	if c.LastGraphError() == nil {
		t.Error("LastGraphError was left nil, so the caller cannot inspect the Graph error at all")
	}

	// The response is reclaimed off the caller's path, so give it a moment.
	for range 40 {
		if body.closed.Load() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	t.Error("the response body was never closed: the connection stays pinned until the peer gives up")
}

// shortReader ends early the way a truncated object does: fewer bytes than
// promised, then a clean io.EOF rather than an error.
type shortReader struct{ left int }

func (r *shortReader) Read(p []byte) (int, error) {
	if r.left <= 0 {
		return 0, io.EOF
	}

	n := min(len(p), r.left)
	r.left -= n

	return n, nil
}

// TestUploadMediaRejectsAShortRead covers the truncation that reports success.
// io.Copy returns nil on a reader that simply ends, so without checking the
// count against fsize the body closes cleanly and Meta stores a partial file
// under an id the caller has no reason to distrust.
func TestUploadMediaRejectsAShortRead(t *testing.T) {
	tr := &bodyReadingTransport{readErr: make(chan error, 1)}
	c := NewClient("token")
	c.HTTPClient = &http.Client{Transport: tr}

	_, err := c.UploadMedia(context.Background(), "phone-id", "image/png", &shortReader{left: 10}, 1<<20, "f.png")
	if err == nil {
		t.Fatal("a reader that delivered 10 of 1048576 bytes was accepted as a complete upload")
	}
	if !strings.Contains(err.Error(), "read 10 of 1048576 bytes") {
		t.Errorf("expected the short read to be named, got %v", err)
	}
}

// TestUploadWriteErrKeepsBothCauses pins the half of the error contract that
// says a cancellation must not erase the failure underneath it. A source stream
// that genuinely died while the context happened to be dead is not retryable,
// and a caller that sees only context.Canceled will retry it forever.
func TestUploadWriteErrKeepsBothCauses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := uploadWriteErr(ctx, "copy file", errFaulty)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled to survive, got %v", err)
	}
	if !errors.Is(err, errFaulty) {
		t.Errorf("expected the underlying cause to survive, got %v", err)
	}
}

// lateClosingTransport reads part of the body, then closes it and answers --
// so the terminating boundary that mw.Close writes is what fails.
type lateClosingTransport struct {
	body    *trackedBody
	reached chan struct{}
}

func (t *lateClosingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Let every field write complete, then interrupt the terminating boundary,
	// so the call that fails is mw.Close -- not the file copy or the field
	// write before it. The pipe is unbuffered, so once the field's bytes have
	// all been read its Write has returned; the next Read blocks until mw.Close
	// starts writing the boundary, and closing then fails that write.
	var seen []byte
	buf := make([]byte, 256)
	for !bytes.Contains(seen, []byte("whatsapp")) {
		n, err := req.Body.Read(buf)
		seen = append(seen, buf[:n]...)
		if err != nil {
			break
		}
	}
	one := make([]byte, 1)
	_, _ = req.Body.Read(one)
	_ = req.Body.Close()
	close(t.reached)

	return &http.Response{StatusCode: http.StatusBadRequest, Body: t.body, Header: make(http.Header)}, nil
}

// TestUploadMediaCleansUpWhenClosingFails covers the exit path where the body
// was written fine and closing it is what fails. Marking the upload finished
// before those closes skipped the whole cleanup on exactly these paths: the
// pipe stayed open, the hand-off channel was never drained, and a live response
// was left with no owner.
func TestUploadMediaCleansUpWhenClosingFails(t *testing.T) {
	body := &trackedBody{Reader: strings.NewReader(`{"error":"bad request"}`)}
	tr := &lateClosingTransport{body: body, reached: make(chan struct{})}
	c := NewClient("token")
	c.HTTPClient = &http.Client{Transport: tr}

	before := runtime.NumGoroutine()

	if _, err := c.UploadMedia(context.Background(), "phone-id", "image/png",
		strings.NewReader("small"), 5, "f.png"); err == nil {
		t.Fatal("expected an error when the body cannot be closed, got nil")
	}

	select {
	case <-tr.reached:
	case <-time.After(5 * time.Second):
		t.Fatal("the transport never ran")
	}

	for range 40 {
		if body.closed.Load() && runtime.NumGoroutine() <= before {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	t.Errorf("cleanup was skipped: response closed=%v, goroutines before=%d now=%d",
		body.closed.Load(), before, runtime.NumGoroutine())
}

// TestUploadMediaTruncationReachesTheServer is the ordering test the fake
// transports cannot do. Framing is what makes the abort visible: against a real
// server the close order decides whether a failed upload arrives as a complete,
// parseable multipart body -- which the server accepts at a clean EOF and has no
// way to recognise as truncated -- or as a stream that ends mid-part and is
// rejected. Reading the body raw, as a stub transport does, cannot tell the two
// apart.
func TestUploadMediaTruncationReachesTheServer(t *testing.T) {
	type outcome struct {
		parts    int
		copyErr  error
		parseErr error
	}
	results := make(chan outcome, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var got outcome

		mr, err := r.MultipartReader()
		if err != nil {
			got.parseErr = err
			results <- got
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		for {
			part, err := mr.NextPart()
			if err != nil {
				if err != io.EOF {
					got.parseErr = err
				}

				break
			}
			if _, err := io.Copy(io.Discard, part); err != nil {
				got.copyErr = err

				break
			}
			got.parts++
		}

		results <- got
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	c := NewClient("token")
	c.GraphAPIVersion = "v1.0"
	// The package already has this helper; it keeps the real transport, and
	// therefore the real chunked framing, underneath.
	c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

	// Dies partway through a file it claims is much larger.
	if _, err := c.UploadMedia(context.Background(), "phone-id", "image/png",
		&faultyReader{left: 4096}, 1<<20, "f.png"); err == nil {
		t.Fatal("expected an error when the source reader fails, got nil")
	}

	select {
	case got := <-results:
		if got.copyErr == nil && got.parseErr == nil {
			t.Errorf("the server parsed %d complete part(s) and saw a clean EOF: a truncated upload arrived looking whole", got.parts)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the server never handled the request")
	}
}

var errTransportDied = errors.New("connection reset by peer")

// dyingTransport takes some of the body and then loses the connection, the way
// a real one does mid-upload.
type dyingTransport struct{}

func (dyingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	buf := make([]byte, 64)
	_, _ = req.Body.Read(buf)
	_ = req.Body.Close()

	return nil, errTransportDied
}

// TestUploadMediaKeepsTheTransportError covers the case where the request dies
// on the wire. The caller is left holding the pipe error that closing the body
// produced locally, which reads as transient I/O -- so the real cause has to
// travel with it, or a dead connection is indistinguishable from a hiccup.
func TestUploadMediaKeepsTheTransportError(t *testing.T) {
	c := NewClient("token")
	c.HTTPClient = &http.Client{Transport: dyingTransport{}}

	_, err := c.UploadMedia(context.Background(), "phone-id", "image/png",
		strings.NewReader(strings.Repeat("x", 1<<20)), 1<<20, "f.png")
	if err == nil {
		t.Fatal("expected an error when the transport dies, got nil")
	}
	if !errors.Is(err, errTransportDied) {
		t.Errorf("the transport failure never reached the caller: %v", err)
	}
}

// TestUploadMediaClearsStaleGraphError pins that a failed upload's Graph error
// does not outlive it. The cleanup defer records one, so without a reset on
// entry LastGraphError would still describe an earlier call.
func TestUploadMediaClearsStaleGraphError(t *testing.T) {
	body := &trackedBody{Reader: strings.NewReader(
		`{"error":{"code":100,"message":"Invalid phone id","type":"OAuthException"}}`)}
	c := NewClient("token")
	c.HTTPClient = &http.Client{Transport: &answeringTransport{body: body}}

	if _, err := c.UploadMedia(context.Background(), "phone-id", "image/png",
		strings.NewReader(strings.Repeat("x", 1<<20)), 1<<20, "f.png"); err == nil {
		t.Fatal("expected the first upload to fail")
	}
	if c.LastGraphError() == nil {
		t.Fatal("expected the first failure to record a Graph error")
	}

	// A second call that never reaches the server must not inherit it.
	c.HTTPClient = &http.Client{Transport: dyingTransport{}}
	if _, err := c.UploadMedia(context.Background(), "phone-id", "image/png",
		strings.NewReader(strings.Repeat("x", 1<<20)), 1<<20, "f.png"); err == nil {
		t.Fatal("expected the second upload to fail")
	}
	if c.LastGraphError() != nil {
		t.Errorf("LastGraphError still holds the previous call's error: %v", c.LastGraphError())
	}
}

// TestUploadMediaSucceeds is the happy path, and it exists because everything
// else here tests a failure. Without it the whole success sequence -- both
// parts written, the terminating boundary sent, the id decoded -- is invisible:
// swapping the two closes at the end, dropping the messaging_product field, or
// returning an empty id all break every real upload while the suite stays green.
func TestUploadMediaSucceeds(t *testing.T) {
	const payload = "the actual image bytes"

	type received struct {
		method    string
		path      string
		auth      string
		file      string
		fileField string
		product   string
		parseErr  error
	}
	got := make(chan received, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rec received
		rec.method, rec.path, rec.auth = r.Method, r.URL.Path, r.Header.Get("Authorization")

		mr, err := r.MultipartReader()
		if err != nil {
			rec.parseErr = err
			got <- rec
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		for {
			part, err := mr.NextPart()
			if err != nil {
				if err != io.EOF {
					rec.parseErr = err
				}

				break
			}

			body, err := io.ReadAll(part)
			if err != nil {
				rec.parseErr = err

				break
			}

			switch part.FormName() {
			case "messaging_product":
				rec.product = string(body)
			default:
				rec.fileField = part.FormName()
				rec.file = string(body)
			}
		}

		got <- rec
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"MEDIA-1"}`))
	}))
	defer srv.Close()

	c := NewClient("token")
	c.GraphAPIVersion = "v99.0"
	c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

	id, err := c.UploadMedia(context.Background(), "phone-id", "image/png",
		strings.NewReader(payload), int64(len(payload)), "f.png")
	if err != nil {
		t.Fatalf("a well-formed upload failed: %v", err)
	}
	if id != "MEDIA-1" {
		t.Errorf("expected the id from the response, got %q", id)
	}

	select {
	case rec := <-got:
		if rec.parseErr != nil {
			t.Errorf("the server could not parse the body: %v", rec.parseErr)
		}
		if rec.file != payload {
			t.Errorf("file part = %q, want %q", rec.file, payload)
		}
		if rec.fileField != "file" {
			t.Errorf("file part sent under field %q, want \"file\"", rec.fileField)
		}
		if rec.product != "whatsapp" {
			t.Errorf("messaging_product = %q, want \"whatsapp\"", rec.product)
		}
		// The request line and credentials matter as much as the body: an
		// unauthenticated PUT to a misspelled edge breaks every upload, and a
		// stub server that accepts anything will not notice.
		if rec.method != http.MethodPost {
			t.Errorf("method = %s, want POST", rec.method)
		}
		if rec.path != "/v99.0/phone-id/media" {
			t.Errorf("path = %q", rec.path)
		}
		if rec.auth != "Bearer token" {
			t.Errorf("Authorization = %q", rec.auth)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the server never handled the request")
	}
}

// TestUploadSessionHappyPaths covers the two resumable-upload methods, which
// this change gave a ctx and rebuilt their requests around -- with no test
// touching either one. Without this, renaming a header, flipping a method, or
// decoding the response into the wrong field breaks every template header
// upload in production while the suite stays green.
func TestUploadSessionHappyPaths(t *testing.T) {
	t.Run("NewUploadSession", func(t *testing.T) {
		type seen struct {
			method, path, auth, ctype string
			body                      map[string]any
			ctxLive                   bool
		}
		got := make(chan seen, 1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := seen{
				method: r.Method,
				path:   r.URL.Path,
				auth:   r.Header.Get("Authorization"),
				ctype:  r.Header.Get("Content-Type"),
			}
			_ = json.NewDecoder(r.Body).Decode(&s.body)
			got <- s
			_, _ = w.Write([]byte(`{"id":"upload:SESSION-1"}`))
		}))
		defer srv.Close()

		c := NewClient("token")
		c.GraphAPIVersion = "v99.0"
		c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

		id, err := c.NewUploadSession(context.Background(), "app-id", NewUploadSessionParams{
			FileLength: 1024, FileName: "header.png", FileType: "image/png",
		})
		if err != nil {
			t.Fatalf("a well-formed call failed: %v", err)
		}
		if id != "upload:SESSION-1" {
			t.Errorf("id = %q, want the id from the response", id)
		}

		s := <-got
		if s.method != http.MethodPost {
			t.Errorf("method = %s, want POST", s.method)
		}
		if s.path != "/v99.0/app-id/uploads" {
			t.Errorf("path = %q", s.path)
		}
		if s.auth != "Bearer token" {
			t.Errorf("Authorization = %q, want Bearer", s.auth)
		}
		if s.ctype != "application/json" {
			t.Errorf("Content-Type = %q", s.ctype)
		}
		// SessionType defaults to attachment; Graph rejects the call without it.
		if s.body["session_type"] != "attachment" {
			t.Errorf("session_type = %v, want attachment", s.body["session_type"])
		}
		if s.body["file_name"] != "header.png" {
			t.Errorf("file_name = %v", s.body["file_name"])
		}
	})

	t.Run("UploadHeaderHandle", func(t *testing.T) {
		const payload = "header bytes"

		type seen struct {
			method, path, auth, ctype, rng, offset, body string
		}
		got := make(chan seen, 1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := io.ReadAll(r.Body)
			got <- seen{
				method: r.Method, path: r.URL.Path,
				auth:  r.Header.Get("Authorization"),
				ctype: r.Header.Get("Content-Type"),
				rng:   r.Header.Get("Content-Range"),
				// canonicalised by net/http
				offset: r.Header.Get("File_offset"),
				body:   string(b),
			}
			_, _ = w.Write([]byte(`{"h":"HANDLE-1"}`))
		}))
		defer srv.Close()

		c := NewClient("token")
		c.GraphAPIVersion = "v99.0"
		c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

		h, err := c.UploadHeaderHandle(context.Background(), "upload:SESSION-1", strings.NewReader(payload))
		if err != nil {
			t.Fatalf("a well-formed call failed: %v", err)
		}
		if h != "HANDLE-1" {
			t.Errorf("handle = %q, want the h from the response", h)
		}

		s := <-got
		if s.method != http.MethodPost {
			t.Errorf("method = %s, want POST", s.method)
		}
		if s.path != "/v99.0/upload:SESSION-1" {
			t.Errorf("path = %q", s.path)
		}
		// OAuth, not Bearer: this endpoint is the one that differs.
		if s.auth != "OAuth token" {
			t.Errorf("Authorization = %q, want OAuth", s.auth)
		}
		if s.ctype != "application/octet-stream" {
			t.Errorf("Content-Type = %q", s.ctype)
		}
		if s.rng != "bytes 0-0/*" {
			t.Errorf("Content-Range = %q", s.rng)
		}
		if s.offset != "0" {
			t.Errorf("file_offset = %q", s.offset)
		}
		if s.body != payload {
			t.Errorf("body = %q, want %q", s.body, payload)
		}
	})
}

// TestSendAndMediaHappyPaths covers the working path of the remaining methods
// whose signatures this change touched. Every other test here drives them into
// a failure, which means the request they actually build was unasserted: the
// HTTP method, the URL, the Authorization header and the decoded response could
// all be broken without a single test noticing.
func TestSendAndMediaHappyPaths(t *testing.T) {
	type seen struct {
		method, path, auth, ctype, accept string
		body                              map[string]any
	}

	newStub := func(t *testing.T, respBody string) (*Client, chan seen) {
		t.Helper()

		got := make(chan seen, 1)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := seen{
				method: r.Method, path: r.URL.Path,
				auth:   r.Header.Get("Authorization"),
				ctype:  r.Header.Get("Content-Type"),
				accept: r.Header.Get("Accept"),
			}
			_ = json.NewDecoder(r.Body).Decode(&s.body)
			got <- s
			_, _ = w.Write([]byte(respBody))
		}))
		t.Cleanup(srv.Close)

		c := NewClient("token")
		c.GraphAPIVersion = "v99.0"
		c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

		return c, got
	}

	msg := &MessageObject{MessagingProduct: "whatsapp", To: "5511999999999", Type: "text"}

	t.Run("SendMessage", func(t *testing.T) {
		c, got := newStub(t, `{"messages":[{"id":"wamid.ABC"}]}`)

		res, err := c.SendMessage(context.Background(), "phone-id", msg)
		if err != nil {
			t.Fatalf("a well-formed send failed: %v", err)
		}
		if len(res.Messages) != 1 || res.Messages[0].ID != "wamid.ABC" {
			t.Errorf("response not decoded: %+v", res)
		}

		s := <-got
		if s.method != http.MethodPost {
			t.Errorf("method = %s, want POST", s.method)
		}
		if s.path != "/v99.0/phone-id/messages" {
			t.Errorf("path = %q", s.path)
		}
		if s.auth != "Bearer token" {
			t.Errorf("Authorization = %q", s.auth)
		}
		if s.ctype != "application/json" {
			t.Errorf("Content-Type = %q", s.ctype)
		}
		if s.body["to"] != "5511999999999" {
			t.Errorf("the message never reached the body: %v", s.body)
		}
	})

	t.Run("SendMessage routes a BSUID recipient", func(t *testing.T) {
		// The send methods rewrite a BSUID from `to` into `recipient` before
		// encoding. Tested in isolation elsewhere, but nothing pinned that the
		// send actually calls it -- so BSUID traffic could silently stop being
		// routed and go out with the wrong field.
		c, got := newStub(t, `{"messages":[{"id":"wamid.BSUID"}]}`)

		const bsuid = "US.13491208655302741918"
		if _, err := c.SendMessage(context.Background(), "phone-id",
			&MessageObject{MessagingProduct: "whatsapp", To: bsuid, Type: "text"}); err != nil {
			t.Fatalf("a well-formed send failed: %v", err)
		}

		s := <-got
		if s.body["recipient"] != bsuid {
			t.Errorf("recipient = %v, want the BSUID moved out of to", s.body["recipient"])
		}
		if _, present := s.body["to"]; present {
			t.Errorf("to must be omitted for a BSUID send, body was %v", s.body)
		}
	})

	t.Run("SendMarketingMessage", func(t *testing.T) {
		c, got := newStub(t, `{"messages":[{"id":"wamid.MKT"}]}`)

		res, err := c.SendMarketingMessage(context.Background(), "phone-id", msg)
		if err != nil {
			t.Fatalf("a well-formed send failed: %v", err)
		}
		if len(res.Messages) != 1 || res.Messages[0].ID != "wamid.MKT" {
			t.Errorf("response not decoded: %+v", res)
		}

		s := <-got
		// A different endpoint from SendMessage; swapping them silently sends
		// marketing traffic down the ordinary path.
		if s.path != "/v99.0/phone-id/marketing_messages" {
			t.Errorf("path = %q", s.path)
		}
		if s.method != http.MethodPost {
			t.Errorf("method = %s, want POST", s.method)
		}
	})

	t.Run("GetMedia", func(t *testing.T) {
		c, got := newStub(t, `{"url":"https://example/media","mime_type":"image/png","file_size":12,"id":"MID"}`)

		res, err := c.GetMedia(context.Background(), "media-id")
		if err != nil {
			t.Fatalf("a well-formed lookup failed: %v", err)
		}
		if res.URL != "https://example/media" || res.MimeType != "image/png" || res.ID != "MID" {
			t.Errorf("response not decoded: %+v", res)
		}

		s := <-got
		if s.method != http.MethodGet {
			t.Errorf("method = %s, want GET", s.method)
		}
		if s.path != "/v99.0/media-id" {
			t.Errorf("path = %q", s.path)
		}
		if s.auth != "Bearer token" {
			t.Errorf("Authorization = %q", s.auth)
		}
	})

	t.Run("DownloadMedia", func(t *testing.T) {
		const content = "raw image bytes"

		got := make(chan seen, 1)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got <- seen{
				method: r.Method, path: r.URL.Path,
				auth:   r.Header.Get("Authorization"),
				accept: r.Header.Get("Accept"),
			}
			_, _ = w.Write([]byte(content))
		}))
		defer srv.Close()

		c := NewClient("token")
		c.HTTPClient = &http.Client{Transport: http.DefaultTransport}

		var out bytes.Buffer
		n, err := c.DownloadMedia(context.Background(),
			&GetMediaResult{URL: srv.URL + "/blob", MimeType: "image/png"}, &out)
		if err != nil {
			t.Fatalf("a well-formed download failed: %v", err)
		}
		if int(n) != len(content) || out.String() != content {
			t.Errorf("wrote %d bytes %q, want %d %q", n, out.String(), len(content), content)
		}

		s := <-got
		if s.method != http.MethodGet {
			t.Errorf("method = %s, want GET", s.method)
		}
		if s.path != "/blob" {
			t.Errorf("path = %q: DownloadMedia must use the URL it was given", s.path)
		}
		if s.auth != "Bearer token" {
			t.Errorf("Authorization = %q", s.auth)
		}
		if s.accept != "image/png" {
			t.Errorf("Accept = %q, want the media's mime type", s.accept)
		}
	})
}

// TestAsGraphErrorWalksTheChain pins the contract 15 consumer call sites depend
// on. It used to be a bare type assertion, which made the error's position in
// the chain load-bearing: wrapping it anywhere -- as the upload cleanup path now
// does -- would silently stop every switch on it from matching.
func TestAsGraphErrorWalksTheChain(t *testing.T) {
	ge := &GraphError{Code: 100, Message: "Invalid phone id"}

	if _, ok := AsGraphError(ge); !ok {
		t.Error("a bare *GraphError was not recognised")
	}
	if _, ok := AsGraphError(fmt.Errorf("upload: %w", ge)); !ok {
		t.Error("a wrapped *GraphError was not recognised")
	}
	if _, ok := AsGraphError(errors.Join(errors.New("copy file: short read"), ge)); !ok {
		t.Error("a joined *GraphError was not recognised")
	}
	if _, ok := AsGraphError(errors.New("something else")); ok {
		t.Error("a plain error was reported as a Graph error")
	}
	if _, ok := AsGraphError(nil); ok {
		t.Error("nil was reported as a Graph error")
	}
}

// TestGraphErrorLifecycleIsPerCall pins that the error accessors describe the
// call the caller just made. errorFromResponse always refreshes the raw body but
// only sets lastGraphError when the payload carries an error code, so without a
// reset on entry a Graph error from one call can be read paired with the raw
// body of a later, entirely different failure -- a mismatch the caller has no
// way to detect.
func TestGraphErrorLifecycleIsPerCall(t *testing.T) {
	graph := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"code":100,"message":"Invalid phone id","type":"OAuthException"}}`))
	}))
	defer graph.Close()

	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`<html>gateway</html>`))
	}))
	defer gateway.Close()

	c := NewClient("token")
	msg := &MessageObject{MessagingProduct: "whatsapp", To: "5511999999999", Type: "text"}

	c.HTTPClient = &http.Client{Transport: rewriteHost{base: graph.URL, next: http.DefaultTransport}}
	if _, err := c.SendMessage(context.Background(), "phone-id", msg); err == nil {
		t.Fatal("expected the Graph error")
	}
	if c.LastGraphError() == nil {
		t.Fatal("expected the Graph error to be recorded")
	}

	// A gateway page carries no error code, so lastGraphError is not refreshed.
	c.HTTPClient = &http.Client{Transport: rewriteHost{base: gateway.URL, next: http.DefaultTransport}}
	if _, err := c.SendMessage(context.Background(), "phone-id", msg); err == nil {
		t.Fatal("expected the gateway error")
	}

	if ge := c.LastGraphError(); ge != nil {
		t.Errorf("LastGraphError still describes the previous call (%v) while LastErrorRawBody is %q",
			ge, c.LastErrorRawBody())
	}
}

// TestSuccessPathRejectsNon200 covers the status gate on the path where the body
// was written fine and the server still said no. Every other failure test here
// arrives through the cleanup defer instead, so without this a 400 decodes as an
// empty JSON object and the call returns ("", nil) -- an empty media id reported
// as success, which the caller happily sends on.
func TestSuccessPathRejectsNon200(t *testing.T) {
	newStub := func(t *testing.T) *Client {
		t.Helper()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.Copy(io.Discard, r.Body)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"code":100,"message":"Invalid phone id","type":"OAuthException"}}`))
		}))
		t.Cleanup(srv.Close)

		c := NewClient("token")
		c.GraphAPIVersion = "v99.0"
		c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

		return c
	}

	t.Run("UploadMedia", func(t *testing.T) {
		const payload = "the actual image bytes"

		id, err := newStub(t).UploadMedia(context.Background(), "phone-id", "image/png",
			strings.NewReader(payload), int64(len(payload)), "f.png")
		if err == nil {
			t.Fatalf("a 400 was reported as success, returning id %q", id)
		}
		if id != "" {
			t.Errorf("id = %q, want empty on failure", id)
		}
		if _, ok := AsGraphError(err); !ok {
			t.Errorf("expected the Graph error, got %v", err)
		}
	})

	t.Run("GetMedia", func(t *testing.T) {
		res, err := newStub(t).GetMedia(context.Background(), "media-id")
		if err == nil {
			t.Fatalf("a 400 was reported as success, returning %+v", res)
		}
		if _, ok := AsGraphError(err); !ok {
			t.Errorf("expected the Graph error, got %v", err)
		}
	})
}

// TestUploadSessionRateLimitSentinel pins the translation a caller backs off on.
// Without it a code-4 Graph error surfaces as an ordinary failure and a client
// that waits on ErrApplicationRateLimitReached keeps hammering Meta straight
// through the limit. It also exercises AsGraphError at these two call sites,
// which is the only place the errors.As rewrite is load-bearing.
func TestUploadSessionRateLimitSentinel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"code":4,"message":"Application request limit reached"}}`))
	}))
	defer srv.Close()

	c := NewClient("token")
	c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

	if _, err := c.NewUploadSession(context.Background(), "app-id", NewUploadSessionParams{FileLength: 1}); !errors.Is(err, ErrApplicationRateLimitReached) {
		t.Errorf("NewUploadSession: expected ErrApplicationRateLimitReached, got %v", err)
	}
	if _, err := c.UploadHeaderHandle(context.Background(), "upload:S", strings.NewReader("x")); !errors.Is(err, ErrApplicationRateLimitReached) {
		t.Errorf("UploadHeaderHandle: expected ErrApplicationRateLimitReached, got %v", err)
	}
}

// TestErrorAccessorsCarryTheResponse pins what the accessors actually hold, not
// just when they are cleared. LastErrorRawBody is the fallback a caller reads
// when Meta answers something that is not a Graph error at all -- a gateway page,
// an HTML block -- and nothing asserted it was ever populated. HTTPStatusCode is
// stamped onto the Graph error and consumers classify retryable versus permanent
// on it: without the stamp every 5xx reads as permanent and the work is dropped
// instead of retried.
func TestErrorAccessorsCarryTheResponse(t *testing.T) {
	t.Run("graph error carries the status and the raw body", func(t *testing.T) {
		const raw = `{"error":{"code":131026,"message":"Message undeliverable","type":"OAuthException"}}`

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(raw))
		}))
		defer srv.Close()

		c := NewClient("token")
		c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

		_, err := c.SendMessage(context.Background(), "phone-id",
			&MessageObject{MessagingProduct: "whatsapp", To: "5511999999999", Type: "text"})
		if err == nil {
			t.Fatal("expected the Graph error")
		}

		ge, ok := AsGraphError(err)
		if !ok {
			t.Fatalf("expected a *GraphError, got %v", err)
		}
		if ge.Code != 131026 {
			t.Errorf("Code = %d, want 131026", ge.Code)
		}
		// Consumers branch on this to decide retry versus discard.
		if ge.HTTPStatusCode != http.StatusServiceUnavailable {
			t.Errorf("HTTPStatusCode = %d, want 503: a 5xx that reports 0 classifies as permanent and the work is dropped", ge.HTTPStatusCode)
		}
		if c.LastErrorRawBody() != raw {
			t.Errorf("LastErrorRawBody = %q, want the response body verbatim", c.LastErrorRawBody())
		}
	})

	t.Run("non-graph body still reaches the caller", func(t *testing.T) {
		const raw = `<html>bad gateway</html>`

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(raw))
		}))
		defer srv.Close()

		c := NewClient("token")
		c.HTTPClient = &http.Client{Transport: rewriteHost{base: srv.URL, next: http.DefaultTransport}}

		if _, err := c.GetMedia(context.Background(), "media-id"); err == nil {
			t.Fatal("expected the gateway error")
		}

		// No error code, so no Graph error -- the raw body is the only thing
		// that says what happened.
		if c.LastGraphError() != nil {
			t.Errorf("LastGraphError = %v, want nil for a body with no error code", c.LastGraphError())
		}
		if c.LastErrorRawBody() != raw {
			t.Errorf("LastErrorRawBody = %q, want the gateway page", c.LastErrorRawBody())
		}
	})
}
