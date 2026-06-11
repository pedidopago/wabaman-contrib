package wsapi

import (
	"encoding/json"

	types "github.com/pedidopago/wabaman-contrib/shared-types"

	"github.com/pedidopago/go-common/mariadb"
)

// Metadata holds contextual information about the account, business, phone, and contact
// associated with a websocket message. It is sent by the server when a client connects.
type Metadata struct {
	Account  *AccountMetadata  `json:"account,omitempty"`
	Business *BusinessMetadata `json:"business,omitempty"`
	Phone    *PhoneMetadata    `json:"phone,omitempty"`
	Contact  *ContactData      `json:"contact,omitempty"`
}

// MetadataToSend is the wire-format variant of [Metadata] with json.RawMessage fields
// for lazy deserialization. Used when sending messages to the websocket API.
type MetadataToSend struct {
	Account  *AccountMetadata       `json:"account,omitempty"`
	Business *BusinessMetadata      `json:"business,omitempty"`
	Phone    *PhoneMetadata         `json:"phone,omitempty"`
	Contact  *ContactMetadataToSend `json:"contact,omitempty"`
}

// AccountMetadata identifies the Wabaman account that owns the websocket connection.
type AccountMetadata struct {
	ID     uint   `json:"id,omitempty"`
	APIKey string `json:"api_key,omitempty"`
}

// BusinessMetadata identifies the business (store) associated with the websocket connection.
type BusinessMetadata struct {
	ID        uint   `json:"id,omitempty"`
	StoreID   string `json:"store_id,omitempty"`
	StoreName string `json:"store_name,omitempty"`
	APIKey    string `json:"api_key,omitempty"`
}

// PhoneMetadata describes the WhatsApp Business phone number (sender) associated with the connection.
type PhoneMetadata struct {
	ID                        uint            `json:"id,omitempty"`
	WhatsAppID                string          `json:"whatsapp_id,omitempty"`
	WhatsAppBusinessAccountID string          `json:"whatsapp_business_account_id,omitempty"`
	PhoneNumber               string          `json:"phone_number,omitempty"`
	BranchID                  string          `json:"branch_id,omitempty"`
	LanguageCode              string          `json:"language_code,omitempty"`
	DriverName                string          `json:"driver_name,omitempty"`
	Metadata                  json.RawMessage `json:"metadata,omitempty"`
}

// ContactMetadata is an alias kept for backward compatibility with consumers
// that reference the old name.
type ContactMetadata = ContactData

// ContactData describes a WhatsApp contact (end-user) and their conversation state.
type ContactData struct {
	ID                           uint64                `json:"id,omitzero"`
	CustomerID                   string                `json:"customer_id,omitzero"`
	CustomerName                 string                `json:"customer_name,omitzero"`
	WABAContactID                string                `json:"waba_contact_id,omitzero"`
	UserID                       string                `json:"user_id,omitzero"` // Business-scoped user ID (BSUID)
	WABAProfileName              string                `json:"waba_profile_name,omitzero"`
	Name                         string                `json:"name,omitzero"`
	ContactPhoneNumber           string                `json:"contact_phone_number,omitzero"`
	IsNewContact                 bool                  `json:"is_new_contact,omitzero"`
	Metadata                     *ContactMetadataField `json:"metadata,omitzero"`
	LastActivity                 mariadb.NullTime      `json:"last_activity,omitzero"`
	LastMessagePreview           string                `json:"last_message_preview,omitzero"`
	LastMessagePreviewOrigin     string                `json:"last_message_preview_origin,omitzero"`
	LastMessagePreviewStatus     string                `json:"last_message_preview_status,omitzero"`
	LastMessageTimestamp         mariadb.NullTime      `json:"last_message_timestamp,omitzero"`
	LastMessageReceivedTimestamp mariadb.NullTime      `json:"last_message_received_timestamp,omitzero"`
	UnreadMessages               *int                  `json:"unread_messages,omitzero"`
	ERPLastSync                  mariadb.NullTime      `json:"erp_last_sync,omitzero"`
	ColorTags                    []ColorTag            `json:"color_tags,omitzero"`
	Last24HWindow                string                `json:"last_24h_window,omitzero"`
	Last24HWindowUnix            int64                 `json:"last_24h_window_unix,omitzero"`
	MarketingEnabled             bool                  `json:"marketing_enabled,omitzero"`
}

