package wsapi

type Metadata struct {
	Account  *AccountMetadata  `json:"account,omitempty"`
	Business *BusinessMetadata `json:"business,omitempty"`
	Phone    *PhoneMetadata    `json:"phone,omitempty"`
}

type AccountMetadata struct {
	ID     uint   `json:"id,omitempty"`
	APIKey string `json:"api_key,omitempty"`
}

type BusinessMetadata struct {
	ID                        uint   `json:"id,omitempty"`
	StoreID                   string `json:"store_id,omitempty"`
	StoreName                 string `json:"store_name,omitempty"`
	WhatsAppBusinessAccountID string `json:"whatsapp_business_account_id,omitempty"`
	FacebookAppID             string `json:"facebook_app_id,omitempty"`
	APIKey                    string `json:"api_key,omitempty"`
}

type PhoneMetadata struct {
	ID          uint   `json:"id,omitempty"`
	WhatsAppID  string `json:"whatsapp_id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	BranchID    string `json:"branch_id,omitempty"`
}
