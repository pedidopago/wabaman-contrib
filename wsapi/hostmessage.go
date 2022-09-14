package wsapi

type HostMessage struct {
	// Internal ID (bigint)
	ID uint64 `json:"id"`
	// Base64 whatsapp message id from graph API.
	// It is guaranteed to be unique.
	WABAMessageID string `json:"waba_message_id"`
	// Internal Phone ID
	PhoneID         uint
	HostPhoneNumber string `json:"host_phone_number"`
	// The id (phone number) of the recipient (client).
	WABARecipientID string        `json:"waba_recipient_id"`
	Type            string        `json:"type"`
	Text            *Text         `json:"text,omitempty"`
	Template        *HostTemplate `json:"template,omitempty"`
}

type HostTemplate struct {
	Name string `json:"name"`
}
