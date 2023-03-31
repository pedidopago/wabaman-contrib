package rest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/pedidopago/go-common/mariadb"
	"github.com/pedidopago/go-common/util"
	"github.com/pedidopago/wabaman-contrib/fbgraph"
	"github.com/pedidopago/wabaman-contrib/wsapi"
)

type MessageType string

const (
	MessageText        MessageType = "text"
	MessageTemplate    MessageType = "template"
	MessageInteractive MessageType = "interactive"
	MessageImage       MessageType = "image"
	MessageVideo       MessageType = "video"
	MessageAudio       MessageType = "audio"
	MessageDocument    MessageType = "document"
	MessageSticker     MessageType = "sticker"
)

type NewMessageStatus string

// send messages statuses
const (
	NewMessageStatusImmediate                  NewMessageStatus = "sent_immediately"
	NewMessageStatusQueuedForTemplate          NewMessageStatus = "queued_for_template"
	NewMessageStatusBlockedByMarketingDisabled NewMessageStatus = "blocked_by_marketing_disabled"
	NewMessageStatusUnknown                    NewMessageStatus = "unknown"
)

type NewMessageRequest struct {
	PhoneID                uint                              `json:"phone_id"`
	BranchID               string                            `json:"branch_id"`
	FromNumber             string                            `json:"from_number"`
	ToNumber               string                            `json:"to_number"`
	Type                   MessageType                       `json:"type"`
	Text                   *fbgraph.TextObject               `json:"text,omitempty"`
	Template               *fbgraph.TemplateObject           `json:"template,omitempty"`
	Interactive            *fbgraph.InteractiveMessageObject `json:"interactive,omitempty"`
	Image                  *fbgraph.MediaObject              `json:"image,omitempty"`
	Audio                  *fbgraph.MediaObject              `json:"audio,omitempty"`
	Document               *fbgraph.MediaObject              `json:"document,omitempty"`
	Video                  *fbgraph.MediaObject              `json:"video,omitempty"`
	Sticker                *fbgraph.MediaObject              `json:"sticker,omitempty"`
	FallbackTemplate       string                            `json:"fallback_template,omitempty"`
	SkipWelcome            bool                              `json:"skip_welcome,omitempty"`
	Origin                 string                            `json:"origin,omitempty"`
	ReadMessages           bool                              `json:"read_messages,omitempty"`
	Verbose                bool                              `json:"verbose,omitempty"`
	ContactMetadata        map[string]any                    `json:"contact_metadata,omitempty"`
	OneShotContactMetadata map[string]any                    `json:"one_shot_contact_metadata,omitempty"`
	AgentID                string                            `json:"agent_id,omitempty"`
	AgentName              string                            `json:"agent_name,omitempty"`
	Context                *NewMessageContext                `json:"context,omitempty"`
}

type NewMessageContext struct {
	MessageID string `json:"message_id,omitempty"`
}

type NewMessageResponse struct {
	ID                uint64           `json:"id,omitempty"`
	MessageID         string           `json:"message_id"`
	SendMessageStatus NewMessageStatus `json:"send_message_status"`
}

type GetMessagesRequest struct {
	PhoneID           uint   `url:"phone_id,omitempty"`
	BranchID          string `url:"branch_id,omitempty"`
	HostPhoneNumber   string `url:"host_phone_number,omitempty"`
	ClientPhoneNumber string `url:"client_phone_number,omitempty"`
	MaxResults        uint64 `url:"max_results,omitempty"`
	Page              uint   `url:"page,omitempty"`
}

type GetMessagesResponse struct {
	Messages   []json.RawMessage `json:"messages,omitempty"`
	MaxResults uint64            `json:"max_results"`
	Page       uint              `json:"page"`
	LastPage   uint              `json:"last_page,omitempty"`
}

type AnyMessage interface {
	GetID() uint64
	GetCreatedAtNano() int64
	GetObjectType() string
}

type NewMessageRequestForRedisQueue struct {
	NewMessageRequest    `json:",inline"`
	AccountID            uint   `json:"account_id"`
	CallbackRedisChannel string `json:"callback_redis_channel"`
}

