package wsapi

import (
	"time"

	"github.com/pedidopago/wabaman-contrib/fbgraph"
)

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

func (m *HostMessage) GetOrigin() string {
	if m == nil {
		return ""
	}

	return m.Origin
}

type FBStatusObjectError struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
	Href  string `json:"href"`
}

type SentMessageFailedReason struct {
	Code   int                   `json:"code,omitempty"`
	Title  string                `json:"title,omitempty"`
	Href   string                `json:"href,omitempty"`
	Errors []FBStatusObjectError `json:"errors,omitempty"`
}

type HostTemplate struct {
	//Deprecated: use GraphObject.Name or Parsed.TemplateName
	Name        string                  `json:"name,omitempty"`
	GraphObject *fbgraph.TemplateObject `json:"graph_object,omitempty"`
	Original    *TemplateRef            `json:"original,omitempty"`
	Parsed      *ParsedTemplate         `json:"parsed,omitempty"`
}

type TemplateCategory string

const (
	TplCategoryTransactional         TemplateCategory = "transactional"
	TplCategoryMarketing             TemplateCategory = "marketing"
	TplCategoryDisposableCredentials TemplateCategory = "disposable_credentials"
)

type TemplateHeaderType string

const (
	TplHeaderTypeNone  TemplateHeaderType = "none"
	TplHeaderTypeText  TemplateHeaderType = "text"
	TplHeaderTypeMedia TemplateHeaderType = "media"
)

type TemplateButtonsType string

const (
	TplButtonsTypeNone         TemplateButtonsType = "none"
	TplButtonsTypeCallToAction TemplateButtonsType = "call_to_action"
	TplButtonsTypeQuickReply   TemplateButtonsType = "quick_reply"
)

type TemplateCTAButtonType string

const (
	TplCTABCall TemplateCTAButtonType = "call"
	TplCTABURL  TemplateCTAButtonType = "url"
)

type TemplateRef struct {
	ID         uint64           `json:"id"`
	BusinessID uint             `json:"business_id"`
	Name       string           `json:"name"`
	Category   TemplateCategory `json:"category"`
	CreatedAt  time.Time        `json:"created_at"`

	Languages []*TemplateRefLanguage `json:"languages"`
}

type TemplateRefLanguage struct {
	LanguageCode string `json:"language_code"`

	Header      *TplRefHeader       `json:"header"`
	Body        string              `json:"body"`
	Footer      string              `json:"footer"`
	ButtonsType TemplateButtonsType `json:"buttons_type"`

	QuickReplyButtons   []TplQuickReplyButton   `json:"quick_reply_buttons,omitempty"`
	CallToActionButtons []TplCallToActionButton `json:"call_to_action_buttons,omitempty"`
}

type TplRefHeader struct {
	Type           TemplateHeaderType `json:"header_type"`
	ContentExample string             `json:"content_example"`
}

type TplQuickReplyButton struct {
	Text string `json:"text"`
}

type TplCallToActionButton struct {
	Type TemplateCTAButtonType `json:"type"`
	Text string                `json:"text"`
	URL  *TplCallToActionURL   `json:"url,omitempty"`
	Call *TplCallToActionCall  `json:"call,omitempty"`
}

type TplCallToActionURL struct {
	Type string `json:"type"` // static,dynamic
	Href string `json:"href"`
}

type TplCallToActionCall struct {
	CC    string `json:"cc"`
	Phone string `json:"phone"`
}

type ParsedTemplate struct {
	TemplateName        string                  `json:"template_name"`
	LanguageCode        string                  `json:"language_code"`
	Header              *ParsedTemplateHeader   `json:"header,omitempty"`
	Body                string                  `json:"body"`
	Footer              string                  `json:"footer"`
	ButtonsType         TemplateButtonsType     `json:"buttons_type"`
	QuickReplyButtons   []TplQuickReplyButton   `json:"quick_reply_buttons"`
	CallToActionButtons []TplCallToActionButton `json:"call_to_action_buttons"`
}

type ParsedTemplateHeader struct {
	HeaderType     TemplateHeaderType `json:"header_type"`
	ContentExample string             `json:"content_example"`
	Content        string             `json:"content"`
	Type           string             `json:"type"`
}
