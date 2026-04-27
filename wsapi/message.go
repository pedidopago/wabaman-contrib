package wsapi

import (
	"encoding/json"
	"time"
)

// MessageType is the identifier of the Message payload.
type MessageType uint8

// All valid message types
const (
	MessageTypePing                       MessageType = 0
	MessageTypePong                       MessageType = 1
	MessageTypeClientMessage              MessageType = 2
	MessageTypeHostMessage                MessageType = 3
	MessageTypeReadByHostReceipt          MessageType = 4
	MessageTypeClientReceipt              MessageType = 5
	MessageTypeContactUpdate              MessageType = 6
	MessageTypeNewContact                 MessageType = 7
	MessageTypeHostNote                   MessageType = 8
	MessageTypeHostNoteUpdated            MessageType = 9
	MessageTypeTag                        MessageType = 10
	MessageTypeTagGroup                   MessageType = 11
	MessageTypeReaction                   MessageType = 12
	MessageTypeContactBroadcast           MessageType = 13
	MessageTypeMessageUpdated             MessageType = 14 // server sends this to the clients
	MessageTypeScheduledMessage           MessageType = 15 // server sends this to the clients
	MessageTypeCancelledScheduledMessages MessageType = 16 // server sends this to the clients
	MessageTypePresenceViewClient         MessageType = 20 // js/ts client sends this to the server
	MessageTypePresenceTypingToClient     MessageType = 21 // js/ts client sends this to the server
	MessageTypePresenceRequest            MessageType = 22 // js/ts client sends this to the server
	MessageTypePresenceResponse           MessageType = 23 // server sends this to the clients
	MessageTypeGetUnreadMessagesRequest   MessageType = 24 // js/ts client sends this to the server
	MessageTypeGetUnreadMessagesResponse  MessageType = 25 // server sends this to the clients
	MessageTypeIncomingCallFromClient     MessageType = 30 // server sends this to the clients
	MessageTypeSetupCallFromBrowser       MessageType = 31 // js/ts client sends this to the server
	MessageTypeTerminateCall              MessageType = 32 // js/ts client sends this to the server (and the server sends it to the clients)
	MessageTypeCallConsumed               MessageType = 33 // server sends this to the clients
	MessageTypeAcceptCall                 MessageType = 34 // js/ts client sends this to the server
	MessageTypeRejectCall                 MessageType = 35 // js/ts client sends this to the server
	MessageTypeSendBrowserCandidate       MessageType = 36 // BIDIRECTIONAL - js/ts client sends this to the server (and the server sends it to the clients)
	MessageTypeCallStarted                MessageType = 37 // server sends this to the clients
	MessageTypeCallEnded                  MessageType = 38 // server sends this to the clients
	MessageTypeCallOnAnswerSDP            MessageType = 39 // server sends this to the clients
	MessageTypeCallStartTimer             MessageType = 40 // server sends this to the clients
	MessageTypeReconnectCall              MessageType = 41 // js/ts client sends this to the server
	MessageTypeActiveCallNotification     MessageType = 42 // server sends this to the clients
	MessageTypeMockClientMessages         MessageType = 230
	MessageTypeGenericError               MessageType = 235
	MessageTypeCloseError                 MessageType = 240
)

