package fbgraph

import (
	"strconv"
	"strings"
)

// valid types:
// 		text
// 		image
//      audio
//      video
//      document
//      template
//      hsm (interactive)

type MessageObject struct {
	MessagingProduct string                    `json:"messaging_product"`
	To               string                    `json:"to"`
	Type             string                    `json:"type"`
	RecipientType    string                    `json:"recipient_type"` // default: individual
	Text             *TextObject               `json:"text,omitempty"`
	Template         *TemplateObject           `json:"template,omitempty"`
	Interactive      *InteractiveMessageObject `json:"interactive,omitempty"`
	Image            *MediaObject              `json:"image,omitempty"`
	Audio            *MediaObject              `json:"audio,omitempty"`
	Document         *MediaObject              `json:"document,omitempty"`
	Video            *MediaObject              `json:"video,omitempty"`
	Sticker          *MediaObject              `json:"sticker,omitempty"`
	// TODO: add more objects at:
	// https://developers.facebook.com/docs/whatsapp/cloud-api/reference/messages#text-object
}

type TextObject struct {
	Body       string `json:"body"`
	PreviewURL bool   `json:"preview_url"`
}

type TemplateObject struct {
	// Namespace  string              `json:"namespace,omitempty"`
	Namespace  string              `json:"-"`
	Name       string              `json:"name"`
	Language   *LanguageObject     `json:"language,omitempty"`
	Components []TemplateComponent `json:"components,omitempty"`
}

type LanguageObject struct {
	Code string `json:"code"`
}

type TemplateComponent struct {
	Type       string                       `json:"type"`
	SubType    string                       `json:"sub_type,omitempty"`
	Parameters []TemplateComponentParameter `json:"parameters"`
	Index      *int                         `json:"index,omitempty"`
}

type TemplateComponentParameter struct {
	Type     string              `json:"type"`
	Image    *ImageParameters    `json:"image,omitempty"`
	Payload  string              `json:"payload,omitempty"`
	Text     string              `json:"text,omitempty"`
	Currency *CurrencyParameters `json:"currency,omitempty"`
	DateTime *DateTimeParameters `json:"date_time,omitempty"`
	Video    *MediaObject        `json:"video,omitempty"`
}

// ImageParameters is present when type = "image"
type ImageParameters struct {
	Link string `json:"link,omitempty"`
}

// VideoParameters is present when type = "video"
type VideoParameters struct {
}

type CurrencyParameters struct {
	FallbackValue string `json:"fallback_value"`
	Code          string `json:"code"`
	// Amount multiplied by 1000
	Amount1000 float64 `json:"amount_1000"`
}

type DateTimeParameters struct {
	FallbackValue string `json:"fallback_value"`
}

type MessageObjectResult struct {
	MessagingProduct string          `json:"messaging_product"`
	Contacts         []ContactResult `json:"contacts"`
	Messages         []MessageResult `json:"messages"`
}

type ContactResult struct {
	Input string `json:"input"`
	WAID  string `json:"wa_id"`
}

type MessageResult struct {
	ID string `json:"id"`
}

type GraphError struct {
	Message        string    `json:"message"`
	Type           string    `json:"type"`
	Code           int       `json:"code"`
	ErrorData      ErrorData `json:"error_data"`
	FBTraceID      string    `json:"fbtrace_id"`
	HTTPStatusCode int       `json:"http_status_code"` // this is not originally in the response
}

func (er *GraphError) Error() string {
	sb := new(strings.Builder)
	sb.WriteString(er.Message)
	sb.WriteString("\n")
	sb.WriteString(er.Type)
	sb.WriteString("\n")
	sb.WriteString("Code: " + strconv.Itoa(er.Code))
	sb.WriteString("\n")
	sb.WriteString("fbtrace_id: " + er.FBTraceID)
	return sb.String()
}

func AsGraphError(err error) (*GraphError, bool) {
	if err == nil {
		return nil, false
	}
	if e, ok := err.(*GraphError); ok {
		return e, true
	}
	return nil, false
}

type ErrorData struct {
	MessagingProduct string `json:"messaging_product"`
	Details          string `json:"details"`
}
