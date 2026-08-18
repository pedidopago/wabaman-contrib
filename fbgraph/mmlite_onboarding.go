package fbgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// MMLiteOnboardingStatus is where a portfolio or a WABA stands in Meta's
// Marketing Messages Lite onboarding.
//
// The two grains do not share a vocabulary, which is the whole reason this is
// one open string type rather than two enums: a portfolio signs the Terms of
// Service, and a WABA becomes sendable. Signing does not by itself make any
// WABA sendable, and Meta reports the two independently.
type MMLiteOnboardingStatus string

// Portfolio-grain statuses. Read from GET /{BUSINESS_ID}.
const (
	// MMLiteNotStarted means nobody has asked this portfolio yet.
	MMLiteNotStarted MMLiteOnboardingStatus = "NOT_STARTED"
	// MMLiteRequestSent means the portfolio has been invited and the alert is
	// waiting for a human inside their WhatsApp Manager.
	MMLiteRequestSent MMLiteOnboardingStatus = "REQUEST_SENT"
	// MMLiteTermsOfServiceSigned means a person with authority in that
	// portfolio accepted the Terms. It cascades to every eligible WABA the
	// portfolio owns.
	MMLiteTermsOfServiceSigned MMLiteOnboardingStatus = "TERM_OF_SERVICE_SIGNED"
)

// WABA-grain statuses. Read from GET /{WABA_ID}.
const (
	// MMLiteEligible means the WABA could be onboarded but is not yet.
	MMLiteEligible MMLiteOnboardingStatus = "ELIGIBLE"
	// MMLiteOnboarded is the only status from which a marketing template may be
	// sent through the Marketing Messages Lite endpoint.
	MMLiteOnboarded MMLiteOnboardingStatus = "ONBOARDED"
	// MMLitePendingValidPaymentMethod means Meta wants a payment method on the
	// portfolio before it will onboard the WABA. Under a Solution Partner
	// credit line this should not appear; if it does, it is a billing question,
	// not an engineering one.
	MMLitePendingValidPaymentMethod MMLiteOnboardingStatus = "PENDING_VALID_PAYMENT_METHOD"
)

// mmLiteOnboardingStatusField is deliberately requested ALONE.
//
// Graph rejects an entire request when the token cannot see any one requested
// field. Appending this to wabaInfoFields would therefore let a token missing
// the permission take the CURRENCY down with it -- and the currency is what
// phone registration and InitiateWABAMigration depend on. Onboarding status is
// a nice-to-have riding along; it must never cost us a field that is not.
const mmLiteOnboardingStatusField = "marketing_messages_onboarding_status"

// GetMMLiteOnboardingStatus reads the Marketing Messages Lite onboarding status
// of either a WABA or a business portfolio.
//
// GET /{WABA_ID}?fields=marketing_messages_onboarding_status
// GET /{BUSINESS_ID}?fields=marketing_messages_onboarding_status
//
// The same edge serves both grains, so one method covers both -- but the
// statuses it can return differ by what you passed, and the caller is expected
// to know which it asked about.
//
// An empty status with a nil error means Meta answered without the field. That
// is "Meta will not tell us", never "NOT_STARTED": treating the absence as a
// status would silently invent a fact about somebody's legal acceptance.
func (c *Client) GetMMLiteOnboardingStatus(ctx context.Context, id string) (MMLiteOnboardingStatus, error) {
	q := make(url.Values)
	q.Set("fields", mmLiteOnboardingStatusField)

	u := fmt.Sprintf("https://graph.facebook.com/%s/%s?%s",
		c.graphVersion(), url.PathEscape(id), q.Encode())

	req, err := NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", c.httpError(resp)
	}

	var out struct {
		Status MMLiteOnboardingStatus `json:"marketing_messages_onboarding_status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return out.Status, nil
}
