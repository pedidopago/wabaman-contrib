package fbgraph

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// PostSubscribedApps is a required step for Embedded SignUp
func (c *Client) PostSubscribedApps(ctx context.Context, wabaID string) error {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	nilContent := strings.NewReader("{}")

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/subscribed_apps", apiVersion, wabaID)

	req, err := NewRequestWithContext(ctx, http.MethodPost, url, nilContent)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		gerr := c.errorFromResponse(resp)
		if ge, ok := AsGraphError(gerr); ok {
			if ge.Code == 4 {
				return ErrApplicationRateLimitReached
			}
		}
		return gerr
	}

	io.Copy(io.Discard, resp.Body)
	return nil
}
