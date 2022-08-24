package wsapi

// MessageType is the identifier of the Message payload.
type MessageType uint8

// All valid message types
const (
	MessageTypePing          MessageType = 0
	MessageTypePong          MessageType = 1
	MessageTypeClientMessage MessageType = 2
	MessageTypeCloseError    MessageType = 240
)

// Message is the root level object of a ms-wabaman websocket.
type Message struct {
	Type          MessageType    `json:"type"`
	Error         *Error         `json:"error,omitempty"`
	ClientMessage *ClientMessage `json:"client_message,omitempty"`
	Metadata      *Metadata      `json:"metadata,omitempty"`
}
