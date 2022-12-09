package fbgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type TokenType string

const (
	TokenTypeAPP  TokenType = "APP"
	TokenTypeUser TokenType = "USER"
)

type GranularScope struct {
	Scope     string   `json:"scope"`
	TargetIDs []string `json:"target_ids"`
}

type TokenInfo struct {
	AppID               string          `json:"app_id"`
	Type                TokenType       `json:"type"`
	Application         string          `json:"application"`
	DataAccessExpiresAt int64           `json:"data_access_expires_at,omitempty"` // Unix TS
	ExpiresAt           int64           `json:"expires_at,omitempty"`             // Unix TS
	IsValid             bool            `json:"is_valid"`
	IssuedAt            int64           `json:"issued_at,omitempty"` // Unix TS
	Scopes              []string        `json:"scopes,omitempty"`
	GranularScopes      []GranularScope `json:"granular_scopes,omitempty"`
	UserID              string          `json:"user_id,omitempty"`
}

func (c *Client) DebugToken(ctx context.Context, inputToken, accessToken string) (TokenInfo, error) {

	emptyd := TokenInfo{}

	url := fmt.Sprintf("https://graph.facebook.com/v15.0/debug_token?input_token=%s&access_token=%s", inputToken, accessToken)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return emptyd, fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return emptyd, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return emptyd, c.errorFromResponse(resp)
	}
	result := struct {
		Data TokenInfo `json:"data"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return emptyd, fmt.Errorf("decode response: %w", err)
	}
	return result.Data, nil
}
