package rest

import (
	"fmt"

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
