package msgdriver

import "github.com/pedidopago/wabaman-contrib/fbgraph"

// SendMessage is an object with the necessary information for the driver to send a message to a client.
type SendMessage struct {
	// Destination is usually the phone number of the client.
	Destination string                            `json:"destination"`
	Type        MessageType                       `json:"type"`
	Text        *fbgraph.TextObject               `json:"text,omitempty"`
	Template    *fbgraph.TemplateObject           `json:"template,omitempty"`
	Interactive *fbgraph.InteractiveMessageObject `json:"interactive,omitempty"`
	Image       *fbgraph.MediaObject              `json:"image,omitempty"`
	Audio       *fbgraph.MediaObject              `json:"audio,omitempty"`
	Document    *fbgraph.MediaObject              `json:"document,omitempty"`
	Video       *fbgraph.MediaObject              `json:"video,omitempty"`
	Sticker     *fbgraph.MediaObject              `json:"sticker,omitempty"`
	Contacts    []fbgraph.ContactObject           `json:"contacts,omitempty"`
	Origin      string                            `json:"origin,omitempty"`
	IsSecret    bool                              `json:"is_secret,omitempty"`
	AgentID     string                            `json:"agent_id,omitempty"`
	AgentName   string                            `json:"agent_name,omitempty"`
	Context     *NewMessageContext                `json:"context,omitempty"`
	Metadata    map[string]any                    `json:"metadata,omitempty"`
}

func (m *SendMessage) GetOrigin() string {
	if m == nil {
		return ""
	}

	return m.Origin
}

type NewMessageContext struct {
	MessageID string `json:"message_id,omitempty"`
}

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
)

type SendMessageResult struct {
	MessageID     string         `json:"message_id"`
	MessageStatus MessageStatus  `json:"message_status,omitempty"`
	ContactID     string         `json:"contact_id"` // This is usually the contact's parsed phone number. AKA WABAContactID (due to legacy WABA implementation)
	ContactName   string         `json:"contact_name"`
	Metadata      map[string]any `json:"metadata,omitempty"` // This can be used to save metadata to this recently created message.
}
