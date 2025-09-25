package fbgraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type MessageTemplateComponentType string

func (ct MessageTemplateComponentType) Equals(other MessageTemplateComponentType) bool {
	return strings.EqualFold(string(ct), string(other))
}

const (
	MTComponentHeader   MessageTemplateComponentType = "HEADER"
	MTComponentBody     MessageTemplateComponentType = "BODY"
	MTComponentFooter   MessageTemplateComponentType = "FOOTER"
	MTComponentButtons  MessageTemplateComponentType = "BUTTONS"
	MTComponentCarousel MessageTemplateComponentType = "CAROUSEL" //TODO: check casing
)

type MessageTemplateComponentFormat string

const (
	MTCFormatVideo    MessageTemplateComponentFormat = "VIDEO"
	MTCFormatImage    MessageTemplateComponentFormat = "IMAGE"
	MTCFormatDocument MessageTemplateComponentFormat = "DOCUMENT"
	MTCFormatText     MessageTemplateComponentFormat = "TEXT"
)

type MessageTemplateButtonType string

const (
	MTBTypeQuickReply  MessageTemplateButtonType = "QUICK_REPLY"
	MTBTypeURL         MessageTemplateButtonType = "URL"
	MTBTypePhoneNumber MessageTemplateButtonType = "PHONE_NUMBER"
	MTBTypeOTP         MessageTemplateButtonType = "OTP"
)

type MessageTemplateCategory string

const (
	MTCategoryUtility         MessageTemplateCategory = "UTILITY"
	MTCategoryMarketing       MessageTemplateCategory = "MARKETING"
	MTCategoryTicketUpdate    MessageTemplateCategory = "TICKET_UPDATE"
	MTCategoryIssueResolution MessageTemplateCategory = "ISSUE_RESOLUTION"
	MTCategoryAuthentication  MessageTemplateCategory = "AUTHENTICATION"
)

var templateParamsDefaultFields = []string{
	"category", "language", "name", "quality_score", "rejected_reason", "status", "content", "components",
}

// fields=category,language,name,quality_score,status,content,components&limit=3&after=MgZDZD
type GetMessageTemplatesParameters struct {
	WhatsAppBusinessAccountID string
	Fields                    []string
	Limit                     int
	After                     string
}

type GetMessageTemplatesResponse struct {
	Data   []MessageTemplate      `json:"data"`
	Paging MessageTemplatesPaging `json:"paging"`
}

type MessageTemplate struct {
	Category       MessageTemplateCategory    `json:"category"`
	Language       string                     `json:"language"`
	Name           string                     `json:"name"`
	QualityScore   *MessageTemplateScore      `json:"quality_score,omitempty"`
	RejectedReason string                     `json:"rejected_reason,omitempty"`
	Status         string                     `json:"status,omitempty"`
	Components     []MessageTemplateComponent `json:"components"`
	ID             string                     `json:"id,omitempty"`
}

func ConvertMessageTemplateToNew(tpl MessageTemplate, allowCategoryChange bool) NewMessageTemplate {
	return NewMessageTemplate{
		MessageTemplate:     tpl,
		AllowCategoryChange: allowCategoryChange,
	}
}

type NewMessageTemplate struct {
	MessageTemplate     `json:",inline"`
	AllowCategoryChange bool `json:"allow_category_change,omitempty"`
}

type MessageTemplateScore struct {
	Score   string   `json:"score"`
	Reasons []string `json:"reasons,omitempty"`
}

type MessageTemplateComponent struct {
	Type                      MessageTemplateComponentType `json:"type"`
	Text                      string                       `json:"text,omitempty"`
	Format                    string                       `json:"format,omitempty"`
	Example                   *MessageTemplateExample      `json:"example,omitempty"`
	Buttons                   []MessageTemplateButton      `json:"buttons,omitempty"`
	Cards                     []MessageTemplateCard        `json:"cards,omitempty"`                       // for type `carousel`
	AddSecurityRecommendation *bool                        `json:"add_security_recommendation,omitempty"` // only for authentication templates
	CodeExpirationMinutes     *int                         `json:"code_expiration_minutes,omitempty"`     // only for authentication templates
}

type MessageTemplateCardComponent struct {
	Type    MessageTemplateComponentType `json:"type"`
	Text    string                       `json:"text,omitempty"`
	Format  string                       `json:"format,omitempty"`
	Example *MessageTemplateExample      `json:"example,omitempty"`
	Buttons []MessageTemplateButton      `json:"buttons,omitempty"`
}

type MessageTemplateExample struct {
	HeaderHandle []string   `json:"header_handle,omitempty"`
	BodyText     [][]string `json:"body_text,omitempty"`
	HeaderText   []string   `json:"header_text,omitempty"`
}