// ContactMetadataField carries a contact metadata payload as raw JSON bytes,
// as a parsed [types.ContactMetadata], or both at once.
//
// On unmarshal it captures the raw bytes without parsing; call [ContactMetadataField.Decode]
// to lazily parse into Parsed. On marshal it emits Parsed when set, otherwise the raw bytes.
// This lets a relay forward an opaque payload (Raw only) while a producer can set a fully
// typed value (Parsed only) through the same field.
type ContactMetadataField struct {
	Raw    json.RawMessage
	Parsed *types.ContactMetadata
}

// NewContactMetadataField wraps a parsed metadata value.
func NewContactMetadataField(m *types.ContactMetadata) *ContactMetadataField {
	return &ContactMetadataField{Parsed: m}
}

// MarshalJSON emits the parsed value when present, otherwise the raw bytes.
func (f ContactMetadataField) MarshalJSON() ([]byte, error) {
	if f.Parsed != nil {
		return json.Marshal(f.Parsed)
	}
	if len(f.Raw) == 0 {
		return []byte("null"), nil
	}
	return f.Raw, nil
}

// UnmarshalJSON stores a copy of the raw bytes and clears any prior parsed value.
func (f *ContactMetadataField) UnmarshalJSON(data []byte) error {
	f.Raw = append(f.Raw[:0], data...)
	f.Parsed = nil
	return nil
}

// Decode parses the raw bytes into Parsed (once) and returns it. It returns nil
// when the field is empty or JSON null. Subsequent calls return the cached value.
func (f *ContactMetadataField) Decode() (*types.ContactMetadata, error) {
	if f == nil {
		return nil, nil
	}
	if f.Parsed != nil {
		return f.Parsed, nil
	}
	if len(f.Raw) == 0 || string(f.Raw) == "null" {
		return nil, nil
	}
	var cm types.ContactMetadata
	if err := json.Unmarshal(f.Raw, &cm); err != nil {
		return nil, err
	}
	f.Parsed = &cm
	return f.Parsed, nil
}

// IsZero reports whether the field holds neither raw bytes nor a parsed value,
// so encoding/json's omitzero drops it.
func (f *ContactMetadataField) IsZero() bool {
	return f == nil || (f.Parsed == nil && len(f.Raw) == 0)
}

// ContactMetadataToSend is the wire-format variant of [ContactMetadata] with json.RawMessage
// fields for lazy deserialization.
type ContactMetadataToSend struct {
	ID                           uint64           `json:"id,omitempty"`
	CustomerID                   string           `json:"customer_id,omitempty"`
	CustomerName                 string           `json:"customer_name,omitempty"`
	WABAContactID                string           `json:"waba_contact_id,omitempty"`
	UserID                       string           `json:"user_id,omitempty"` // Business-scoped user ID (BSUID)
	WABAProfileName              string           `json:"waba_profile_name,omitempty"`
	Name                         string           `json:"name,omitempty"`
	ContactPhoneNumber           string           `json:"contact_phone_number,omitempty"`
	IsNewContact                 bool             `json:"is_new_contact,omitempty"`
	Metadata                     json.RawMessage  `json:"metadata,omitempty"`
	LastActivity                 mariadb.NullTime `json:"last_activity,omitempty"`
	LastMessagePreview           string           `json:"last_message_preview,omitempty"`
	LastMessagePreviewOrigin     string           `json:"last_message_preview_origin,omitempty"`
	LastMessagePreviewStatus     string           `json:"last_message_preview_status,omitempty"`
	LastMessageTimestamp         mariadb.NullTime `json:"last_message_timestamp,omitempty"`
	LastMessageReceivedTimestamp mariadb.NullTime `json:"last_message_received_timestamp,omitempty"`
	UnreadMessages               *int             `json:"unread_messages,omitempty"`
	ERPLastSync                  mariadb.NullTime `json:"erp_last_sync,omitempty"`
	ColorTags                    []ColorTag       `json:"color_tags,omitempty"`
	Last24HWindow                string           `json:"last_24h_window,omitempty"`
	Last24HWindowUnix            int64            `json:"last_24h_window_unix,omitempty"`
	MarketingEnabled             bool             `json:"marketing_enabled,omitempty"`
}

// ColorTag is a simplification of CTag
type ColorTag struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
