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
	MessageTypeHostNoteUpdated    MessageType = 9
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
	HostNoteUpdated   *HostNoteUpdated   `json:"host_note_updated,omitempty"`
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
	MessageID     uint64    `json:"message_id"`
	ContactID     uint64    `json:"contact_id,omitempty"`
	WABAContactID string    `json:"waba_contact_id,omitempty"`
	ReadAt        time.Time `json:"read_at"`
	Metadata      string    `json:"metadata"`
}

type ClientReceipt struct {
	Type            string                   `json:"type"` // sent, read, delivered, failed
	MessageID       uint64                   `json:"message_id"`
	WABAContactID   string                   `json:"waba_contact_id"`
	WABAProfileName string                   `json:"waba_profile_name"`
	UpdatedAt       time.Time                `json:"updated_at"`
	FailedReason    *SentMessageFailedReason `json:"failed_reason,omitempty"`
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

type HostNoteFormat string

const (
	HostNoteFormatText    HostNoteFormat = "TEXT"
	HostNoteFormatButtons HostNoteFormat = "BUTTONS"
)

type HostNote struct {
	ID            uint64           `json:"id"`
	ContactID     uint64           `json:"contact_id"`
	WABAContactID string           `json:"waba_contact_id"`
	PhoneID       uint             `json:"phone_id"`
	Format        HostNoteFormat   `json:"format"`
	Title         string           `json:"title"`
	TitleIcon     string           `json:"title_icon"`
	Description   string           `json:"description"`
	Text          string           `json:"text"`
	Buttons       []HostNoteButton `json:"buttons,omitempty"`
	AgentID       string           `json:"agent_id"`
	AgentName     string           `json:"agent_name,omitempty"`
	Origin        string           `json:"origin,omitempty"`
	Type          string           `json:"type"`
	CreatedAt     time.Time        `json:"created_at"`
	CreatedAtNano int64            `json:"created_at_nano"`
	Metadata      map[string]any   `json:"metadata,omitempty"`
	ObjectType    string           `json:"object_type,omitempty"`
}

func (m *HostNote) GetID() uint64 {
	return m.ID
}

func (m *HostNote) GetCreatedAtNano() int64 {
	return m.CreatedAtNano
}

func (m *HostNote) GetObjectType() string {
	if m.ObjectType == "" {
		return "host_note"
	}
	return m.ObjectType
}

type HostNoteButton struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	Selected   bool   `json:"selected"`
	SelectedBy struct {
		AgentID   string `json:"agent_id,omitempty"`
		AgentName string `json:"agent_name,omitempty"`
	} `json:"selected_by,omitempty"`
	SelectedAt time.Time `json:"selected_at,omitempty"`
}

type HostNoteUpdated struct {
	HostNote         *HostNote `json:"host_note"`
	SelectedButtonID string    `json:"selected_button_id,omitempty"`
	AgentID          string    `json:"agent_id,omitempty"`
	AgentName        string    `json:"agent_name,omitempty"`
}