// Message is the root level object of a ms-wabaman websocket.
type Message struct {
	Type                       MessageType                 `json:"type"`
	Error                      *Error                      `json:"error,omitempty"`
	ClientMessage              *ClientMessage              `json:"client_message,omitempty"`
	HostMessage                *HostMessage                `json:"host_message,omitempty"`
	ReadByHostReceipt          *ReadByHostReceipt          `json:"read_by_host_receipt,omitempty"`
	ClientReceipt              *ClientReceipt              `json:"client_receipt,omitempty"`
	ContactUpdate              *ContactUpdate              `json:"contact_update,omitempty"`
	NewContact                 *NewContact                 `json:"new_contact,omitempty"`
	HostNote                   *HostNote                   `json:"host_note,omitempty"`
	HostNoteUpdated            *HostNoteUpdated            `json:"host_note_updated,omitempty"`
	Metadata                   *Metadata                   `json:"metadata,omitempty"`
	ClientMockData             *ClientMockData             `json:"client_mock_data,omitempty"`
	Tag                        *TagEventData               `json:"tag,omitempty"`
	TagGroup                   *TagEventData               `json:"tag_group,omitempty"`
	Reaction                   *ReactionEventData          `json:"reaction,omitempty"`
	ContactBroadcast           *ContactBroadcast           `json:"contact_broadcast,omitempty"`
	MessageUpdated             *MessageUpdated             `json:"message_updated,omitempty"`
	PresenceViewClient         *PresenceViewClient         `json:"presence_view_client,omitempty"`
	PresenceTypingToClient     *PresenceTypingToClient     `json:"presence_typing_to_client,omitempty"`
	PresenceRequest            *PresenceRequest            `json:"presence_request,omitempty"`
	PresenceResponse           *PresenceResponse           `json:"presence_response,omitempty"`
	ScheduledMessage           *ScheduledMessageStub       `json:"scheduled_message,omitempty"`
	CancelledScheduledMessages *CancelledScheduledMessages `json:"cancelled_scheduled_messages,omitempty"`
	GetUnreadMessagesRequest   *GetUnreadMessagesRequest   `json:"get_unread_messages_request,omitempty"`
	GetUnreadMessagesResponse  *GetUnreadMessagesResponse  `json:"get_unread_messages_response,omitempty"`
	IncomingCallFromClient     *IncomingCallFromClient     `json:"incoming_call_from_client,omitempty"`
	SetupCallFromBrowser       *SetupCallFromBrowser       `json:"setup_call_from_browser,omitempty"`
	TerminateCall              *TerminateCall              `json:"terminate_call,omitempty"`
	CallConsumed               *CallConsumed               `json:"call_consumed,omitempty"`
	AcceptCall                 *AcceptCall                 `json:"accept_call,omitempty"`
	RejectCall                 *RejectCall                 `json:"reject_call,omitempty"`
	SendBrowserCandidate       *SendBrowserCandidate       `json:"send_browser_candidate,omitempty"`
	CallStarted                *CallStarted                `json:"call_started,omitempty"`
	CallEnded                  *CallEnded                  `json:"call_ended,omitempty"`
	CallOnAnswerSDP            *CallOnAnswerSDP            `json:"call_on_answer_sdp,omitempty"`
	CallStartTimer             *CallStartTimer             `json:"call_start_timer,omitempty"`
	ReconnectCall              *ReconnectCall              `json:"reconnect_call,omitempty"`
	ActiveCallNotification     *ActiveCallNotification     `json:"active_call_notification,omitempty"`
}

// ToJSON marshals the Message to a JSON string.
func (e Message) ToJSON() string {
	d, _ := json.Marshal(e)
	return string(d)
}

// FromJSON unmarshals a JSON string into the Message.
func (e *Message) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), e)
}

// MessageToSend is the root level object of a ms-wabaman websocket.
// This is nearly identical to Message, but with netadata fields as json.RawMessage.
// This is used to send messages to the websocket API.
type MessageToSend struct {
	Type                       MessageType                 `json:"type"`
	Error                      *Error                      `json:"error,omitempty"`
	ClientMessage              *ClientMessage              `json:"client_message,omitempty"`
	HostMessage                *HostMessage                `json:"host_message,omitempty"`
	ReadByHostReceipt          *ReadByHostReceipt          `json:"read_by_host_receipt,omitempty"`
	ClientReceipt              *ClientReceipt              `json:"client_receipt,omitempty"`
	ContactUpdate              *ContactUpdateToSend        `json:"contact_update,omitempty"`
	NewContact                 *NewContactToSend           `json:"new_contact,omitempty"`
	HostNote                   *HostNote                   `json:"host_note,omitempty"`
	HostNoteUpdated            *HostNoteUpdated            `json:"host_note_updated,omitempty"`
	Metadata                   *MetadataToSend             `json:"metadata,omitempty"`
	ClientMockData             *ClientMockData             `json:"client_mock_data,omitempty"`
	Tag                        *TagEventData               `json:"tag,omitempty"`
	TagGroup                   *TagEventData               `json:"tag_group,omitempty"`
	Reaction                   *ReactionEventData          `json:"reaction,omitempty"`
	ContactBroadcast           *ContactBroadcast           `json:"contact_broadcast,omitempty"`
	MessageUpdated             *MessageUpdated             `json:"message_updated,omitempty"`
	PresenceViewClient         *PresenceViewClient         `json:"presence_view_client,omitempty"`
	PresenceTypingToClient     *PresenceTypingToClient     `json:"presence_typing_to_client,omitempty"`
	PresenceRequest            *PresenceRequest            `json:"presence_request,omitempty"`
	PresenceResponse           *PresenceResponse           `json:"presence_response,omitempty"`
	ScheduledMessage           *ScheduledMessageStub       `json:"scheduled_message,omitempty"`
	CancelledScheduledMessages *CancelledScheduledMessages `json:"cancelled_scheduled_messages,omitempty"`
}

