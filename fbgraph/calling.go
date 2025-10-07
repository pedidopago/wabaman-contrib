package fbgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/*
{
  "calling": {
    "status": "NOT_SET",
    "call_icon_visibility": "NOT_SET",
    "callback_permission_status": "NOT_SET"
  },
  "storage_configuration": {
    "status": "DEFAULT"
  }
}
*/

type CallingStatus string

const (
	CallingStatusNotSet   CallingStatus = "NOT_SET"
	CallingStatusEnabled  CallingStatus = "ENABLED"
	CallingStatusDisabled CallingStatus = "DISABLED"
)

var callingStatusValid = map[CallingStatus]bool{
	CallingStatusNotSet:   false,
	CallingStatusEnabled:  true,
	CallingStatusDisabled: true,
}

func (s CallingStatus) IsValid() bool {
	return callingStatusValid[s]
}

type CallIconButtonVisibility string

const (
	// This status is only valid when the Whatsapp number is not configured for receiving calls.
	CallVisibilityNotSet CallIconButtonVisibility = "NOT_SET"
	// The Call button icon will be displayed in the chat menu bar and the business info page, allowing for unsolicited calls to the business by WhatsApp users.
	CallVisibilityDefault CallIconButtonVisibility = "DEFAULT"
	// The call button icon is hidden in the chat menu bar and the business info page, and all other entry points external to the chat are also disabled. Consumers cannot make unsolicited calls to the business.
	// Your business can still send interactive messages or template messages with a Calling API CTA button.
	CallVisibilityDisableAll CallIconButtonVisibility = "DISABLE ALL"
)

var callIconButtonVisibilityValid = map[CallIconButtonVisibility]bool{
	CallVisibilityNotSet:     false,
	CallVisibilityDefault:    true,
	CallVisibilityDisableAll: true,
}

func (v CallIconButtonVisibility) IsValid() bool {
	return callIconButtonVisibilityValid[v]
}

type CallHoursStatus string

const (
	CallHoursStatusNotSet   CallHoursStatus = "NOT_SET"
	CallHoursStatusEnabled  CallHoursStatus = "ENABLED"
	CallHoursStatusDisabled CallHoursStatus = "DISABLED"
)

var callHoursStatusValid = map[CallHoursStatus]bool{
	CallHoursStatusNotSet:   false,
	CallHoursStatusEnabled:  true,
	CallHoursStatusDisabled: true,
}

func (s CallHoursStatus) IsValid() bool {
	return callHoursStatusValid[s]
}

type CallHoursObject struct {
	Status               CallHoursStatus             `json:"status"`
	TimezoneID           string                      `json:"timezone_id"`
	WeeklyOperatingHours []WeeklyOperatingHourObject `json:"weekly_operating_hours"`
	HolidaySchedule      []HolidayScheduleObject     `json:"holiday_schedule,omitempty"`
}

type DayOfWeek string

const (
	DayOfWeekSunday    DayOfWeek = "SUNDAY"
	DayOfWeekMonday    DayOfWeek = "MONDAY"
	DayOfWeekTuesday   DayOfWeek = "TUESDAY"
	DayOfWeekWednesday DayOfWeek = "WEDNESDAY"
	DayOfWeekThursday  DayOfWeek = "THURSDAY"
	DayOfWeekFriday    DayOfWeek = "FRIDAY"
	DayOfWeekSaturday  DayOfWeek = "SATURDAY"
)

var dayOfWeekValid = map[DayOfWeek]bool{
	DayOfWeekSunday:    true,
	DayOfWeekMonday:    true,
	DayOfWeekTuesday:   true,
	DayOfWeekWednesday: true,
	DayOfWeekThursday:  true,
	DayOfWeekFriday:    true,
	DayOfWeekSaturday:  true,
}

func (d DayOfWeek) IsValid() bool {
	return dayOfWeekValid[d]
}

type OpenCloseTime string

func NewOpenCloseTime(hour, minute int) OpenCloseTime {
	if hour < 0 || hour > 23 {
		hour = hour % 24
	}
	if minute < 0 || minute > 59 {
		minute = minute % 60
	}

	return OpenCloseTime(fmt.Sprintf("%02d%02d", hour, minute))
}

func (t OpenCloseTime) IsValid() bool {
	if len(t) != 4 {
		return false
	}

	for _, c := range t {
		if c < '0' || c > '9' {
			return false
		}
	}

	if t[0] > '2' {
		return false
	}

	if t[0] == '2' && t[1] > '4' {
		return false
	}

	if t[2] > '6' {
		return false
	}

	return true
}

type WeeklyOperatingHourObject struct {
	DayOfWeek DayOfWeek     `json:"day_of_week"`
	OpenTime  OpenCloseTime `json:"open_time"`
	CloseTime OpenCloseTime `json:"close_time"`
}

type HolidayScheduleObject struct {
	Date      string        `json:"date"` // YYYY-MM-DD
	StartTime OpenCloseTime `json:"start_time"`
	EndTime   OpenCloseTime `json:"end_time"`
}

