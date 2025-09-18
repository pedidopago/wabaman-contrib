package fbgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SMBSyncType string

const (
	SMBSyncTypeHistory         SMBSyncType = "history"
	SMBSyncTypeSMBAppStateSync SMBSyncType = "smb_app_state_sync"
)

type SMBAppDataResult struct {
	MessagingProduct string `json:"messaging_product"`
	RequestID        string `json:"request_id"`
	Success          bool   `json:"success"`
}

// This endpoint is used as part of messaging history synchronization process when onboarding
// business customers who have a WhatsApp Business app account and phone number.
func (c *Client) PostSMBAppData(ctx context.Context, whatsappID string, syncType SMBSyncType) (*SMBAppDataResult, error) {
	apiVersion := "v23.0"
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/smb_app_data", apiVersion, whatsappID)

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(map[string]any{
		"messaging_product": "whatsapp",
		"sync_type":         syncType,
	}); err != nil {
		return nil, fmt.Errorf("encode message: %w", err)
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

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.errorFromResponse(resp)
	}

	result := &SMBAppDataResult{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}
