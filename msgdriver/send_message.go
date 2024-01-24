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
	Origin      string                            `json:"origin,omitempty"`
	IsSecret    bool                              `json:"is_secret,omitempty"`
	AgentID     string                            `json:"agent_id,omitempty"`
	AgentName   string                            `json:"agent_name,omitempty"`
	Context     *NewMessageContext                `json:"context,omitempty"`
}

type NewMessageContext struct {
	MessageID string `json:"message_id,omitempty"`
}

type SendMessageResult struct {
	MessageID string `json:"message_id"`
}
