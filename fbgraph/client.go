package fbgraph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
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

func (c *Client) errorFromResponse(resp *http.Response) error {
	herr := &GraphError{}
	jbdbuff := new(bytes.Buffer)
	io.Copy(jbdbuff, resp.Body)
	if err := json.Unmarshal(jbdbuff.Bytes(), herr); err != nil {
		return fmt.Errorf("http status: %d (%s); %w - %s", resp.StatusCode, resp.Status, err, jbdbuff.String())
	}
	herr.HTTPStatusCode = resp.StatusCode
	return herr
}
