package wsapi

// MessageType is the identifier of the Message payload.
type MessageType uint8

// All valid message types
const (
	MessageTypePing               MessageType = 0
	MessageTypePong               MessageType = 1
	MessageTypeClientMessage      MessageType = 2
	MessageTypeHostMessage        MessageType = 3
	MessageTypeMockClientMessages MessageType = 230
	MessageTypeCloseError         MessageType = 240
)

// Message is the root level object of a ms-wabaman websocket.
type Message struct {
	Type           MessageType     `json:"type"`
	Error          *Error          `json:"error,omitempty"`
	ClientMessage  *ClientMessage  `json:"client_message,omitempty"`
	HostMessage    *HostMessage    `json:"host_message,omitempty"`
	Metadata       *Metadata       `json:"metadata,omitempty"`
	ClientMockData *ClientMockData `json:"client_mock_data,omitempty"`
}

type ClientMockData struct {
	Count             int    `json:"count"`
	Interval          int    `json:"interval"`
	ClientPhoneNumber string `json:"client_phone_number"`
	ClientName        string `json:"client_name"`
	HostPhoneNumber   string `json:"host_phone_number"`
	Type              string `json:"type"`
	Text              *Text  `json:"text,omitempty"`
}
