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
//      contacts
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
	Contacts         []ContactObject           `json:"contacts,omitempty"`
	Context          *MessageContext           `json:"context,omitempty"`

	Reaction *ReactionObject `json:"reaction,omitempty"`
	// TODO: add more objects at:
	// https://developers.facebook.com/docs/whatsapp/cloud-api/reference/messages#text-object
}

type ReactionObject struct {
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

type MessageContext struct {
	MessageID string `json:"message_id,omitempty"`
}

type TextObject struct {
	Body       string `json:"body"`
	PreviewURL bool   `json:"preview_url"`
}

type ComponentHolder interface {
	GetComponents() []TemplateComponent
}

type TemplateObject struct {
	// Namespace  string              `json:"namespace,omitempty"`
	Namespace  string              `json:"-"`
	Name       string              `json:"name"`
	Language   *LanguageObject     `json:"language,omitempty"`
	Components []TemplateComponent `json:"components,omitempty"`
}

func (t *TemplateObject) GetComponents() []TemplateComponent {
	return t.Components
}

type ContactObject struct {
	Addresses []ContactAddress `json:"addresses,omitempty"`
	Birthday  string           `json:"birthday,omitempty"`
	Emails    []ContactEmail   `json:"emails,omitempty"`
	Name      ContactName      `json:"name"`
	Org       ContactOrg       `json:"org,omitempty"`
	Phones    []ContactPhone   `json:"phones,omitempty"`
	URLs      []ContactURL     `json:"urls,omitempty"`
}

type ContactName struct {
	FormattedName string `json:"formatted_name"`        // full name
	FirstName     string `json:"first_name"`            // first name
	LastName      string `json:"last_name"`             // last name
	MiddleName    string `json:"middle_name,omitempty"` // middle name
	Suffix        string `json:"suffix,omitempty"`      // suffix
	Prefix        string `json:"prefix,omitempty"`      // prefix
}

type ContactPhone struct {
	Phone string `json:"phone"` // phone number
	Type  string `json:"type"`  // phone type
	WAID  string `json:"wa_id"` // whatsapp ID
}

type ContactOrg struct {
	Company    string `json:"company"`    // company name
	Department string `json:"department"` // department
	Title      string `json:"title"`      // job title
}

type ContactURL struct {
	URL  string `json:"url"`  // URL
	Type string `json:"type"` // URL type
}

type ContactAddress struct {
	Street      string `json:"street"`       // street number and name
	City        string `json:"city"`         // city
	State       string `json:"state"`        // state code
	ZIP         string `json:"zip"`          // zip code
	Country     string `json:"country"`      // country name
	CountryCode string `json:"country_code"` // country code
	Type        string `json:"type"`         // address type
}

type ContactEmail struct {
	Email string `json:"email"` // email address
	Type  string `json:"type"`  // email type
}

type LanguageObject struct {
	Code string `json:"code"`
}

type TemplateComponent struct {
	Type       string                       `json:"type"`
	SubType    string                       `json:"sub_type,omitempty"`
	Parameters []TemplateComponentParameter `json:"parameters"`
	Index      *int                         `json:"index,omitempty"`
	Cards      []TemplateCardComponent      `json:"cards,omitempty"`
}

type TemplateCardComponent struct {
	CardIndex  *int                `json:"card_index,omitempty"`
	Components []TemplateComponent `json:"components"`
}

func (t *TemplateCardComponent) GetComponents() []TemplateComponent {
	return t.Components
}

type TemplateComponentParameter struct {
	Type     string              `json:"type"`
	Image    *ImageParameters    `json:"image,omitempty"`
	Payload  string              `json:"payload,omitempty"`
	Text     string              `json:"text,omitempty"`
	Currency *CurrencyParameters `json:"currency,omitempty"`
	DateTime *DateTimeParameters `json:"date_time,omitempty"`
	Video    *VideoParameters    `json:"video,omitempty"`
	Document *MediaObject        `json:"document,omitempty"`
}

// ImageParameters is present when type = "image"
type ImageParameters struct {
	Link string `json:"link,omitempty"`
}

// VideoParameters is present when type = "video"
type VideoParameters struct {
	Link string `json:"link,omitempty"`
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
	Message        string    `json:"message"`          // present in v17+
	Type           string    `json:"type"`             // present in v17+
	Code           int       `json:"code"`             // present in v17+
	ErrorSubcode   int       `json:"error_subcode"`    // present in v17+
	IsTransient    bool      `json:"is_transient"`     // present in v17+
	ErrorUserTitle string    `json:"error_user_title"` // present in v17+
	ErrorUserMsg   string    `json:"error_user_msg"`   // present in v17+
	ErrorData      ErrorData `json:"error_data,omitempty"`
	FBTraceID      string    `json:"fbtrace_id"`
	HTTPStatusCode int       `json:"http_status_code"` // this is not originally in the response
}

// type AutoGenerated struct {
// 	Error struct {
// 		Message        string `json:"message"`
// 		Type           string `json:"type"`
// 		Code           int    `json:"code"`
// 		ErrorSubcode   int    `json:"error_subcode"`
// 		IsTransient    bool   `json:"is_transient"`
// 		ErrorUserTitle string `json:"error_user_title"`
// 		ErrorUserMsg   string `json:"error_user_msg"`
// 		FbtraceID      string `json:"fbtrace_id"`
// 	} `json:"error"`
// }

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