// ClientMockData is used to generate fake client messages for testing purposes.
type ClientMockData struct {
	Count             int    `json:"count"`
	Interval          int    `json:"interval"`
	ClientPhoneNumber string `json:"client_phone_number"`
	ClientName        string `json:"client_name"`
	HostPhoneNumber   string `json:"host_phone_number"`
	Type              string `json:"type"`
	Text              *Text  `json:"text,omitempty"`
}

// ReadByHostReceipt is sent by an agent to mark a message as read.
// The server forwards it to the WhatsApp Cloud API to send a read receipt.
type ReadByHostReceipt struct {
	MessageID     uint64    `json:"message_id"`
	WABAMessageID string    `json:"waba_message_id"`
	ContactID     uint64    `json:"contact_id,omitempty"`
	WABAContactID string    `json:"waba_contact_id,omitempty"`
	ReadAt        time.Time `json:"read_at"`
	Metadata      string    `json:"metadata"`
}

// ClientReceipt represents a delivery/read status update for a message sent to a client.
// The Type field indicates the status: "sent", "read", "delivered", or "failed".
type ClientReceipt struct {
	Type            string                   `json:"type"` // sent, read, delivered, failed
	MessageID       uint64                   `json:"message_id"`
	WABAMessageID   string                   `json:"waba_message_id"`
	WABAContactID   string                   `json:"waba_contact_id"`
	WABAProfileName string                   `json:"waba_profile_name"`
	UpdatedAt       time.Time                `json:"updated_at"`
	FailedReason    *SentMessageFailedReason `json:"failed_reason,omitempty"`
}

// ContactUpdate is broadcast to connected clients when a contact's information changes.
type ContactUpdate struct {
	ContactID          uint64         `json:"contact_id"`
	HostPhoneID        uint           `json:"host_phone_id"`
	WABAContactID      string         `json:"waba_contact_id"`
	UserID             string         `json:"user_id,omitempty"` // Business-scoped user ID (BSUID)
	WABAProfileName    string         `json:"waba_profile_name"`
	ContactPhoneNumber string         `json:"contact_phone_number,omitempty"`
	CustomerID         string         `json:"customer_id"`
	CustomerName       string         `json:"customer_name"`
	Name               string         `json:"name"`
	Metadata           map[string]any `json:"metadata"`
	ColorTags          []ColorTag     `json:"color_tags,omitempty"`
	UpdatedFields      []string       `json:"updated_fields"`
	FieldsBefore       map[string]any `json:"fields_before,omitempty"`
	FieldsAfter        map[string]any `json:"fields_after,omitempty"`
}

// ContactUpdateToSend is the wire-format variant of [ContactUpdate] with json.RawMessage
// fields for lazy deserialization.
type ContactUpdateToSend struct {
	ContactID          uint64          `json:"contact_id"`
	HostPhoneID        uint            `json:"host_phone_id"`
	WABAContactID      string          `json:"waba_contact_id"`
	UserID             string          `json:"user_id,omitempty"` // Business-scoped user ID (BSUID)
	WABAProfileName    string          `json:"waba_profile_name"`
	ContactPhoneNumber string          `json:"contact_phone_number,omitempty"`
	CustomerID         string          `json:"customer_id"`
	CustomerName       string          `json:"customer_name"`
	Name               string          `json:"name"`
	Metadata           json.RawMessage `json:"metadata"`
	ColorTags          []ColorTag      `json:"color_tags,omitempty"`
	UpdatedFields      []string        `json:"updated_fields"`
	FieldsBefore       map[string]any  `json:"fields_before,omitempty"`
	FieldsAfter        map[string]any  `json:"fields_after,omitempty"`
}