type MessageTemplateButton struct {
	Type          MessageTemplateButtonType `json:"type"`
	OTPType       OTPType                   `json:"otp_type,omitempty"`      // for OTP types
	AutofillText  string                    `json:"autofill_text,omitempty"` // for OTP type ONE_TAP
	Text          string                    `json:"text"`
	URL           string                    `json:"url,omitempty"`
	Example       []string                  `json:"example,omitempty"`
	PhoneNumber   string                    `json:"phone_number,omitempty"`
	PackageName   string                    `json:"package_name,omitempty"`   // for OTP type ONE_TAP
	SignatureHash string                    `json:"signature_hash,omitempty"` // for OTP type ONE_TAP
}

// for type `carousel`
//
// At least 1 button required, maximum 2; button types can be mixed
type MessageTemplateCard struct {
	CardIndex  *int                           `json:"card_index,omitzero"`
	Components []MessageTemplateCardComponent `json:"components"`
}

type OTPType string

const (
	OTPTypeOneTap   OTPType = "ONE_TAP"
	OTPTypeCopyCode OTPType = "COPY_CODE"
)

type MessageTemplatesPaging struct {
	Cursors  MessageTemplatesCursors `json:"cursors"`
	Next     string                  `json:"next,omitempty"`
	Previous string                  `json:"previous,omitempty"`
}

type MessageTemplatesCursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

func (c *Client) GetMessageTemplate(ctx context.Context, id string) (*MessageTemplate, error) {
	apiversion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiversion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s", apiversion, id)

	req, err := NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
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
	result := &MessageTemplate{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

func (c *Client) GetMessageTemplates(ctx context.Context, params GetMessageTemplatesParameters) (*GetMessageTemplatesResponse, error) {
	apiversion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiversion = c.GraphAPIVersion
	}

	if params.Fields == nil {
		params.Fields = templateParamsDefaultFields
	}

	encfields := make(url.Values)
	encfields.Set("fields", strings.Join(params.Fields, ","))
	if params.After != "" {
		encfields.Set("after", params.After)
	}
	if params.Limit != 0 {
		encfields.Set("limit", fmt.Sprintf("%d", params.Limit))
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/message_templates?%s", apiversion, params.WhatsAppBusinessAccountID, encfields.Encode())

	req, err := NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
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
	result := &GetMessageTemplatesResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

func (c *Client) CreateMessageTemplate(ctx context.Context, wabaID string, template NewMessageTemplate) (id string, err error) {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/message_templates", apiVersion, wabaID)

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(template); err != nil {
		return "", fmt.Errorf("encode template: %w", err)
	}

	req, err := NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		gerr := c.errorFromResponse(resp)

		if ge, ok := AsGraphError(gerr); ok {
			if ge.Code == 4 {
				return "", ErrApplicationRateLimitReached
			}
		}

		return "", gerr
	}

	result := struct {
		ID string `json:"id"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.ID, nil
}

func (c *Client) UpdateMessageTemplateCategory(ctx context.Context, templateID string, newCategory MessageTemplateCategory) error {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s", apiVersion, templateID)

	cpstruct := struct {
		Category MessageTemplateCategory `json:"category"`
	}{Category: newCategory}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cpstruct); err != nil {
		return fmt.Errorf("encode template update: %w", err)
	}

	req, err := NewRequest(http.MethodPost, url, buf)
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

	result := struct {
		Success bool `json:"success"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (c *Client) UpdateMessageTemplate(ctx context.Context, templateID string, components []MessageTemplateComponent) error {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	url := fmt.Sprintf("https://graph.facebook.com/%s/%s", apiVersion, templateID)

	cpstruct := struct {
		Components []MessageTemplateComponent `json:"components"`
	}{Components: components}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cpstruct); err != nil {
		return fmt.Errorf("encode template update: %w", err)
	}

	req, err := NewRequest(http.MethodPost, url, buf)
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

	result := struct {
		Success bool `json:"success"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (c *Client) DeleteMessageTemplate(ctx context.Context, whatsappBusinessAccountID, templateName string) error {
	c.lastGraphError = nil
	c.lastErrorRawBody = ""

	apiVersion := DefaultGraphAPIVersion
	if c.GraphAPIVersion != "" {
		apiVersion = c.GraphAPIVersion
	}

	urlv := make(url.Values)
	urlv.Set("name", templateName)
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/message_templates?%s", apiVersion, whatsappBusinessAccountID, urlv.Encode())

	req, err := NewRequest(http.MethodDelete, url, nil)
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

	result := struct {
		Success bool `json:"success"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
