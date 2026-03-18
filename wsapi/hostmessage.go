package wsapi

import (
	"time"

	"github.com/pedidopago/wabaman-contrib/fbgraph"
)

// HostMessage is a message sent from the business (agent) to a WhatsApp client.
type HostMessage struct {
	// Internal ID (bigint)
	ID uint64 `json:"id"`
	// Base64 whatsapp message id from graph API.
	// It is guaranteed to be unique.
	WABAMessageID string `json:"waba_message_id"`
	// Internal Phone ID
	PhoneID         uint   `json:"phone_id"`
	HostPhoneNumber string `json:"host_phone_number"`
	// The id (phone number) of the recipient (client).
	WABARecipientID         string                            `json:"waba_recipient_id"`
	// The business-scoped user ID (BSUID) of the recipient.
	UserID                  string                            `json:"user_id,omitempty"`
	WABATimestamp           time.Time                         `json:"waba_timestamp"`
	OriginalFailedMessageID string                            `json:"original_failed_message_id,omitempty"`
	FailedMessageRetryChain uint                              `json:"failed_message_retry_chain,omitempty"`
	Type                    string                            `json:"type"`
	Text                    *Text                             `json:"text,omitempty"`
	Document                *Document                         `json:"document,omitempty"`
	Video                   *Video                            `json:"video,omitempty"`
	Image                   *Image                            `json:"image,omitempty"`
	Audio                   *Audio                            `json:"audio,omitempty"`
	Sticker                 *Sticker                          `json:"sticker,omitempty"`
	Template                *HostTemplate                     `json:"template,omitempty"`
	Interactive             *fbgraph.InteractiveMessageObject `json:"interactive,omitempty"`
	Contacts                []fbgraph.ContactObject           `json:"contacts,omitempty"`
	ObjectType              string                            `json:"object_type,omitempty"`
	AgentID                 string                            `json:"agent_id,omitempty"`
	AgentName               string                            `json:"agent_name,omitempty"`
	FailedReason            *SentMessageFailedReason          `json:"failed_reason,omitempty"`
	Preview                 string                            `json:"preview,omitempty"`
	Origin                  string                            `json:"origin,omitempty"`
	CreatedAt               time.Time                         `json:"created_at,omitempty"`
	CreatedAtNano           int64                             `json:"created_at_nano,omitempty"`
	Context                 *MessageContext                   `json:"context,omitempty"`
	Reactions               []MessageReaction                 `json:"reactions,omitempty"`
	Metadata                map[string]any                    `json:"metadata,omitempty"`
}

// GetOrigin returns the origin of the host message, or an empty string if nil.
func (m *HostMessage) GetOrigin() string {
	if m == nil {
		return ""
	}

	return m.Origin
}

// FBStatusObjectError represents an error returned by the WhatsApp Cloud API
// in a message status webhook callback.
type FBStatusObjectError struct {
	Code           int    `json:"code"`
	Title          string `json:"title"`
	Href           string `json:"href"`
	LocalizedError string `json:"localized_error,omitempty"` // Pedido Pago custom field
}

// SentMessageFailedReason contains the error details when a message fails to be delivered.
type SentMessageFailedReason struct {
	Code           int                   `json:"code,omitempty"`
	Title          string                `json:"title,omitempty"`
	Href           string                `json:"href,omitempty"`
	Errors         []FBStatusObjectError `json:"errors,omitempty"`
	LocalizedError string                `json:"localized_error,omitempty"` // Pedido Pago custom field
}

// HostTemplate represents a WhatsApp message template as sent by the host.
// It contains the original template reference, the graph API object, and the parsed result.
type HostTemplate struct {
	//Deprecated: use GraphObject.Name or Parsed.TemplateName
	Name        string                  `json:"name,omitempty"`
	GraphObject *fbgraph.TemplateObject `json:"graph_object,omitempty"`
	Original    *TemplateRef            `json:"original,omitempty"`
	Parsed      *ParsedTemplate         `json:"parsed,omitempty"`
}

// TemplateCategory classifies a WhatsApp message template per the Cloud API categories.
type TemplateCategory string

const (
	TplCategoryTransactional         TemplateCategory = "transactional"
	TplCategoryMarketing             TemplateCategory = "marketing"
	TplCategoryDisposableCredentials TemplateCategory = "disposable_credentials"
)

// TemplateHeaderType specifies the type of header in a WhatsApp message template.
type TemplateHeaderType string

const (
	TplHeaderTypeNone  TemplateHeaderType = "none"
	TplHeaderTypeText  TemplateHeaderType = "text"
	TplHeaderTypeMedia TemplateHeaderType = "media"
)

// Deprecated: use TemplateButtonType
type TemplateButtonsType string

