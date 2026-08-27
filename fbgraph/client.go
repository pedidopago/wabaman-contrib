package fbgraph

import (
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	DebugTrace             bool
	DefaultGraphAPIVersion = "v23.0"
)

var DefaultHTTPClient = &http.Client{
	Timeout: time.Second * 120,
}

// Client is not safe for concurrent use. Each goroutine must use its own
// Client instance; shared instances require external synchronization.
type Client struct {
	HTTPClient      *http.Client
	AccessToken     string
	GraphAPIVersion string // use this to override the default API version (check DefaultGraphAPIVersion)

	lastGraphError   *GraphError
	lastErrorRawBody string
}

func NewClient(accessToken string) *Client {
	return &Client{
		HTTPClient:  DefaultHTTPClient,
		AccessToken: accessToken,
	}
}

func (c *Client) LastGraphError() *GraphError {
	return c.lastGraphError
}

func (c *Client) LastErrorRawBody() string {
	return c.lastErrorRawBody
}

// SendMessageFn is a function type for sending messages.
// (*Client).SendMessage is the default implementation of this function.
// (*Client).SendMarketingMessage is a specialized implementation for marketing messages.
type SendMessageFn func(ctx context.Context, phoneID string, msg *MessageObject) (*MessageObjectResult, error)

// SendMessage is not idempotent. A context cancelled between the write and the
// response returns an error for a message Meta may already have delivered, so
// context.Canceled here does not mean "not sent" and must not trigger a resend.
func (c *Client) SendMessage(ctx context.Context, phoneID string, msg *MessageObject) (*MessageObjectResult, error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}
	msg.routeBSUIDRecipient()

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", apiVersion, phoneID)
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(msg); err != nil {
		return nil, fmt.Errorf("encode message: %w", err)
	}
	if DebugTrace {
		println("fbgraph SendMessage", url, "\n", buf.String())
	}
	req, err := NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, c.errorFromResponse(resp)
	}
	result := &MessageObjectResult{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// SendMarketingMessage uses the Marketing Messages API to send a marketing message to a phone ID.
// You need to have access to the Marketing Messages API and have the appropriate permissions.
//
// See https://developers.facebook.com/documentation/business-messaging/whatsapp/marketing-messages/overview
func (c *Client) SendMarketingMessage(ctx context.Context, phoneID string, msg *MessageObject) (*MessageObjectResult, error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}
	msg.routeBSUIDRecipient()

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/marketing_messages", apiVersion, phoneID)
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(msg); err != nil {
		return nil, fmt.Errorf("encode message: %w", err)
	}
	if DebugTrace {
		println("fbgraph SendMarketingMessage", url, "\n", buf.String())
	}
	req, err := NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, c.errorFromResponse(resp)
	}
	result := &MessageObjectResult{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// UploadMedia streams r to Meta as a multipart upload and returns the media id.
//
// fsize must be the exact length of r. It is checked against the bytes actually
// copied, because a reader that simply ends reports io.EOF rather than an error
// -- so without it a short read is indistinguishable from a complete upload and
// Meta stores a truncated file under an id the caller has no reason to distrust.
// Pass 0 only when the length is genuinely unknown, which disables the check.
//
// ctx bounds the HTTP request, not r. A source reader that blocks forever
// blocks this call: cancellation has nothing to interrupt inside r.Read.
func (c *Client) UploadMedia(ctx context.Context, phoneID string, mimeType string, r io.Reader, fsize int64, filename string) (id string, err error) {
	// The cleanup defer below can set these, so a failed upload would otherwise
	// leave its Graph error readable after a later successful one.
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/media", apiVersion, phoneID)

	// Built before the goroutine on purpose. It does not read the body, and
	// doing it here means an unusable context -- a nil one, say -- fails while
	// nothing is writing to the pipe yet. Built inside the goroutine instead,
	// that failure leaves the multipart writer below parked on a pipe no one
	// will ever read, and UploadMedia never returns at all.
	req, err := NewRequestWithContext(ctx, http.MethodPost, url, pr)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.ContentLength = -1
	//TODO: calculate content length like in https://gist.github.com/cryptix/9dd094008b6236f4fc57
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	// do the request concurrently
	var resp *http.Response
	// Every exit path below now drains this -- the success path reads it, and
	// the cleanup defer reads it for the rest -- so the buffer is not load
	// bearing today, and no test can tell it apart from an unbuffered channel.
	// It stays as cheap insurance: this send is inside a goroutine, so an early
	// return added later that skips the drain would park it for the life of the
	// process, which is the exact defect this function already had once.
	done := make(chan error, 1)
	go func() {
		var rerr error
		resp, rerr = c.HTTPClient.Do(req)
		if rerr != nil {
			done <- fmt.Errorf("request failed: %w", rerr)

			return
		}
		done <- nil
	}()

	allok := false
	defer func() {
		if !allok {
			// Poison the pipe before closing the writer. mw.Close writes the
			// terminating boundary, so the other order sends a complete,
			// parseable multipart body that the server accepts at a clean EOF
			// -- a truncated upload it has no way to recognise as truncated.
			// This ordering is the whole truncation defense.
			// cmp.Or, because CloseWithError(nil) degrades to a plain Close --
			// a clean EOF, which is the truncation this ordering exists to
			// prevent. err is non-nil on every ordinary arrival here, but a
			// panic between the goroutine starting and allok reaches this with
			// nothing set.
			_ = pw.CloseWithError(cmp.Or(err, errUploadAborted))
			_ = mw.Close()

			// The request can still have succeeded: net/http prefers a received
			// response over a request-body write error, so a server that
			// answers early -- a 400 on a bad phone id, an oversized file --
			// leaves a response here that none of these paths reach the <-done
			// to own. Unclaimed, its connection and read loop stay pinned until
			// the peer gives up.
			switch rerr := <-done; {
			case rerr != nil:
				// The request itself failed, and the pipe error the caller is
				// holding is only the local symptom of it. Keep both: a dead
				// connection reads as transient local I/O otherwise, exactly
				// the mislabel this block exists to prevent.
				err = errors.Join(err, rerr)
			case resp != nil:
				// The server's own reason beats the pipe error that surfaced
				// locally: a 400 for a bad phone id is permanent, while
				// "read/write on closed pipe" reads as transient and gets
				// retried forever. This also populates LastGraphError.
				if resp.StatusCode != http.StatusOK {
					// Both causes, not just the server's: the write-side error
					// is the only thing that says which stage failed and, for a
					// short read, how many bytes actually arrived. Safe to wrap
					// now that AsGraphError walks the chain. The nesting is for
					// readability, not line count -- GraphError.Error() is
					// multi-line, so nothing here is single-line either way.
					err = fmt.Errorf("%w (graph: %w)", err, c.errorFromResponse(resp))
				}
				_ = resp.Body.Close()
			}
		}
	}()

	fh := make(textproto.MIMEHeader)
	fh.Set("Content-Type", mimeType)
	fh.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("file"), escapeQuotes(filename)))
	fpart, err := mw.CreatePart(fh)
	if err != nil {
		return "", uploadWriteErr(ctx, "create form file", err)
	}
	written, err := io.Copy(fpart, r)
	if err != nil {
		return "", uploadWriteErr(ctx, "copy file", err)
	}
	// Short only: rejecting a longer stream would fail uploads whose recorded
	// size merely lags the bytes. A zero or negative fsize disables the check on
	// its own, since written is never negative.
	if written < fsize {
		return "", fmt.Errorf("copy file: read %d of %d bytes", written, fsize)
	}
	if err := mw.WriteField("messaging_product", "whatsapp"); err != nil {
		return "", uploadWriteErr(ctx, "write field", err)
	}
	if err := mw.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}
	// io.PipeWriter.Close never fails, so there is nothing to check here.
	_ = pw.Close()
	// Only now: both closes can fail, and setting this above them skipped the
	// cleanup on exactly those paths -- pipe left open, done never drained, a
	// live response left unclaimed.
	allok = true
	if err := <-done; err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", c.errorFromResponse(resp)
	}
	idstruct := struct {
		ID string `json:"id"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&idstruct); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	return idstruct.ID, nil
}

// errUploadAborted marks a body abandoned without an error of its own, so the
// pipe is still closed with a failure rather than a clean end.
var errUploadAborted = errors.New("upload aborted before completion")

// uploadWriteErr prefers the context error while writing the multipart body.
// When the request fails, net/http closes the pipe reader plainly, so the
// writer only ever sees io.ErrClosedPipe -- the cancellation that actually
// caused it would never reach the caller, who then cannot tell a dead context
// from a genuine media failure and drops the job instead of retrying it.
func uploadWriteErr(ctx context.Context, stage string, err error) error {
	// ctx is never nil here: NewRequestWithContext rejects that before the pipe
	// exists, so this is only ever reached with a usable context.
	if cerr := ctx.Err(); cerr != nil {
		// Both causes, not just the context: a source stream that genuinely
		// died while the context happened to be dead would otherwise read as a
		// clean cancellation, and the caller would retry media that can never
		// succeed. Formatted on one line so logs stay greppable.
		return fmt.Errorf("%s: %w (context: %w)", stage, err, cerr)
	}

	return fmt.Errorf("%s: %w", stage, err)
}

func (c *Client) GetMedia(ctx context.Context, mediaID string) (*GetMediaResult, error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s", apiVersion, mediaID)
	req, err := NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, c.errorFromResponse(resp)
	}
	result := &GetMediaResult{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

func (c *Client) DownloadMedia(ctx context.Context, mr *GetMediaResult, out io.Writer) (nwritten int64, err error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	req, err := NewRequestWithContext(ctx, http.MethodGet, mr.URL, nil)
	if err != nil {
		return 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Set("Accept", mr.MimeType)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	log.Debug().Interface("media", mr).Int("http_status_code", resp.StatusCode).Interface("response_headers", resp.Header).Msg("downloading media")
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		nwritten, err := io.Copy(out, resp.Body)
		if err != nil {
			return nwritten, err
		}

		return nwritten, nil
	}

	log.Warn().Interface("media", mr).Int("http_status_code", resp.StatusCode).Interface("response_headers", resp.Header).Msg("fbgraph download media failed")

	return 0, c.errorFromResponse(resp)
}

type GetMediaResult struct {
	MessagingProduct string  `json:"messaging_product"`
	URL              string  `json:"url"`
	MimeType         string  `json:"mime_type"`
	Sha256           string  `json:"sha256"`
	FileSize         float64 `json:"file_size"`
	ID               string  `json:"id"`
}

func (mr *GetMediaResult) VerifyChecksum(r io.Reader) bool {
	computedHex, err := hex.DecodeString(mr.Sha256)
	if err != nil {
		return false
	}

	hh := sha256.New()
	if _, err := io.Copy(hh, r); err != nil {
		return false
	}

	h256sum := hh.Sum(nil)

	return bytes.Equal(computedHex, h256sum)
}

type NewUploadSessionParams struct {
	// The file length in bytes
	FileLength int64 `json:"file_length"`
	// The name of the file to be uploaded
	FileName string `json:"file_name"`
	// The MIME type of the file to be uploaded
	FileType string `json:"file_type"`
	// The type of upload session that is being requested by the app
	//
	// default: attachment
	SessionType string `json:"session_type"`
}

func (c *Client) NewUploadSession(ctx context.Context, fbAppID string, params NewUploadSessionParams) (id string, err error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	if params.SessionType == "" {
		params.SessionType = "attachment"
	}

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/uploads", apiVersion, fbAppID)
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(params); err != nil {
		return "", fmt.Errorf("encode message: %w", err)
	}
	req, err := NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		gerr := c.errorFromResponse(resp)

		if ge, ok := AsGraphError(gerr); ok {
			if ge.Code == 4 {
				return "", ErrApplicationRateLimitReached
			}
		}

		return "", gerr
	}
	idstruct := struct {
		ID string `json:"id"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&idstruct); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	return idstruct.ID, nil
}

