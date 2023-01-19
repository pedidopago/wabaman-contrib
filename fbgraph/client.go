package fbgraph

import (
	"bytes"
	"encoding/json"
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
	DebugTrace bool
)

var DefaultHTTPClient = &http.Client{
	Timeout: time.Second * 120,
}

type Client struct {
	HTTPClient  *http.Client
	AccessToken string
}

func NewClient(accessToken string) *Client {
	return &Client{
		HTTPClient:  DefaultHTTPClient,
		AccessToken: accessToken,
	}
}

func (c *Client) SendMessage(phoneID string, msg *MessageObject) (*MessageObjectResult, error) {
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}
	url := fmt.Sprintf("https://graph.facebook.com/v13.0/%s/messages", phoneID)
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(msg); err != nil {
		return nil, fmt.Errorf("encode message: %w", err)
	}
	if DebugTrace {
		println("fbgraph SendMessage", url, "\n", buf.String())
	}
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
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

func (c *Client) UploadMedia(phoneID string, mimeType string, r io.Reader, fsize int64, filename string) (id string, err error) {
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	url := fmt.Sprintf("https://graph.facebook.com/v14.0/%s/media", phoneID)

	// do the request concurrently
	var resp *http.Response
	done := make(chan error)
	go func() {
		req, err := http.NewRequest(http.MethodPost, url, pr)
		if err != nil {
			done <- fmt.Errorf("new request: %w", err)
			return
		}
		req.ContentLength = -1
		//TODO: calculate content length like in https://gist.github.com/cryptix/9dd094008b6236f4fc57
		req.Header.Set("Content-Type", mw.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			done <- fmt.Errorf("request failed: %w", err)
			return
		}
		done <- nil
	}()
	allok := false
	defer func() {
		if !allok {
			mw.Close()
			pw.Close()
		}
	}()

	fh := make(textproto.MIMEHeader)
	fh.Set("Content-Type", mimeType)
	fh.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("file"), escapeQuotes(filename)))
	fpart, err := mw.CreatePart(fh)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(fpart, r); err != nil {
		return "", fmt.Errorf("copy file: %w", err)
	}
	if err := mw.WriteField("messaging_product", "whatsapp"); err != nil {
		return "", fmt.Errorf("write field: %w", err)
	}
	allok = true
	if err := mw.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}
	if err := pw.Close(); err != nil {
		return "", fmt.Errorf("close pipe writer: %w", err)
	}
	if err := <-done; err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()
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

func (c *Client) GetMedia(mediaID string) (*GetMediaResult, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v14.0/%s", mediaID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, c.errorFromResponse(resp)
	}
	result := &GetMediaResult{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

func (c *Client) DownloadMedia(mr *GetMediaResult, out io.Writer) error {
	req, err := http.NewRequest(http.MethodGet, mr.URL, nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Set("Accept", mr.MimeType)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	log.Debug().Interface("media", mr).Int("http_status_code", resp.StatusCode).Interface("response_headers", resp.Header).Msg("downloading media")
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		nwritten, err := io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		mwritten := int64(mr.FileSize)
		log.Debug().Int64("nwritten", nwritten).Int64("mwritten", mwritten).Msg("downloaded media")
		if mwritten != nwritten {
			return fmt.Errorf("was expecting %d bytes, but received %d", mwritten, nwritten)
		}
		return nil
	}
	log.Warn().Interface("media", mr).Int("http_status_code", resp.StatusCode).Interface("response_headers", resp.Header).Msg("fbgraph download media failed")
	return c.errorFromResponse(resp)
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
	//FIXME: implement this
	//TODO: implement this
	return true
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

func (c *Client) NewUploadSession(ownerID string, params NewUploadSessionParams) (id string, err error) {
	if params.SessionType == "" {
		params.SessionType = "attachment"
	}
	url := fmt.Sprintf("https://graph.facebook.com/v14.0/%s/uploads", ownerID)
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(params); err != nil {
		return "", fmt.Errorf("encode message: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
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

func (c *Client) errorFromResponse(resp *http.Response) error {
	eparent := struct {
		Error GraphError `json:"error"`
	}{}
	jbdbuff := new(bytes.Buffer)
	io.Copy(jbdbuff, resp.Body)
	if err := json.Unmarshal(jbdbuff.Bytes(), &eparent); err != nil {
		return fmt.Errorf("http status: %d (%s); %w - %s", resp.StatusCode, resp.Status, err, jbdbuff.String())
	}
	if eparent.Error.Code == 0 {
		return fmt.Errorf("http status: %d (%s); %s", resp.StatusCode, resp.Status, jbdbuff.String())
	}
	eparent.Error.HTTPStatusCode = resp.StatusCode
	return &eparent.Error
}
