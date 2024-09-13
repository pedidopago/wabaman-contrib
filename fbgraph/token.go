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

func (c *Client) DebugToken(ctx context.Context, inputToken string) (TokenInfo, error) {

	emptyd := TokenInfo{}

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/debug_token?input_token=%s&access_token=%s", apiVersion, inputToken, c.AccessToken)

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

func (c *Client) NewPermanentAccessToken(ctx context.Context, appID, appSecret, tempToken string) (string, error) {

	url := fmt.Sprintf("https://graph.facebook.com/v13.0/oauth/access_token?grant_type=fb_exchange_token&client_id=%s&client_secret=%s&fb_exchange_token=%s", appID, appSecret, tempToken)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", c.errorFromResponse(resp)
	}
	result := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	return result.AccessToken, nil
}
