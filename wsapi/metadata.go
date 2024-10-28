package wsapi

import "github.com/pedidopago/go-common/mariadb"

type Metadata struct {
	Account  *AccountMetadata  `json:"account,omitempty"`
	Business *BusinessMetadata `json:"business,omitempty"`
	Phone    *PhoneMetadata    `json:"phone,omitempty"`
	Contact  *ContactMetadata  `json:"contact,omitempty"`
}

type AccountMetadata struct {
	ID     uint   `json:"id,omitempty"`
	APIKey string `json:"api_key,omitempty"`
}

type BusinessMetadata struct {
	ID        uint   `json:"id,omitempty"`
	StoreID   string `json:"store_id,omitempty"`
	StoreName string `json:"store_name,omitempty"`
	APIKey    string `json:"api_key,omitempty"`
}

type PhoneMetadata struct {
	ID                        uint   `json:"id,omitempty"`
	WhatsAppID                string `json:"whatsapp_id,omitempty"`
	WhatsAppBusinessAccountID string `json:"whatsapp_business_account_id,omitempty"`
	PhoneNumber               string `json:"phone_number,omitempty"`
	BranchID                  string `json:"branch_id,omitempty"`
	LanguageCode              string `json:"language_code,omitempty"`
}

type ContactMetadata struct {
	ID                           uint64           `json:"id,omitempty"`
	CustomerID                   string           `json:"customer_id,omitempty"`
	CustomerName                 string           `json:"customer_name,omitempty"`
	WABAContactID                string           `json:"waba_contact_id,omitempty"`
	WABAProfileName              string           `json:"waba_profile_name,omitempty"`
	Name                         string           `json:"name,omitempty"`
	ContactPhoneNumber           string           `json:"contact_phone_number,omitempty"`
	IsNewContact                 bool             `json:"is_new_contact,omitempty"`
	Metadata                     map[string]any   `json:"metadata,omitempty"`
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
