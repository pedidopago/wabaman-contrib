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

const (
	MTComponentHeader  MessageTemplateComponentType = "HEADER"
	MTComponentBody    MessageTemplateComponentType = "BODY"
	MTComponentFooter  MessageTemplateComponentType = "FOOTER"
	MTComponentButtons MessageTemplateComponentType = "BUTTONS"
)

type MessageTemplateComponentFormat string

const (
	MTCFormatVideo    MessageTemplateComponentFormat = "VIDEO"
	MTCFormatImage    MessageTemplateComponentFormat = "IMAGE"
	MTCFormatDocument MessageTemplateComponentFormat = "DOCUMENT" //TODO: check
)

type MessageTemplateButtonType string

const (
	MTBTypeQuickReply  MessageTemplateButtonType = "QUICK_REPLY"
	MTBTypeURL         MessageTemplateButtonType = "URL"
	MTBTypePhoneNumber MessageTemplateButtonType = "PHONE_NUMBER"
)

type MessageTemplateCategory string

const (
	MTCategoryTransactional        MessageTemplateCategory = "TRANSACTIONAL"
	MTCategoryMarketing            MessageTemplateCategory = "MARKETING"
	MTCategoryTicketUpdate         MessageTemplateCategory = "TICKET_UPDATE"
	MTCategoryIssueResolution      MessageTemplateCategory = "ISSUE_RESOLUTION"
	MTCategoryIssueOneTimePassword MessageTemplateCategory = "OTP"
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

type MessageTemplateScore struct {
	Score   string   `json:"score"`
	Reasons []string `json:"reasons,omitempty"`
}

type MessageTemplateComponent struct {
	Type    MessageTemplateComponentType `json:"type"`
	Text    string                       `json:"text,omitempty"`
	Format  string                       `json:"format,omitempty"`
	Example *MessageTemplateExample      `json:"example,omitempty"`
	Buttons []MessageTemplateButton      `json:"buttons,omitempty"`
}

type MessageTemplateExample struct {
	HeaderHandle Slice[string] `json:"header_handle"`
	BodyText     Slice[string] `json:"body_text"`
}

type MessageTemplateButton struct {
	Type        MessageTemplateButtonType `json:"type"`
	Text        string                    `json:"text"`
	URL         string                    `json:"url,omitempty"`
	Example     []string                  `json:"example,omitempty"`
	PhoneNumber string                    `json:"phone_number,omitempty"`
}

type MessageTemplatesPaging struct {
	Cursors  MessageTemplatesCursors `json:"cursors"`
	Next     string                  `json:"next,omitempty"`
	Previous string                  `json:"previous,omitempty"`
}

type MessageTemplatesCursors struct {
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

func (c *Client) GetMessageTemplates(ctx context.Context, params GetMessageTemplatesParameters) (*GetMessageTemplatesResponse, error) {
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

	url := fmt.Sprintf("https://graph.facebook.com/v15.0/%s/message_templates?%s", params.WhatsAppBusinessAccountID, encfields.Encode())

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

func (c *Client) CreateMessageTemplate(ctx context.Context, wabaID string, template MessageTemplate) (id string, err error) {
	url := fmt.Sprintf("https://graph.facebook.com/v15.0/%s/message_templates", wabaID)

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(template); err != nil {
		return "", fmt.Errorf("encode template: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
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

func (c *Client) UpdateMessageTemplate(ctx context.Context, templateID string, components []TemplateComponent) error {
	url := fmt.Sprintf("https://graph.facebook.com/v15.0/%s", templateID)

	cpstruct := struct {
		Components []TemplateComponent `json:"components"`
	}{Components: components}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(cpstruct); err != nil {
		return fmt.Errorf("encode template update: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
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
		return c.errorFromResponse(resp)
	}

	result := struct {
		Success bool `json:"success"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func (c *Client) DeleteMessageTemplate(ctx context.Context, templateID string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v15.0/%s", templateID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
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
		return c.errorFromResponse(resp)
	}

	result := struct {
		Success bool `json:"success"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
