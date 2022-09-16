package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pedidopago/go-common/util"
	"github.com/pedidopago/wabaman-contrib/rest"
)

const (
	DefaultBaseURL = "https://api.first.v2.pedidopago.com.br/wabaman"
)

type Client struct {
	JWT     string
	BaseURL string
}

func (c *Client) NewMessage(ctx context.Context, req *rest.NewMessageRequest) (*rest.NewMessageResponse, error) {
	output := &rest.NewMessageResponse{}
	if err := c.post(ctx, "/api/v1/message", req, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *Client) urlPrefix() string {
	return util.Default(c.BaseURL, DefaultBaseURL)
}

func (c *Client) post(ctx context.Context, suffix string, input, output any) error {
	d, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshal input: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.urlPrefix()+suffix, bytes.NewReader(d))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	if c.JWT != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.JWT))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.errorFromResponse(ctx, resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func (c *Client) errorFromResponse(ctx context.Context, resp *http.Response) error {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return fmt.Errorf("copy response body: %w", err)
	}
	jerr := &rest.ErrorResponse{}
	if err := json.Unmarshal(buf.Bytes(), jerr); err != nil {
		// not a valid json response!
		return &rest.ErrorResponse{
			StatusCode: resp.StatusCode,
			Raw:        buf.String(),
		}
	}
	jerr.StatusCode = resp.StatusCode
	return jerr
}
