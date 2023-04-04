package wsapi

import "time"

// ClientMessage is a message that was sent from a client to a WhatsApp business.
type ClientMessage struct {
	// Internal ID (bigint)
	ID uint64 `json:"id"`
	// Base64 whatsapp message id from graph API.
	// It is guaranteed to be unique.
	WABAMessageID string `json:"waba_message_id"`
	// The id (phone number) of the sender.
	WABAFromID string `json:"waba_from_id"`
	// The profile name of the sender.
	WABAProfileName string    `json:"waba_profile_name"`
	WABATimestamp   time.Time `json:"waba_timestamp"`
	Type            string    `json:"type"`

	Text        *Text        `json:"text,omitempty"`
	Document    *Document    `json:"document,omitempty"`
	Video       *Video       `json:"video,omitempty"`
	Image       *Image       `json:"image,omitempty"`
	Audio       *Audio       `json:"audio,omitempty"`
	Sticker     *Sticker     `json:"sticker,omitempty"`
	Interactive *Interactive `json:"interactive,omitempty"`
	Button      *Button      `json:"button,omitempty"`
	Preview     string       `json:"preview,omitempty"`

	Context   *MessageContext   `json:"context,omitempty"`
	Reactions []MessageReaction `json:"reactions,omitempty"`

	CreatedAt     time.Time  `json:"created_at"`
	CreatedAtNano int64      `json:"created_at_nano,omitempty"`
	ReadAt        *time.Time `json:"read_at,omitempty"`
	ObjectType    string     `json:"object_type,omitempty"`
}

// Text is present in a message if type=text
type Text struct {
	Body string `json:"body"`
}

// Document is present in a message if type=document
type Document struct {
	ID        string `json:"id"` // whatsapp ID
	MimeType  string `json:"mime_type"`
	Sha256    string `json:"sha256"`
	Caption   string `json:"caption"`
	Filename  string `json:"filename"`
	PublicURL string `json:"public_url"`
}

// Audio is present in a message if type=audio
type Audio struct {
	ID        string `json:"id"` // whatsapp ID
	MimeType  string `json:"mime_type"`
	PublicURL string `json:"public_url"`
}

// Video is present in a message if type=video
type Video struct {
	ID        string `json:"id"` // whatsapp ID
	MimeType  string `json:"mime_type"`
	Sha256    string `json:"sha256"`
	Caption   string `json:"caption"`
	PublicURL string `json:"public_url"`
}

// Image is present in a message if type=image
type Image struct {
	ID        string `json:"id"` // whatsapp ID
	MimeType  string `json:"mime_type"`
	Sha256    string `json:"sha256"`
	Caption   string `json:"caption"`
	PublicURL string `json:"public_url"`
}

// Sticker is present in a message if type=sticker
type Sticker struct {
	ID        string `json:"id"` // whatsapp ID
	MimeType  string `json:"mime_type"`
	Sha256    string `json:"sha256"`
	PublicURL string `json:"public_url"`
}

type InteractiveType string

const (
	InteractiveButtonReply InteractiveType = "button_reply"
	InteractiveListReply   InteractiveType = "list_reply"
)

// Interactive is present in a message if type=interactive
type Interactive struct {
	Type        InteractiveType `json:"type"`
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"` // if type=list_reply
}

// Button is present in a message if type=button
type Button struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}