// NewContact is broadcast to connected clients when a new contact is created
// (e.g. when a previously unknown user sends their first message).
type NewContact struct {
	ContactID          uint64         `json:"contact_id"`
	HostPhoneID        uint           `json:"host_phone_id"`
	WABAContactID      string         `json:"waba_contact_id"`
	UserID             string         `json:"user_id,omitempty"` // Business-scoped user ID (BSUID)
	WABAProfileName    string         `json:"waba_profile_name"`
	ContactPhoneNumber string         `json:"contact_phone_number,omitempty"`
	CustomerID         string         `json:"customer_id"`
	Metadata           map[string]any `json:"metadata"`
}

// NewContactToSend is the wire-format variant of [NewContact] with json.RawMessage
// fields for lazy deserialization.
type NewContactToSend struct {
	ContactID          uint64          `json:"contact_id"`
	HostPhoneID        uint            `json:"host_phone_id"`
	WABAContactID      string          `json:"waba_contact_id"`
	UserID             string          `json:"user_id,omitempty"` // Business-scoped user ID (BSUID)
	WABAProfileName    string          `json:"waba_profile_name"`
	ContactPhoneNumber string          `json:"contact_phone_number,omitempty"`
	CustomerID         string          `json:"customer_id"`
	Metadata           json.RawMessage `json:"metadata"`
}

// HostNoteFormat identifies the content format of a [HostNote].
type HostNoteFormat string

const (
	HostNoteFormatText    HostNoteFormat = "TEXT"
	HostNoteFormatButtons HostNoteFormat = "BUTTONS"
	HostNoteFormatImages  HostNoteFormat = "IMAGES"
	HostNoteFormatCustom  HostNoteFormat = "CUSTOM" // Custom JSON
)