func (c *Client) UploadHeaderHandle(ctx context.Context, uploadSessionID string, r io.Reader) (h string, err error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s", apiVersion, uploadSessionID)
	req, err := NewRequestWithContext(ctx, http.MethodPost, url, r)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", c.AccessToken))
	req.Header.Set("Content-Range", "bytes 0-0/*")
	req.Header.Set("file_offset", "0")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		gerr := c.errorFromResponse(resp)

		if ge, ok := AsGraphError(gerr); ok {
			if ge.Code == 4 {
				return "", ErrApplicationRateLimitReached
			}
		}

		return "", gerr
	}
	hstruct := struct {
		H string `json:"h"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&hstruct); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return hstruct.H, nil
}

func (c *Client) errorFromResponse(resp *http.Response) error {
	if DebugTrace {
		fmt.Printf("HTTP STATUS %d\nHEADERS:\n", resp.StatusCode)
		for k, v := range resp.Header {
			fmt.Printf("  %s: %s\n", k, strings.Join(v, ", "))
		}
	}

	eparent := struct {
		Error GraphError `json:"error"`
	}{}
	jbdbuff := new(bytes.Buffer)
	_, _ = io.Copy(jbdbuff, resp.Body)

	c.lastErrorRawBody = jbdbuff.String()

	if DebugTrace {
		fmt.Printf("BODY:\n%s\n", jbdbuff.String())
	}

	if err := json.Unmarshal(jbdbuff.Bytes(), &eparent); err != nil {
		return fmt.Errorf("http status: %d (%s); %w - %s", resp.StatusCode, resp.Status, err, jbdbuff.String())
	}
	if eparent.Error.Code == 0 {
		return fmt.Errorf("http status: %d (%s); %s", resp.StatusCode, resp.Status, jbdbuff.String())
	}
	eparent.Error.HTTPStatusCode = resp.StatusCode
	c.lastGraphError = &eparent.Error
	return c.lastGraphError
}
