package fbgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Meta's migration status values, returned by the Currency Migration API.
const (
	MigrationStatusInitiated        = "INITIATED"
	MigrationStatusAccepted         = "ACCEPTED"
	MigrationStatusInProgress       = "IN_PROGRESS"
	MigrationStatusReadyToComplete  = "READY_TO_COMPLETE"
	MigrationStatusCompleted        = "COMPLETED"
	MigrationStatusFailed           = "FAILED"
	MigrationStatusFailedToComplete = "FAILED_TO_COMPLETE"
	MigrationStatusRejected         = "REJECTED"
)

// MigrationIntentRequest starts a WABA currency migration.
type MigrationIntentRequest struct {
	// Currency is the target billing currency of the cloned WABA, e.g. "BRL".
	// Meta rejects the request when it matches the source WABA's currency.
	Currency string `json:"currency"`
	// ExtendedCreditID attaches a Line of Credit. Omit for credit-card billing;
	// Lines of Credit are always denominated in USD regardless of the target
	// currency, so this is not the target currency's credit line.
	ExtendedCreditID string `json:"extended_credit_id,omitempty"`
}

// MigrationIntentResponse is the answer to a migration intent.
type MigrationIntentResponse struct {
	MigrationID     string `json:"migration_id"`
	MigrationStatus string `json:"migration_status"`
}

// MigrationDestinationWABA describes the cloned WABA. It only appears once Meta
// has actually created it, so it is absent early in the migration.
type MigrationDestinationWABA struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	Currency                 string `json:"currency"`
	TimezoneID               string `json:"timezone_id"`
	MessageTemplateNamespace string `json:"message_template_namespace"`
}

// MigrationStatusResponse is the state of a migration, as reported by Meta.
type MigrationStatusResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	// DestinationWABA is nil until Meta has created the clone.
	DestinationWABA *MigrationDestinationWABA `json:"destination_waba,omitempty"`
}

// SetPaymentMethodMigrationIntent starts a WABA currency migration: Meta
// validates the source WABA, creates a clone in the target currency, and begins
// copying templates, Flows, app installation and users.
//
// This call does NOT move the phone numbers -- the destination WABA cannot send
// anything until ResumeMigration is called. Coexistence and Authorized Agent
// WABAs are not eligible and are rejected here.
//
// POST /{WABA_ID}/set_payment_method_migration_intent
func (c *Client) SetPaymentMethodMigrationIntent(ctx context.Context, sourceWABAID string, params MigrationIntentRequest) (*MigrationIntentResponse, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal migration intent: %w", err)
	}

	u := fmt.Sprintf("https://graph.facebook.com/%s/%s/set_payment_method_migration_intent",
		c.graphVersion(), url.PathEscape(sourceWABAID))

	req, err := NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, c.httpError(resp)
	}

	out := &MigrationIntentResponse{}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return out, nil
}

// GetMigrationStatus reads the current state of a migration. Poll until the
// status reaches READY_TO_COMPLETE (asset cloning finished, safe to complete) or
// a terminal value.
//
// GET /{MIGRATION_ID}
func (c *Client) GetMigrationStatus(ctx context.Context, migrationID string) (*MigrationStatusResponse, error) {
	u := fmt.Sprintf("https://graph.facebook.com/%s/%s",
		c.graphVersion(), url.PathEscape(migrationID))

	req, err := NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, c.httpError(resp)
	}

	out := &MigrationStatusResponse{}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return out, nil
}

// ResumeMigration completes a migration by detaching the phone numbers from the
// source WABA and attaching them to the destination.
//
// IRREVERSIBLE. After this succeeds the source WABA can never send messages
// again; there is no undo and no rollback. Only call it once the status is
// READY_TO_COMPLETE, and only when the caller genuinely intends the cut-over.
//
// POST /{MIGRATION_ID}/resume_migration
func (c *Client) ResumeMigration(ctx context.Context, migrationID string) (*MigrationStatusResponse, error) {
	u := fmt.Sprintf("https://graph.facebook.com/%s/%s/resume_migration",
		c.graphVersion(), url.PathEscape(migrationID))

	req, err := NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, c.httpError(resp)
	}

	out := &MigrationStatusResponse{}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return out, nil
}

