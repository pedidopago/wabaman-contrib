package fbgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Conversation Routing: when more than one responder (partner apps, Meta
// Business Agent) serves a phone number, each thread has one owner.
// https://developers.facebook.com/documentation/business-messaging/whatsapp/conversation-routing/overview

// ThreadControlAction is what POST /{phone-number-id}/thread_control does.
type ThreadControlAction string

const (
	// ThreadControlPass hands the thread to another role (the escalation
	// partner when no target is given). Caller must own the thread.
	ThreadControlPass ThreadControlAction = "pass"
	// ThreadControlRelease returns the thread to idle. Caller must own it.
	ThreadControlRelease ThreadControlAction = "release"
	// ThreadControlTake moves the thread to the escalation partner. Only the
	// escalation partner may call it.
	ThreadControlTake ThreadControlAction = "take"
)

// RoutingRole identifies a responder's role on a phone number.
type RoutingRole string

const (
	RoutingRoleAIAgent         RoutingRole = "ai_agent"
	RoutingRoleCustomerService RoutingRole = "customer_service"
	RoutingRoleMarketing       RoutingRole = "marketing"
	RoutingRoleUtility         RoutingRole = "utility"
	RoutingRoleCTWA            RoutingRole = "ctwa"
	RoutingRoleEscalation      RoutingRole = "escalation"
	// RoutingRoleCalls is reported by routing_roles but never owns a thread.
	RoutingRoleCalls RoutingRole = "calls"
)

// ThreadControlRequest is the body of POST /{phone-number-id}/thread_control.
// Set exactly one of To (phone number / wa_id) or Recipient (BSUID).
type ThreadControlRequest struct {
	MessagingProduct string              `json:"messaging_product"`
	Action           ThreadControlAction `json:"action"`
	To               string              `json:"to,omitempty"`
	Recipient        string              `json:"recipient,omitempty"`
	// Metadata is forwarded verbatim to the receiving app (max 2000 chars).
	Metadata    string                   `json:"metadata,omitempty"`
	ControlPass *ThreadControlPassTarget `json:"control_pass,omitempty"`
}

// ThreadControlPassTarget targets a pass at a role instead of the escalation partner.
type ThreadControlPassTarget struct {
	TargetRole RoutingRole `json:"target_role"`
}

// ThreadControlResponse is Meta's answer to a successful thread_control call.
type ThreadControlResponse struct {
	MessagingProduct string `json:"messaging_product"`
	RequestID        string `json:"request_id"`
}

// ThreadControl passes, releases or takes a thread. A Graph error comes back
// as *GraphError (see AsGraphError).
func (c *Client) ThreadControl(ctx context.Context, phoneNumberID string, r ThreadControlRequest) (*ThreadControlResponse, error) {
	if (r.To == "") == (r.Recipient == "") {
		return nil, fmt.Errorf("thread control: set exactly one of to or recipient")
	}
	r.MessagingProduct = "whatsapp"
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("encode thread control: %w", err)
	}
	var out ThreadControlResponse
	if err := c.routingDo(ctx, http.MethodPost, phoneNumberID+"/thread_control", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RoutingRoles is GET /{phone-number-id}/routing_roles: the roles our
// business holds on the number. Roles is empty (not an error) when routing is
// on but we hold none; routing being off is a Graph error (code 100).
type RoutingRoles struct {
	Roles     []RoutingRole `json:"roles"`
	IsStandby bool          `json:"is_standby"`
}

// GetRoutingRoles reads our routing roles on a phone number. Rate limited per
// number: do not call it before every send.
func (c *Client) GetRoutingRoles(ctx context.Context, phoneNumberID string) (*RoutingRoles, error) {
	var out RoutingRoles
	if err := c.routingDo(ctx, http.MethodGet, phoneNumberID+"/routing_roles", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) routingDo(ctx context.Context, method, path string, body []byte, out any) error {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}
	var rd io.Reader
	if body != nil {
		rd = bytes.NewReader(body)
	}
	req, err := NewRequestWithContext(ctx, method, fmt.Sprintf("https://graph.facebook.com/%s/%s", apiVersion, path), rd)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return c.errorFromResponse(resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
