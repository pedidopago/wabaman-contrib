package wsapi

import "time"

// MessageType is the identifier of the Message payload.
type MessageType uint8

// All valid message types
const (
	MessageTypePing               MessageType = 0
	MessageTypePong               MessageType = 1
	MessageTypeClientMessage      MessageType = 2
	MessageTypeHostMessage        MessageType = 3
	MessageTypeReadByHostReceipt  MessageType = 4
	MessageTypeClientReceipt      MessageType = 5
	MessageTypeContactUpdate      MessageType = 6
	MessageTypeNewContact         MessageType = 7
	MessageTypeHostNote           MessageType = 8
	MessageTypeMockClientMessages MessageType = 230
	MessageTypeCloseError         MessageType = 240
)

// Message is the root level object of a ms-wabaman websocket.
type Message struct {
	Type              MessageType        `json:"type"`
	Error             *Error             `json:"error,omitempty"`
	ClientMessage     *ClientMessage     `json:"client_message,omitempty"`
	HostMessage       *HostMessage       `json:"host_message,omitempty"`
	ReadByHostReceipt *ReadByHostReceipt `json:"read_by_host_receipt,omitempty"`
	ClientReceipt     *ClientReceipt     `json:"client_receipt,omitempty"`
	ContactUpdate     *ContactUpdate     `json:"contact_update,omitempty"`
	NewContact        *NewContact        `json:"new_contact,omitempty"`
	HostNote          *HostNote          `json:"host_note,omitempty"`
	Metadata          *Metadata          `json:"metadata,omitempty"`
	ClientMockData    *ClientMockData    `json:"client_mock_data,omitempty"`
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

type ReadByHostReceipt struct {
	MessageID uint64    `json:"message_id"`
	ReadAt    time.Time `json:"read_at"`
	Metadata  string    `json:"metadata"`
}

type ClientReceipt struct {
	Type            string    `json:"type"` // sent, read, delivered
	MessageID       uint64    `json:"message_id"`
	WABAContactID   string    `json:"waba_contact_id"`
	WABAProfileName string    `json:"waba_profile_name"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ContactUpdate struct {
	ContactID       uint64                 `json:"contact_id"`
	HostPhoneID     uint                   `json:"host_phone_id"`
	WABAContactID   string                 `json:"waba_contact_id"`
	WABAProfileName string                 `json:"waba_profile_name"`
	CustomerID      string                 `json:"customer_id"`
	CustomerName    string                 `json:"customer_name"`
	Name            string                 `json:"name"`
	Metadata        map[string]interface{} `json:"metadata"`
	UpdatedFields   []string               `json:"updated_fields"`
}

type NewContact struct {
	ContactID       uint64                 `json:"contact_id"`
	HostPhoneID     uint                   `json:"host_phone_id"`
	WABAContactID   string                 `json:"waba_contact_id"`
	WABAProfileName string                 `json:"waba_profile_name"`
	CustomerID      string                 `json:"customer_id"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type HostNote struct {
	Text          string    `json:"text"`
	AgentID       string    `json:"agent_id"`
	AgentName     string    `json:"agent_name,omitempty"`
	Origin        string    `json:"origin,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedAtNano int64     `json:"created_at_nano"`
}