// Configure call signaling via signal initiation protocol (SIP).
//
// Note: When SIP is enabled, you cannot use calling related endpoints and will not receive calling related webhooks.
//
// https://developers.facebook.com/docs/whatsapp/cloud-api/calling/sip
type SipServerObject struct {
	Hostname             string            `json:"hostname"`
	Port                 int               `json:"port"`
	RequestURIUserParams map[string]string `json:"request_uri_user_params,omitempty"`
}

type CallingSip struct {
	Status  string            `json:"status"` // ENABLED, DISABLED (default)
	Servers []SipServerObject `json:"servers,omitzero"`
}

type WhatsappSettings struct {
	Calling struct {
		Status                   CallingStatus            `json:"status"`
		CallIconVisibility       CallIconButtonVisibility `json:"call_icon_visibility"`
		CallHours                CallHoursObject          `json:"call_hours"`
		CallbackPermissionStatus CallingStatus            `json:"callback_permission_status"` // ENABLED, DISABLED
		Sip                      *CallingSip              `json:"sip,omitempty"`
	} `json:"calling"`
}

func (c *Client) GetWhatsappSettings(ctx context.Context, whatsappID string) (*WhatsappSettings, error) {
	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/settings", apiVersion, whatsappID)

	req, err := NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, c.errorFromResponse(resp)
	}

	result := new(WhatsappSettings)

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// UpdateWhatsappSettings: Use this endpoint to update call settings configuration for an individual business phone number.
//
//Possible errors that can occur:
//
// Permissions/Authorization errors
// Invalid status
// Invalid schedule for call_hours
// Holiday given in call_hours is a past date
// Timezone is invalid in call_hours
// weekly_operating_hours in call_hours cannot be empty
// Date format in holiday_schedule for call_hours is invalid
// More than 2 entries not allowed in weekly_operating_hours schedule in call_hours
// Overlapping schedule in call_hours is not allowed
// Calling restriction errors

func (c *Client) UpdateWhatsappSettings(ctx context.Context, whatsappID string, settings *WhatsappSettings) error {
	c.lastErrorRawBody = ""
	c.lastGraphError = nil

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/settings", apiVersion, whatsappID)

	jd, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	rbuf := bytes.NewBuffer(jd)
	rbuf.Write(jd)

	req, err := NewRequestWithContext(ctx, http.MethodPost, url, rbuf)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.errorFromResponse(resp)
	}

	io.Copy(io.Discard, resp.Body)

	return nil
}

type MessageLimitingTier string

const (
	Tier50        MessageLimitingTier = "TIER_50"
	Tier250       MessageLimitingTier = "TIER_250"
	Tier1K        MessageLimitingTier = "TIER_1K"
	Tier10K       MessageLimitingTier = "TIER_10K"
	Tier100K      MessageLimitingTier = "TIER_100K"
	TierNotSet    MessageLimitingTier = "TIER_NOT_SET" // Indicates the business phone number has not been used to send a message yet.
	TierUnlimited MessageLimitingTier = "TIER_UNLIMITED"
)

func (c *Client) GetMessagingLimitingTier(ctx context.Context, whatsappID string) (MessageLimitingTier, error) {
	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s?fields=messaging_limit_tier", apiVersion, whatsappID)

	req, err := NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", c.errorFromResponse(resp)
	}

	result := struct {
		MessagingLimitTier MessageLimitingTier `json:"messaging_limit_tier"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.MessagingLimitTier, nil
}

func (tier MessageLimitingTier) CanEnableCallingAPI() bool {
	switch tier {
	case TierNotSet, Tier50, Tier250, "":
		return false
	}

	return true
}

type WebhookFieldDesc struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type WebhookObject struct {
	Object      string             `json:"object"`
	CallbackURL string             `json:"callback_url"`
	Active      bool               `json:"active"`
	Fields      []WebhookFieldDesc `json:"fields"`
}

func (c *Client) GetAppSubscribedWebhooks(ctx context.Context, appID, appSecret string) ([]WebhookObject, error) {
	c.lastErrorRawBody = ""
	c.lastGraphError = nil

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	specialToken := fmt.Sprintf("%s|%s", appID, appSecret)

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/subscriptions", apiVersion, appID)

	req, err := NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", specialToken))
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.errorFromResponse(resp)
	}

	result := struct {
		Data []WebhookObject `json:"data"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Data, nil
}

func IsSubscribedToCalls(objs []WebhookObject, minVersion string) (bool, error) {
	for _, obj := range objs {
		if obj.Object == "whatsapp_business_account" {
			if !obj.Active {
				continue
			}

			for _, field := range obj.Fields {
				if field.Name == "calls" {
					if minVersion == "" {
						return true, nil
					}

					cmp, err := CompareGraphAPIVersions(field.Version, minVersion)
					if err != nil {
						return false, fmt.Errorf("compare versions: %w", err)
					}

					if cmp >= 0 {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}
