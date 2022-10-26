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
	PhoneID         uint
	HostPhoneNumber string `json:"host_phone_number"`
	// The id (phone number) of the recipient (client).
	WABARecipientID string                            `json:"waba_recipient_id"`
	WABATimestamp   time.Time                         `json:"waba_timestamp"`
	Type            string                            `json:"type"`
	Text            *Text                             `json:"text,omitempty"`
	Document        *Document                         `json:"document,omitempty"`
	Video           *Video                            `json:"video,omitempty"`
	Image           *Image                            `json:"image,omitempty"`
	Audio           *Audio                            `json:"audio,omitempty"`
	Sticker         *Sticker                          `json:"sticker,omitempty"`
	Template        *HostTemplate                     `json:"template,omitempty"`
	Interactive     *fbgraph.InteractiveMessageObject `json:"interactive,omitempty"`
}

type HostTemplate struct {
	//Deprecated: use GraphObject.Name
	Name        string                  `json:"name"`
	GraphObject *fbgraph.TemplateObject `json:"graph_object"`
	Original    *TemplateRef            `json:"original"`
}

type TemplateRef struct {
	ID         uint64    `json:"id"`
	BusinessID uint      `json:"business_id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	CreatedAt  time.Time `json:"created_at"`

	Languages []*TemplateRefLanguage `json:"languages"`
}

type TemplateRefLanguage struct {
	LanguageCode string `json:"language_code"`

	Header      *TplRefHeader `json:"header"`
	Body        string        `json:"body"`
	Footer      string        `json:"footer"`
	ButtonsType string        `json:"buttons_type"`

	QuickReplyButtons   []TplQuickReplyButton   `json:"quick_reply_buttons"`
	CallToActionButtons []TplCallToActionButton `json:"call_to_action_buttons"`
}

type TplRefHeader struct {
	Type           string `json:"header_type"`
	ContentExample string `json:"content_example"`
}

type TplQuickReplyButton struct {
	Text string `json:"text"`
}

type TplCallToActionButton struct {
	Type string `json:"type"`
	Text string `json:"text"`
	URL  struct {
		Type string `json:"type"` // static,dynamic
		Href string `json:"href"`
	} `json:"url,omitempty"`
	Call struct {
		CC    string `json:"cc"`
		Phone string `json:"phone"`
	} `json:"call,omitempty"`
}