// HostNoteImage is an image attachment within a [HostNote] of format IMAGES.
type HostNoteImage struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// HostNote is an internal note created by an agent on a contact's conversation.
// Notes are not sent to the WhatsApp client; they are visible only to agents.
// The Format field determines which content field (Text, Buttons, Images, or Custom) is populated.
type HostNote struct {
	ID            uint64         `json:"id"`
	ContactID     uint64         `json:"contact_id"`
	WABAContactID string         `json:"waba_contact_id"`
	PhoneID       uint           `json:"phone_id"`
	Format        HostNoteFormat `json:"format"` // primary "type"
	Title         string         `json:"title"`
	TitleIcon     string         `json:"title_icon"`
	Description   string         `json:"description"`

	// these are mutually exclusive: (depends on the format)

	Text    string           `json:"text,omitempty"`
	Buttons []HostNoteButton `json:"buttons,omitempty"`
	Images  []HostNoteImage  `json:"images,omitempty"`
	Custom  json.RawMessage  `json:"custom,omitempty"`

	AgentID       string         `json:"agent_id"`
	AgentName     string         `json:"agent_name,omitempty"`
	Origin        string         `json:"origin,omitempty"`
	Type          string         `json:"type"` // subtype
	CreatedAt     time.Time      `json:"created_at"`
	CreatedAtNano int64          `json:"created_at_nano"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	ObjectType    string         `json:"object_type,omitempty"`
}

// GetID returns the internal ID of the host note.
func (m *HostNote) GetID() uint64 {
	return m.ID
}

// GetCreatedAtNano returns the creation timestamp in nanoseconds.
func (m *HostNote) GetCreatedAtNano() int64 {
	return m.CreatedAtNano
}

// GetObjectType returns the object type, defaulting to "host_note" when unset.
func (m *HostNote) GetObjectType() string {
	if m.ObjectType == "" {
		return "host_note"
	}
	return m.ObjectType
}

// GetOrigin returns the origin of the host note, or an empty string if nil.
func (m *HostNote) GetOrigin() string {
	if m == nil {
		return ""
	}
	return m.Origin
}

// HostNoteButton is an actionable button within a [HostNote] of format BUTTONS.
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

// HostNoteUpdated is broadcast to connected clients when a host note is modified
// (e.g. a button is selected or fields are changed).
type HostNoteUpdated struct {
	HostNote         *HostNote `json:"host_note"`
	SelectedButtonID string    `json:"selected_button_id,omitempty"`
	UpdatedFields    []string  `json:"updated_fields,omitempty"`
	AgentID          string    `json:"agent_id,omitempty"`
	AgentName        string    `json:"agent_name,omitempty"`
}

// TagEventAction represents the kind of mutation performed on a tag or tag group.
type TagEventAction string

const (
	TagEventActionCreated TagEventAction = "created"
	TagEventActionUpdated TagEventAction = "updated"
	TagEventActionDeleted TagEventAction = "deleted"
)

// TagEventData is broadcast when a tag or tag group is created, updated, or deleted.
type TagEventData struct {
	BusinessID        uint           `json:"business_id,omitempty"`
	StoreID           string         `json:"store_id,omitempty"`
	Action            TagEventAction `json:"action,omitempty"`
	PreviousName      string         `json:"previous_name,omitempty"`
	NewName           string         `json:"new_name,omitempty"`
	Name              string         `json:"name,omitempty"`
	Color             string         `json:"color,omitempty"`
	PreviousColor     string         `json:"previous_color,omitempty"`
	NewColor          string         `json:"new_color,omitempty"`
	ID                uint64         `json:"id,omitempty"`
	IsVisible         *bool          `json:"is_visible,omitempty"`
	PreviousIsVisible *bool          `json:"previous_is_visible,omitempty"`
}

// ReactionEventData is broadcast when a reaction (emoji) is added or removed from a message.
type ReactionEventData struct {
	WABAMessageID string    `json:"waba_message_id"`
	WABAContactID string    `json:"waba_contact_id"`
	Emoji         string    `json:"emoji"`
	AgentID       string    `json:"agent_id,omitempty"`
	AgentName     string    `json:"agent_name,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

// MessageContext contains information about the original message that a reply or forward refers to.
type MessageContext struct {
	MessageID           string `json:"message_id,omitempty"`
	From                string `json:"from,omitempty"`
	Forwarded           bool   `json:"forwarded,omitempty"`
	FrequentlyForwarded bool   `json:"frequently_forwarded,omitempty"`
}

// MessageReaction represents a single emoji reaction on a message.
type MessageReaction struct {
	ID            string    `json:"id,omitempty"`
	WABAContactID string    `json:"waba_contact_id"`
	Emoji         string    `json:"emoji"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	AgentID       string    `json:"agent_id,omitempty"`
	AgentName     string    `json:"agent_name,omitempty"`
	Status        string    `json:"status,omitempty"`
	Error         string    `json:"error,omitempty"`
}

// ContactBroadcast is a generic event scoped to a contact, carrying an opaque JSON payload.
type ContactBroadcast struct {
	ContactID uint64          `json:"contact_id,omitempty"`
	Type      string          `json:"type,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

// MessageUpdated is broadcast when a previously sent message is modified
// (e.g. marked as confidential).
type MessageUpdated struct {
	WABAMessageID          string `json:"waba_message_id"`
	WABAContactID          string `json:"waba_contact_id"`
	PhoneID                uint   `json:"phone_id"`
	IsFromHost             bool   `json:"is_from_host"`
	ID                     uint64 `json:"id"`
	IsMessageCConfidential bool   `json:"is_message_confidential"`
}

// ScheduledMessageStub is a lightweight representation of a scheduled message,
// broadcast to clients when a message is scheduled for future delivery.
type ScheduledMessageStub struct {
	ID            uint64 `json:"id"`
	PhoneID       uint   `json:"phone_id"`
	BranchID      string `json:"branch_id,omitempty"`
	WABAContactID string `json:"waba_contact_id"`
	ContactID     uint64 `json:"contact_id,omitempty"`
	Scheduled     struct {
		At         string `json:"at,omitempty"`
		Persistent bool   `json:"persistent,omitempty"`
	} `json:"scheduled"`
}

// CancelledScheduledMessages is broadcast when one or more scheduled messages are cancelled.
type CancelledScheduledMessages struct {
	PhoneID             uint     `json:"phone_id"`
	BranchID            string   `json:"branch_id,omitempty"`
	WABAContactID       string   `json:"waba_contact_id"`
	ContactID           uint64   `json:"contact_id,omitempty"`
	CancelledMessageIDs []uint64 `json:"cancelled_message_ids"`
}
