package fbgraph

import (
	"context"
	"encoding/json"
	"fmt"
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

type CallHoursStatus string

const (
	CallHoursStatusNotSet   CallHoursStatus = "NOT_SET"
	CallHoursStatusEnabled  CallHoursStatus = "ENABLED"
	CallHoursStatusDisabled CallHoursStatus = "DISABLED"
)

type CallHoursObject struct {
	Status               CallHoursStatus             `json:"status"`
	TimezoneID           string                      `json:"timezone_id"`
	WeeklyOperatingHours []WeeklyOperatingHourObject `json:"weekly_operating_hours"`
	HolidaySchedule      []HolidayScheduleObject     `json:"holiday_schedule"`
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

type SipServerObject struct {
	Hostname             string            `json:"hostname"`
	Port                 int               `json:"port"`
	RequestURIUserParams map[string]string `json:"request_uri_user_params,omitempty"`
}

type WhatsappSettings struct {
	Calling struct {
		Status                   CallingStatus            `json:"status"`
		CallIconVisibility       CallIconButtonVisibility `json:"call_icon_visibility"`
		CallHours                CallHoursObject          `json:"call_hours"`
		CallbackPermissionStatus string                   `json:"callback_permission_status"` // ENABLED, DISABLED
		Sip                      struct {
			Status  string            `json:"status"` // ENABLED, DISABLED (default)
			Servers []SipServerObject `json:"servers,omitzero"`
		} `json:"sip"`
	} `json:"calling"`
}

func (c *Client) GetWhatsappSettings(ctx context.Context, whatsappID string) (*WhatsappSettings, error) {
	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/settings", apiVersion, whatsappID)

	req, err := http.NewRequest(http.MethodGet, url, nil)

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