type NewMessageRequestForRedisQueueResponse struct {
	Status       int                 `json:"status"`
	ErrorMessage string              `json:"error_message"`
	Response     *NewMessageResponse `json:"response"`
}

func (r NewMessageRequestForRedisQueueResponse) JSONString() string {
	d, _ := json.Marshal(r)
	if d == nil {
		return ""
	}
	return string(d)
}

type NewMediaResponse struct {
	MediaID string `json:"media_id"`
}

type UpdateContactRequest struct {
	CustomerID      string         `json:"customer_id,omitempty"`
	CustomerName    string         `json:"customer_name,omitempty"`
	WABAProfileName string         `json:"waba_profile_name,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
	OneShotMetadata map[string]any `json:"one_shot_metadata,omitempty"`
	Origin          string         `json:"origin,omitempty"`
}

type UpdateContactResponse struct {
	Contact *Contact `json:"contact"`
}

type GetContactsRequest struct {
	BusinessID      uint     `query:"business_id"`
	BranchID        string   `query:"branch_id"`
	PhoneID         uint     `query:"phone_id"`
	CustomerIDs     []string `query:"customer_id"`
	WABAContactIDs  []string `query:"waba_contact_id"`
	HostPhoneNumber string   `query:"host_phone_number"`
	MaxResults      uint64   `query:"max_results"`
	Page            uint     `query:"page"`
	FetchMessages   bool     `query:"fetch_messages"`
	FetchLastPage   bool     `query:"fetch_last_page"`
	Origin          string   `query:"origin"`
}

func (req GetContactsRequest) BuildQuery() url.Values {
	q := make(url.Values)
	if iszero, _ := util.IsZero(req.BusinessID); !iszero {
		q.Set("business_id", fmt.Sprintf("%d", req.BusinessID))
	}
	if iszero, _ := util.IsZero(req.BranchID); !iszero {
		q.Set("branch_id", req.BranchID)
	}
	if iszero, _ := util.IsZero(req.PhoneID); !iszero {
		q.Set("phone_id", fmt.Sprintf("%d", req.PhoneID))
	}
	if iszero, _ := util.IsZero(req.CustomerIDs); !iszero {
		for _, id := range req.CustomerIDs {
			q.Add("customer_id", id)
		}
	}
	if iszero, _ := util.IsZero(req.WABAContactIDs); !iszero {
		for _, id := range req.WABAContactIDs {
			q.Add("waba_contact_id", id)
		}
	}
	if iszero, _ := util.IsZero(req.HostPhoneNumber); !iszero {
		q.Set("host_phone_number", req.HostPhoneNumber)
	}
	if iszero, _ := util.IsZero(req.MaxResults); !iszero {
		q.Set("max_results", fmt.Sprintf("%d", req.MaxResults))
	}
	if iszero, _ := util.IsZero(req.Page); !iszero {
		q.Set("page", fmt.Sprintf("%d", req.Page))
	}
	if iszero, _ := util.IsZero(req.FetchMessages); !iszero {
		q.Set("fetch_messages", fmt.Sprintf("%t", req.FetchMessages))
	}
	if iszero, _ := util.IsZero(req.FetchLastPage); !iszero {
		q.Set("fetch_last_page", fmt.Sprintf("%t", req.FetchLastPage))
	}
	if iszero, _ := util.IsZero(req.Origin); !iszero {
		q.Set("origin", req.Origin)
	}
	return q
}

type GetContactsResponse struct {
	Contacts   []*Contact `json:"contacts"`
	MaxResults uint64     `json:"max_results"`
	Page       uint       `json:"page"`
	LastPage   uint       `json:"last_page,omitempty"`
}

type GetContactsV2Request struct {
	BusinessID               uint     `url:"business_id,omitempty"`
	StoreID                  string   `url:"store_id,omitempty"`
	BranchID                 string   `url:"branch_id,omitempty"`
	PhoneID                  uint     `url:"phone_id,omitempty"`
	ContactIDs               []uint64 `url:"contact_id,omitempty"`
	CustomerIDs              []string `url:"customer_id,omitempty"`
	WABAContactIDs           []string `url:"waba_contact_id,omitempty"`
	ExactWABAContactIDs      bool     `url:"exact_waba_contact_ids,omitempty"`
	ExactNames               []string `url:"exact_name,omitempty"`
	Name                     string   `url:"name,omitempty"`
	HostPhoneNumber          string   `url:"host_phone_number,omitempty"`
	Tags                     []string `url:"tag,omitempty"`
	LastMessagePreviewStatus string   `url:"last_message_preview_status,omitempty"`
	MaxResults               uint64   `url:"max_results,omitempty"`
	Page                     uint     `url:"page,omitempty"`
	Origin                   string   `url:"origin,omitempty"`

	Fixed          bool `url:"fixed,omitempty"`
	UnreadMessages bool `url:"unread_messages,omitempty"`

	// metadata items

	MetaInquiryStatus     string `url:"md_inquiry_status,omitempty"`
	MetaSellerName        string `url:"md_seller_name,omitempty"`
	MetaActiveChatbot     *bool  `url:"md_active_chatbot,omitempty"`
	MetaLastCouponOffered string `url:"md_last_coupon_offered,omitempty"`
	MetaCPF               string `url:"md_cpf,omitempty"`

	// ranges

	LastMessageReceivedFrom time.Time `url:"last_message_received_from,omitempty"`
	LastMessageReceivedTo   time.Time `url:"last_message_received_to,omitempty"`
	LastMessageSentFrom     time.Time `url:"last_message_sent_from,omitempty"`
	LastMessageSentTo       time.Time `url:"last_message_sent_to,omitempty"`
	LastMessageFrom         time.Time `url:"last_message_from,omitempty"`
	LastMessageTo           time.Time `url:"last_message_to,omitempty"`

	// expert options

	// If true, no sorting is applied. Do not use this unless you know what you are doing.
	NoSorting bool `url:"no_sorting,omitempty"`
	// If true, no cache is used. Do not use this unless you know what you are doing.
	NoCache bool `url:"no_cache,omitempty"`
}

type GetContactsV2Response struct {
	Contacts   []*ContactV2 `json:"contacts"`
	MaxResults uint64       `json:"max_results"`
	Page       uint         `json:"page"`
	LastPage   uint         `json:"last_page,omitempty"`
}

type CheckIntegrationRequest struct {
	BranchID           string
	ContactPhoneNumber string
}

type CheckIntegrationResponse struct {
	BusinessID uint     `json:"business_id"`
	PhoneIDs   []uint   `json:"phone_ids,omitempty"`
	ContactIDs []uint64 `json:"contact_ids,omitempty"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
	Raw        string `json:"raw,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if e.StatusCode > 0 && e.Message != "" {
		return fmt.Sprintf("status: %d message: %s", e.StatusCode, e.Message)
	}
	if e.Message != "" {
		return e.Message
	}
	if e.StatusCode != 0 && e.Raw != "" {
		return fmt.Sprintf("status: %d raw: %s", e.StatusCode, e.Raw)
	}
	if e.Raw != "" {
		return "raw response: " + e.Raw
	}
	return "unknown error"
}

type Contact struct {
	ID uint64 `json:"id"`
	// The id of the phone object (table phone -> phone.id)
	PhoneID uint `json:"phone_id"`
	// The whatsapp id (phone number) of the contact.
	WABAContactID string `json:"waba_contact_id"`
	// The profile name of the contact.
	WABAProfileName mariadb.NullString `json:"waba_profile_name"`
	// The timestamp of the last time the contact was 'seen' online.
	LastActivity mariadb.NullTime `json:"last_activity"`
	// The ID of the last message received from the contact.
	LastMessageReceivedID mariadb.NullInt64 `json:"last_message_received_id"`
	// The ID of the last message sent to the contact.
	LastMessageSentID mariadb.NullInt64 `json:"last_message_sent_id"`
	// The timestamp of the last message received from the contact.
	LastMessageReceivedTimestamp mariadb.NullTime `json:"last_message_received_timestamp"`
	// The timestamp of the last message sent to the contact.
	LastMessageSentTimestamp mariadb.NullTime `json:"last_message_sent_timestamp"`
	// The short version of the last message sent/received
	LastMessagePreview mariadb.NullString `json:"last_message_preview"`
	// The mariadb enum if the message was sent from the contact to the host or viceversa (host|client)
	LastMessagePreviewOrigin     mariadb.NullString `json:"last_message_preview_origin"`
	LastMessagePreviewStatus     string             `json:"last_message_preview_status"`
	LastMessagePreviewWhatsAppID mariadb.NullString `json:"last_message_preview_whatsapp_id,omitempty"`
	// Contact Metadata
	Metadata map[string]any `json:"metadata"`
	// The customer_id of ms_customer
	CustomerID mariadb.NullString `json:"customer_id"`
	// The customer_name of ms_customer
	CustomerName mariadb.NullString `json:"customer_name"`
	// The datetime this contact was created on the database.
	CreatedAt time.Time `json:"created_at"`
	// The datetime this contact was last updated on the database.
	UpdatedAt      time.Time `json:"updated_at"`
	UnreadMessages *int      `json:"unread_messages,omitempty"`
	Name           string    `json:"name,omitempty"`
}

type ContactV2 struct {
	*Contact        `json:",inline"`
	AccountID       uint      `json:"account_id"`
	BusinessID      uint      `json:"business_id"`
	StoreID         string    `json:"store_id"`
	BranchID        string    `json:"branch_id"`
	HostPhoneNumber string    `json:"host_phone_number"`
	Tags            []string  `json:"tags"`
	AgentTags       []string  `json:"agent_tags"`
	LastMessage     time.Time `json:"last_message_timestamp,omitempty"`
	UnreadMessages  int       `json:"unread_messages"`
}

type GetBusinessesRequest struct {
	ID      uint   `json:"id,omitempty"`
	StoreID string `json:"store_id,omitempty"`
}

type GetBusinessesResponse struct {
	Businesses []*Business `json:"businesses"`
}

type GetPhonesRequest struct {
	ID          uint   `json:"id,omitempty"`
	BusinessID  uint   `json:"business_id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	BranchID    string `json:"branch_id,omitempty"`

	WithStatistics bool `json:"with_statistics,omitempty"`
}

