package rest

import "github.com/pedidopago/wabaman-contrib/fbgraph"

type MessageType string

const (
	MessageText     MessageType = "text"
	MessageTemplate MessageType = "template"
)

type NewMessageStatus string

// send messages statuses
const (
	NewMessageStatusImmediate         NewMessageStatus = "sent_immediately"
	NewMessageStatusQueuedForTemplate NewMessageStatus = "queued_for_template"
	NewMessageStatusUnknown           NewMessageStatus = "unknown"
)

type NewMessage struct {
	BranchID         string                  `json:"branch_id" validate:"required"`
	FromNumber       string                  `json:"from_number"`
	ToNumber         string                  `json:"to_number"`
	Type             MessageType             `json:"type"`
	Text             *fbgraph.TextObject     `json:"text,omitempty"`
	Template         *fbgraph.TemplateObject `json:"template,omitempty"`
	FallbackTemplate string                  `json:"fallback_template,omitempty"`
}

type NewMessageResponse struct {
	MessageID         string           `json:"message_id"`
	SendMessageStatus NewMessageStatus `json:"send_message_status"`
}