// WABAPhoneNumber is a phone number attached to a WhatsApp Business Account.
type WABAPhoneNumber struct {
	// ID is the phone-number-id used for every send and webhook lookup.
	ID string `json:"id"`
	// DisplayPhoneNumber is the human-formatted number, e.g. "+55 11 3271-0305".
	// It is the only identifier a currency migration cannot change, so it is the
	// key for matching Meta's numbers back to local rows.
	DisplayPhoneNumber string `json:"display_phone_number"`
	VerifiedName       string `json:"verified_name"`
	QualityRating      string `json:"quality_rating"`
	CodeVerification   string `json:"code_verification_status"`
	PlatformType       string `json:"platform_type"`
}

type getWABAPhoneNumbersResponse struct {
	Data   []WABAPhoneNumber `json:"data"`
	Paging struct {
		Cursors struct {
			After string `json:"after"`
		} `json:"cursors"`
		Next string `json:"next"`
	} `json:"paging"`
}

// GetWABAPhoneNumbers lists the phone numbers currently attached to a WABA.
//
// Used after a migration completes to verify which numbers actually moved and
// whether their phone-number-ids survived: Meta's docs imply ids are preserved
// when numbers move between WABAs, but nothing guarantees it, and a changed id
// silently breaks both sending and webhook resolution.
//
// GET /{WABA_ID}/phone_numbers
// It pages to the end rather than reading the first page only: a truncated list
// is indistinguishable from "these numbers did not migrate", which is a
// conclusion callers act on.
func (c *Client) GetWABAPhoneNumbers(ctx context.Context, wabaID string) ([]WABAPhoneNumber, error) {
	all := make([]WABAPhoneNumber, 0, 16)
	after := ""

	// Bounded so a rotating cursor cannot spin forever; 200 per page covers a
	// WABA far larger than Meta's own phone-number ceiling.
	for page := 0; page < 50; page++ {
		q := make(url.Values)
		q.Set("fields", "id,display_phone_number,verified_name,quality_rating,code_verification_status,platform_type")
		q.Set("limit", "200")
		if after != "" {
			q.Set("after", after)
		}

		u := fmt.Sprintf("https://graph.facebook.com/%s/%s/phone_numbers?%s",
			c.graphVersion(), url.PathEscape(wabaID), q.Encode())

		req, err := NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return nil, fmt.Errorf("new request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			err := c.httpError(resp)
			_ = resp.Body.Close()
			return nil, err
		}

		out := &getWABAPhoneNumbersResponse{}
		derr := json.NewDecoder(resp.Body).Decode(out)
		_ = resp.Body.Close()
		if derr != nil {
			return nil, fmt.Errorf("decode response: %w", derr)
		}

		all = append(all, out.Data...)

		// Graph returns a cursor on the last page too; `next` is what actually
		// signals there is more.
		if out.Paging.Next == "" || out.Paging.Cursors.After == "" || out.Paging.Cursors.After == after {
			break
		}
		after = out.Paging.Cursors.After
	}

	return all, nil
}

// WABAInfo is the subset of a WhatsApp Business Account we read back from Meta.
type WABAInfo struct {
	ID string `json:"id"`
	// Currency is the WABA's billing currency, fixed at creation. Meta is the
	// only source of truth for it: our stored copy is a cache that drifts when a
	// WABA is created or migrated outside our flow.
	Currency string `json:"currency"`
	Name     string `json:"name"`
	// TimezoneID matters for billing: volume tiers reset at midnight in the
	// WABA's own timezone, not ours.
	TimezoneID string `json:"timezone_id"`
}

// GetWABAInfo reads a WABA's billing currency and timezone.
//
// GET /{WABA_ID}?fields=id,currency,name,timezone_id
func (c *Client) GetWABAInfo(ctx context.Context, wabaID string) (*WABAInfo, error) {
	q := make(url.Values)
	q.Set("fields", "id,currency,name,timezone_id")

	u := fmt.Sprintf("https://graph.facebook.com/%s/%s?%s",
		c.graphVersion(), url.PathEscape(wabaID), q.Encode())

	req, err := NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, c.httpError(resp)
	}

	out := &WABAInfo{}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return out, nil
}
