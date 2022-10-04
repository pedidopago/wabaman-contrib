package rest

import (
	"fmt"
	"time"

	"github.com/pedidopago/wabaman-contrib/fbgraph"
)

type MessageType string

const (
	MessageText        MessageType = "text"
	MessageTemplate    MessageType = "template"
	MessageInteractive MessageType = "interactive"
	MessageImage       MessageType = "image"
	MessageVideo       MessageType = "video"
	MessageAudio       MessageType = "audio"
	MessageDocument    MessageType = "document"
)

type NewMessageStatus string

// send messages statuses
const (
	NewMessageStatusImmediate         NewMessageStatus = "sent_immediately"
	NewMessageStatusQueuedForTemplate NewMessageStatus = "queued_for_template"
	NewMessageStatusUnknown           NewMessageStatus = "unknown"
)

type NewMessageRequest struct {
	BranchID         string                            `json:"branch_id" validate:"required"`
	FromNumber       string                            `json:"from_number"`
	ToNumber         string                            `json:"to_number"`
	Type             MessageType                       `json:"type"`
	Text             *fbgraph.TextObject               `json:"text,omitempty"`
	Template         *fbgraph.TemplateObject           `json:"template,omitempty"`
	Interactive      *fbgraph.InteractiveMessageObject `json:"interactive,omitempty"`
	Image            *fbgraph.MediaObject              `json:"image,omitempty"`
	Audio            *fbgraph.MediaObject              `json:"audio,omitempty"`
	Document         *fbgraph.MediaObject              `json:"document,omitempty"`
	Video            *fbgraph.MediaObject              `json:"video,omitempty"`
	FallbackTemplate string                            `json:"fallback_template,omitempty"`
}

type NewMessageResponse struct {
	MessageID         string           `json:"message_id"`
	SendMessageStatus NewMessageStatus `json:"send_message_status"`
}

type NewMediaResponse struct {
	MediaID string `json:"media_id"`
}

type UpdateContactRequest struct {
	CustomerID string `json:"customer_id"`
}

type UpdateContactResponse struct {
	Contact *Contact `json:"contact"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
	Raw        string `json:"raw,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if e.StatusCode > 0 && e.Message != "" {
		return fmt.Sprintf("status: %d message: %s", e.StatusCode, e.Message)
	}
	if e.Message != "" {
		return e.Message
	}
	if e.StatusCode != 0 && e.Raw != "" {
		return fmt.Sprintf("status: %d raw: %s", e.StatusCode, e.Raw)
	}
	if e.Raw != "" {
		return "raw response: " + e.Raw
	}
	return "unknown error"
}

type Contact struct {
	ID                           uint64    `json:"id,omitempty"`                              // "id": 1,
	PhoneID                      uint      `json:"phone_id,omitempty"`                        // "phone_id": 1,
	WabaContactID                string    `json:"waba_contact_id,omitempty"`                 // "waba_contact_id": "5511941011935",
	WabaProfileName              string    `json:"waba_profile_name,omitempty"`               // "waba_profile_name": "Gabriel Ochsenhofer",
	LastActivity                 time.Time `json:"last_activity,omitempty"`                   // "last_activity": "2022-10-03T19:32:19Z",
	LastMessageReceivedID        uint64    `json:"last_message_received_id,omitempty"`        // "last_message_received_id": 145,
	LastMessageSentId            uint64    `json:"last_message_sent_id,omitempty"`            // "last_message_sent_id": 918,
	LastMessageReceivedTimestamp time.Time `json:"last_message_received_timestamp,omitempty"` // "last_message_received_timestamp": "2022-10-03T19:32:19Z",
	LastMessageSentTimestamp     time.Time `json:"last_message_sent_timestamp,omitempty"`     // "last_message_sent_timestamp": "2022-10-04T02:30:23Z",
	CustomerID                   string    `json:"customer_id,omitempty"`                     // "customer_id": "01F5E1TNWH1TCTGJ1VW71X1NA8",
	CreatedAt                    time.Time `json:"created_at,omitempty"`                      // "created_at": "2022-09-02T14:04:08Z",
	UpdatedAt                    time.Time `json:"updated_at,omitempty"`                      // "updated_at": "2022-10-04T02:30:26Z",
	//TODO: add fields below
	// HostMessages string `json:"host_messages,omitempty"` // "host_messages": null,
	// ClientMessages string `json:"client_messages,omitempty"` // "client_messages": null
}