// Deprecated: use TemplateButtonType
const (
	TplButtonsTypeNone         TemplateButtonsType = "none"
	TplButtonsTypeCallToAction TemplateButtonsType = "call_to_action"
	TplButtonsTypeQuickReply   TemplateButtonsType = "quick_reply"
)

// TemplateButtonType specifies the type of button in a WhatsApp message template.
type TemplateButtonType string

const (
	TemplateButtonCall       TemplateButtonType = "call"
	TemplateButtonURL        TemplateButtonType = "url"
	TemplateButtonQuickReply TemplateButtonType = "quick_reply"
)

// TemplateRef is a reference to a WhatsApp message template stored in Wabaman.
type TemplateRef struct {
	ID         uint64           `json:"id"`
	BusinessID uint             `json:"business_id"`
	Name       string           `json:"name"`
	Category   TemplateCategory `json:"category"`
	CreatedAt  time.Time        `json:"created_at"`

	Languages []*TemplateRefLanguage `json:"languages"`
}

// TemplateRefLanguage is a single language variant of a [TemplateRef].
type TemplateRefLanguage struct {
	LanguageCode string `json:"language_code"`

	Header *TplRefHeader `json:"header"`
	Body   string        `json:"body"`
	Footer string        `json:"footer"`

	// Deprecated: use TemplateButtonType
	ButtonsType TemplateButtonsType `json:"buttons_type"`

	// Deprecated: use Buttons
	QuickReplyButtons []TplQuickReplyButton `json:"quick_reply_buttons,omitempty"`
	// Deprecated: use Buttons
	CallToActionButtons []TplCallToActionButton `json:"call_to_action_buttons,omitempty"`

	Buttons []TplButton `json:"buttons,omitempty"`
	Cards   []TplCard   `json:"cards,omitempty"`
}

// TplRefHeader describes the header section of a template language variant.
type TplRefHeader struct {
	Type           TemplateHeaderType `json:"header_type"`
	ContentExample string             `json:"content_example"`
}

// Deprecated: use TplButton
type TplQuickReplyButton struct {
	Text string `json:"text"`
}

// Deprecated: use TplButton
type TplCallToActionButton struct {
	Type TemplateButtonType   `json:"type"`
	Text string               `json:"text"`
	URL  *TplCallToActionURL  `json:"url,omitempty"`
	Call *TplCallToActionCall `json:"call,omitempty"`
}

// TplCallToActionURL contains the URL configuration for a call-to-action template button.
type TplCallToActionURL struct {
	Type string `json:"type"` // static,dynamic
	Href string `json:"href"`
}

// TplCallToActionCall contains the phone number configuration for a call-to-action template button.
type TplCallToActionCall struct {
	CC    string `json:"cc"`
	Phone string `json:"phone"`
}

// ParsedTemplate is the final rendered form of a message template after parameter substitution.
type ParsedTemplate struct {
	TemplateName string                `json:"template_name"`
	LanguageCode string                `json:"language_code"`
	Header       *ParsedTemplateHeader `json:"header,omitempty"`
	Body         string                `json:"body"`
	Footer       string                `json:"footer"`
	Buttons      []TplButton           `json:"buttons,omitempty"`
	Cards        []TplCard             `json:"cards,omitempty"`

	// Deprecated fields

	// Deprecated: use Buttons
	ButtonsType TemplateButtonsType `json:"buttons_type"`

	// Deprecated: use Buttons
	QuickReplyButtons []TplQuickReplyButton `json:"quick_reply_buttons" deprecated:"true"`
	// Deprecated: use Buttons
	CallToActionButtons []TplCallToActionButton `json:"call_to_action_buttons" deprecated:"true"`
}

// ParsedTemplateHeader is the rendered header of a [ParsedTemplate].
type ParsedTemplateHeader struct {
	HeaderType     TemplateHeaderType `json:"header_type"`
	ContentExample string             `json:"content_example"`
	Content        string             `json:"content"`
	Type           string             `json:"type"`
}

// TplCard is a single card within a carousel template message.
type TplCard struct {
	Header *ParsedTemplateHeader `json:"header,omitempty"`
	Body   string                `json:"body"`
	//Footer  string                `json:"footer,omitempty"` //gabs: removed footer since there is no support: https://developers.facebook.com/docs/whatsapp/business-management-api/message-templates/media-card-carousel-templates/
	Buttons []TplButton `json:"buttons,omitempty"`
}

// TplButton represents a button within a WhatsApp message template.
type TplButton struct {
	Type TemplateButtonType   `json:"type"`
	Text string               `json:"text"`
	URL  *TplCallToActionURL  `json:"url,omitempty"`
	Call *TplCallToActionCall `json:"call,omitempty"`
}