type GetPhonesResponse struct {
	Phone []*Phone `json:"phone"`
}

type NewBusinessRequest struct {
	StoreID           string `json:"store_id" validate:"required"`
	StoreName         string `json:"store_name" validate:"required"`
	AccessToken       string `json:"access_token"`
	FacebookAppID     string `json:"facebook_app_id"`
	FacebookAppSecret string `json:"facebook_app_secret"`
	// Phones            []??? `json:"phones,omitempty"`
}

type NewPhoneRequest struct {
	BusinessID                uint              `json:"business_id" validate:"required"`
	WhatsAppID                string            `json:"whatsapp_id" validate:"required"`
	WhatsAppBusinessAccountID string            `json:"whatsapp_business_account_id" validate:"required"`
	TemplateNamespace         string            `json:"template_namespace"`
	PhoneNumber               string            `json:"phone_number" validate:"required"`
	BranchID                  string            `json:"branch_id" validate:"required"`
	AccessToken               util.SecretString `json:"access_token"`
	DefaultTplHeaderImage     string            `json:"default_tpl_header_image"`
	DefaultTplHeaderVideo     string            `json:"default_tpl_header_video"`
	DefaultReheatTemplateName string            `json:"default_reheat_template_name"`
	FbAppID                   string            `json:"fb_app_id"`
	FbAppSecret               string            `json:"fb_app_secret"`
	AlertEmail                string            `json:"alert_email"`
	AlertDiscord              string            `json:"alert_discord"`
}

