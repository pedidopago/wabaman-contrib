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
	Name string `json:"name"`
}
