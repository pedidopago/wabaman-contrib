package fbgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Conversions API (CAPI) for Business Messaging. Used to report Click-to-WhatsApp
// ad conversions (e.g. a lead) back to Meta so ad optimization sees the outcome.
// Docs: developers.facebook.com/docs/marketing-api/conversions-api/business-messaging

// ConversionUserData is the CAPI user_data block. For business messaging the
// ctwa_clid is the attribution key and is NOT hashed.
type ConversionUserData struct {
	CtwaClid                  string `json:"ctwa_clid"`
	WhatsAppBusinessAccountID string `json:"whatsapp_business_account_id"`
}

// ConversionEvent is a single CAPI event. ActionSource is "business_messaging"
// and MessagingChannel is "whatsapp" for CTWA conversions.
type ConversionEvent struct {
	EventName        string             `json:"event_name"`
	EventTime        int64              `json:"event_time"`
	ActionSource     string             `json:"action_source"`
	MessagingChannel string             `json:"messaging_channel"`
	UserData         ConversionUserData `json:"user_data"`
}

func (c *Client) graphVersion() string {
	if c.GraphAPIVersion != "" {
		return c.GraphAPIVersion
	}
	return DefaultGraphAPIVersion
}

// CreateDataset provisions a Conversions API dataset on a WhatsApp Business
// Account and returns its id. POST /{WABA_ID}/dataset.
func (c *Client) CreateDataset(ctx context.Context, wabaID string) (datasetID string, err error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/dataset", c.graphVersion(), wabaID)

	req, err := NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", c.errorFromResponse(resp)
	}

	result := struct {
		ID string `json:"id"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.ID, nil
}

// GetDatasetID returns the id of the WABA's existing Conversions API dataset, or
// an empty string when none exists. GET /{WABA_ID}/dataset. Used to reuse a
// dataset the client provisioned themselves before creating a new one.
func (c *Client) GetDatasetID(ctx context.Context, wabaID string) (datasetID string, err error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/dataset", c.graphVersion(), wabaID)

	req, err := NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", c.errorFromResponse(resp)
	}

	// The edge returns a paged list; the WABA has at most one CTWA dataset today.
	result := struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Data) == 0 {
		return "", nil
	}
	return result.Data[0].ID, nil
}

// SendConversionEvents posts CAPI events to a dataset and returns how many Meta
// accepted (events_received). POST /{DATASET_ID}/events. Meta rejects the whole
// batch if any event_time is older than 7 days, so callers should send small
// batches and guard the window before calling.
func (c *Client) SendConversionEvents(ctx context.Context, datasetID string, events []ConversionEvent) (received int, err error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/events", c.graphVersion(), datasetID)

	body := struct {
		Data []ConversionEvent `json:"data"`
	}{Data: events}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return 0, fmt.Errorf("encode: %w", err)
	}

	req, err := NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return 0, fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, c.errorFromResponse(resp)
	}

	result := struct {
		EventsReceived int `json:"events_received"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decode response: %w", err)
	}

	return result.EventsReceived, nil
}