type Business struct {
	ID                uint      `json:"id"`
	AccountID         uint      `json:"account_id"`
	StoreID           string    `json:"store_id"`
	StoreName         string    `json:"store_name"`
	FBAppID           string    `json:"fb_app_id"`
	APIKey            string    `json:"api_key"`
	UseTemplate24Rule bool      `json:"use_template_24_rule"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Phones            []*Phone  `json:"phones,omitempty"`
}

type Phone struct {
	ID                          uint      `json:"id"`
	BusinessID                  uint      `json:"business_id"`
	WhatsAppID                  string    `json:"whatsapp_id"`
	WhatsAppBusinessAccountID   string    `json:"whatsapp_business_account_id"`
	PhoneNumber                 string    `json:"phone_number"`
	BranchID                    string    `json:"branch_id"`
	BranchName                  string    `json:"branch_name,omitempty"`
	TemplateNamespace           string    `json:"template_namespace,omitempty"`
	DefaultTplHeaderImage       string    `json:"default_tpl_header_image,omitempty"`
	DefaultTplHeaderVideo       string    `json:"default_tpl_header_video,omitempty"`
	DefaultReheatTemplate       string    `json:"default_reheat_template,omitempty"`
	FBAppID                     string    `json:"fb_app_id"`
	AlertEmail                  string    `json:"alert_email,omitempty"`
	AlertDiscord                string    `json:"alert_discord,omitempty"`
	TemplateDefaultCompanyName  string    `json:"template_default_company_name,omitempty"`
	TemplateDefaultStoreURL     string    `json:"template_default_store_url,omitempty"`
	TemplateDefaultContactPhone string    `json:"template_default_contact_phone,omitempty"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`

	Statistics *PhoneStatistics `json:"statistics,omitempty"`
}

type PhoneStatistics struct {
	TotalContacts int `json:"total_contacts"`
}

type UpdateContactOptions struct {
	WABAContactID string
	BranchID      string
	Silent        bool
	Async         bool
}

type UpdateContactOption func(*UpdateContactOptions)

func UCWithWABAContactID(id string) UpdateContactOption {
	return func(o *UpdateContactOptions) {
		o.WABAContactID = id
	}
}

func UCWithBranchID(id string) UpdateContactOption {
	return func(o *UpdateContactOptions) {
		o.BranchID = id
	}
}

func UCWithSilent(silent bool) UpdateContactOption {
	return func(o *UpdateContactOptions) {
		o.Silent = silent
	}
}

func UCWithAsync(async bool) UpdateContactOption {
	return func(o *UpdateContactOptions) {
		o.Async = async
	}
}

type NewNoteRequest struct {
	ContactID     uint64                 `json:"contact_id,omitempty"`
	WABAContactID string                 `json:"waba_contact_id,omitempty"`
	BranchID      string                 `json:"branch_id,omitempty"`
	PhoneID       uint                   `json:"phone_id,omitempty"`
	AgentName     string                 `json:"agent_name,omitempty"`
	Text          string                 `json:"text"`
	Origin        string                 `json:"origin,omitempty"`
	Format        wsapi.HostNoteFormat   `json:"format,omitempty"`
	Type          string                 `json:"type"`
	Buttons       []wsapi.HostNoteButton `json:"buttons,omitempty"`
	Title         string                 `json:"title,omitempty"`
	TitleIcon     string                 `json:"title_icon,omitempty"`
	Description   string                 `json:"description,omitempty"`
	Metadata      map[string]any         `json:"metadata,omitempty"`
}

type NewNoteResponse wsapi.HostNote

type SentMessage struct {
	ID                 uint64             `json:"id"`
	PhoneID            uint               `json:"phone_id"`
	WabaMessageID      string             `json:"waba_message_id"`
	WabaRecipientID    mariadb.NullString `json:"waba_recipient_id"`
	WabaProfileName    mariadb.NullString `json:"waba_profile_name"`
	WabaTimestamp      time.Time          `json:"waba_timestamp"`
	LastStatusName     mariadb.NullString `json:"last_status_name"`
	TsStatusSent       mariadb.NullTime   `json:"ts_status_sent"`
	TsStatusDelivered  mariadb.NullTime   `json:"ts_status_delivered"`
	TsStatusRead       mariadb.NullTime   `json:"ts_status_read"`
	TsStatusFailed     mariadb.NullTime   `json:"ts_status_failed"`
	PricingBillable    bool               `json:"pricing_billable"`
	PricingModel       mariadb.NullString `json:"pricing_model"`
	PricingCategory    mariadb.NullString `json:"pricing_category"`
	WabaConversationID mariadb.NullString `json:"waba_conversation_id"`
	Type               mariadb.NullString `json:"type"`
	TextBody           mariadb.NullString `json:"text_body"`
	MediaCaption       mariadb.NullString `json:"media_caption"`
	MediaMimeType      mariadb.NullString `json:"media_mime_type"`
	MediaID            string             `json:"media_id"`
	DocumentFilename   mariadb.NullString `json:"document_filename"`
	S3FilePublicURL    mariadb.NullString `json:"s3_file_public_url"`
	S3FileKey          mariadb.NullString `json:"s3_file_key"`
	S3BucketName       mariadb.NullString `json:"s3_bucket_name"`
	TemplateName       mariadb.NullString `json:"template_name"`
	TemplateLangCode   mariadb.NullString `json:"template_lang_code"`
	AgentID            mariadb.NullString `json:"agent_id"`
	AgentName          mariadb.NullString `json:"agent_name"`
	Origin             mariadb.NullString `json:"origin"`
	CreatedAt          time.Time          `json:"created_at"`
	CreatedAtNano      int64              `json:"created_at_nano"`
	ObjectType         string             `json:"object_type"`

	Interactive  *fbgraph.InteractiveMessageObject `json:"interactive"`
	Template     *TemplateCopy                     `json:"template,omitempty"`
	FailedReason *wsapi.SentMessageFailedReason    `json:"failed_reason,omitempty"`
}

type TemplateCopy struct {
	ParsedTemplate *wsapi.ParsedTemplate `json:"parsed,omitempty"`
}

func (m *SentMessage) GetID() uint64 {
	return m.ID
}

func (m *SentMessage) GetCreatedAtNano() int64 {
	return m.CreatedAtNano
}

func (m *SentMessage) GetObjectType() string {
	return m.ObjectType
}

type ReceivedMessage struct {
	ID                     uint64             `json:"id"`
	PhoneID                uint               `json:"phone_id"`
	WABAMessageID          string             `json:"waba_message_id"`
	WABAFromID             string             `json:"waba_from_id"`
	WABAProfileName        string             `json:"waba_profile_name"`
	WABATimestamp          time.Time          `json:"waba_timestamp"`
	Type                   MessageType        `json:"type"`
	TextBody               *string            `json:"text_body"`
	MediaCaption           *string            `json:"media_caption"`
	MediaMimeType          *string            `json:"media_mime_type"`
	MediaSha256B64         *string            `json:"media_sha256_b64"`
	MediaID                *string            `json:"media_id"`
	DocumentFilename       *string            `json:"document_filename"`
	InteractiveType        *string            `json:"interactive_type"`
	InteractiveID          *string            `json:"interactive_id"`
	InteractiveTitle       *string            `json:"interactive_title"`
	InteractiveDescription *string            `json:"interactive_description"`
	ButtonPayload          *string            `json:"button_payload"`
	ButtonText             *string            `json:"button_text"`
	S3FilePublicURL        *string            `json:"s3_file_public_url"`
	S3FileKey              *string            `json:"s3_file_key"`
	S3BucketName           *string            `json:"s3_bucket_name"`
	CreatedAt              time.Time          `json:"created_at"`
	CreatedAtNano          int64              `json:"created_at_nano"`
	ReadAt                 mariadb.NullTime   `json:"read_at"`
	ReadAtMetadata         mariadb.NullString `json:"read_at_metadata"`
	ObjectType             string             `json:"object_type"`
}

func (m *ReceivedMessage) GetID() uint64 {
	return m.ID
}

func (m *ReceivedMessage) GetCreatedAtNano() int64 {
	return m.CreatedAtNano
}

func (m *ReceivedMessage) GetObjectType() string {
	return m.ObjectType
}
