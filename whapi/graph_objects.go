package whapi

import (
	"encoding/json"
	"fmt"

	"github.com/pedidopago/wabaman-contrib/fbgraph"
	wtypes "github.com/pedidopago/wabaman-contrib/shared-types"
	"github.com/pedidopago/wabaman-contrib/wsapi"
	"github.com/rs/zerolog/log"
)

type WebhookObjectType string

const (
	WebhookObjectWhatsappBusinessAccount WebhookObjectType = "whatsapp_business_account"
)

//TODO: coexist follow-up https://developers.facebook.com/docs/whatsapp/embedded-signup/custom-flows/onboarding-business-app-users

// Disparalha:
// https://developers.facebook.com/apps/482087044545101/whatsapp-business/wa-settings/?business_id=1033582667966734

// TODO: check this page periodically for updates:
// https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks/components

// WebhookObject - Webhooks are triggered when a customer performs an action or
// the status for a message a business sends a customer changes. You get a webhooks notification:
//
//	When a customer performs an action
//	  Sends a text message to the business
//	  Sends an image, video, audio, document, or sticker to the business
//	  Sends contact information to the business
//	  Sends location information to the business
//	  Clicks a reply button set up by the business
//	  Clicks a call-to-actions button on an Ad that Clicks to WhatsApp
//	  Clicks an item on a business' list
//	  Updates their profile information such as their phone number
//
//	When the status for a message received by a business changes (includes pricing information)
//	  delivered
//	  read
//	  sent
type WebhookObject struct {
	Object WebhookObjectType `json:"object"`
	Entry  []EntryObject     `json:"entry"`
}

type EntryObject struct {
	ID      string         `json:"id"` // WHATSAPP-BUSINESS-ACCOUNT-ID
	Changes []ChangeObject `json:"changes"`
}

type ChangeObjectField string

const (
	ChangeObjectFieldMessages         ChangeObjectField = "messages"
	ChangeObjectFieldContacts         ChangeObjectField = "contacts"
	ChangeObjectFieldErrors           ChangeObjectField = "errors"
	ChangeObjectFieldStatuses         ChangeObjectField = "statuses"
	ChangeObjectFieldHistory          ChangeObjectField = "history"
	ChangeObjectFieldSMBMessageEchoes ChangeObjectField = "smb_message_echoes"
	ChangeObjectFieldSMBAppStateSync  ChangeObjectField = "smb_app_state_sync"
	ChangeObjectFieldCalls            ChangeObjectField = "calls"
)

type ChangeObject struct {
	Field ChangeObjectField `json:"field"`
	Value *ValueObject      `json:"value"`
}

// ValueObject - The value object contains details for the change that triggered
// the webhook. This object is nested within the changes array of the entry array.
type ValueObject struct {
	// The value is whatsapp.
	MessagingProduct string `json:"messaging_product,omitempty"`
	// Array of contacts objects with information for the customer who sent a
	// message to the business.
	Contacts []ContactObject `json:"contacts,omitempty"`
	// Array of error objects with information received when a message failed.
	Errors []ErrorObject `json:"errors,omitempty"`
	// Information about a message received by the business that is subscribed
	// to the webhook.
	Messages []MessageObject `json:"messages,omitempty"`
	// Status for a message that was sent by the business that is subscribed to the webhook.
	Statuses []StatusObject `json:"statuses,omitempty"`
	// History for a message that was sent by the business that is subscribed to the webhook.
	Histories []HistoryObject `json:"history,omitempty"`
	// Message echoes for a message that was sent by the business that is subscribed to the webhook.
	MessageEchoes []MessageEHObject `json:"message_echoes,omitempty"`
	// State sync for a message that was sent by the business that is subscribed to the webhook.
	StateSync []StateSyncObject `json:"state_sync,omitempty"`
	// UserPreferences https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks/reference/user_preferences
	UserPreferences []UserPreferencesObject `json:"user_preferences,omitempty"`
	// Calls https://developers.facebook.com/docs/whatsapp/cloud-api/calling/user-initiated-calls
	Calls []CallObject `json:"calls,omitempty"`
	// Metadata for the business that is subscribed to the webhook.
	Metadata ValueObjectMetadata `json:"metadata,omitempty"`

	// template specific fields:

	Event                   string `json:"event,omitempty"`
	MessageTemplateID       uint64 `json:"message_template_id,omitempty"`
	MessageTemplateName     string `json:"message_template_name,omitempty"`
	MessageTemplateLanguage string `json:"message_template_language,omitempty"`
	Reason                  string `json:"reason,omitempty"`
	MessageTemplateCategory string `json:"message_template_category,omitempty"`
	// only included if template disabled
	DisableInfo struct {
		DisableDate string `json:"disable_date,omitempty"`
	} `json:"disable_info,omitempty"`
	// only included if template locked or unlocked
	OtherInfo struct {
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"other_info,omitempty"`
	// only included if template rejected with INVALID_FORMAT reason
	RejectionInfo struct {
		Reason         string `json:"reason,omitempty"`
		Recommendation string `json:"recommendation,omitempty"`
	} `json:"rejection_info,omitempty"`
}

func (v ValueObject) GetContactProfileName(waid string) string {
	for _, c := range v.Contacts {
		if c.WAID == waid {
			return c.Profile.Name
		}
	}
	return ""
}

// ValueObject.Event types when `value` is a template event
const (
	// Indicates the template has been approved and can now be sent in template messages.
	TemplateEventApproved = "APPROVED"
	// ndicates the template has been archived to keep the list of templates in WhatsApp manager clean.
	TemplateEventArchived = "ARCHIVED"
	// Indicates the template has been deleted.
	TemplateEventDeleted = "DELETED"
	// Indicates the template has been disabled due to user feedback.
	TemplateEventDisabled = "DISABLED"
	// Indicates the template has received negative feedback and is at risk of being disabled.
	TemplateEventFlagged = "FLAGGED"
	// Indicates the template is in the appeal process.
	TemplateEventInAppeal = "IN_APPEAL"
	// Indicates the WhatsApp Business Account template is at its template limit.
	TemplateEventLimitExceeded = "LIMIT_EXCEEDED"
	// Indicates the template has been locked and cannot be edited.
	TemplateEventLocked = "LOCKED"
	// Indicates the template has been paused.
	TemplateEventPaused = "PAUSED"
	// Indicates the template is undergoing template review.
	TemplateEventPending = "PENDING"
	// Indicates the template is no longer flagged or disabled and can be sent in template messages again.
	TemplateEventReinstated = "REINSTATED"
	// Indicates template has been deleted via WhatsApp Manager.
	TemplateEventPendingDeletion = "PENDING_DELETION"
	// Indicates the template has been rejected. You can edit the template to have it
	// undergo template review again or appeal the rejection.
	TemplateEventRejected = "REJECTED"
)

type ValueObjectMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number,omitempty"` // PHONE-NUMBER
	PhoneNumberID      string `json:"phone_number_id,omitempty"`      // PHONE-NUMBER-ID
}

// ContactObject contains information for the customer who sent a message to the business.
type ContactObject struct {
	// The customer's WhatsApp ID. A business can respond to a message using this ID.
	WAID string `json:"wa_id"`
	// An object containing customer profile information.
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
}

type MessageObjectType string

// All MessageObjectTypes
const (
	MOTypeAudio       MessageObjectType = "audio"
	MOTypeButton      MessageObjectType = "button"
	MOTypeDocument    MessageObjectType = "document"
	MOTypeText        MessageObjectType = "text"
	MOTypeImage       MessageObjectType = "image"
	MOTypeInteractive MessageObjectType = "interactive"
	MOTypeContacts    MessageObjectType = "contacts"
	MOTypeSticker     MessageObjectType = "sticker"
	MOTypeSystem      MessageObjectType = "system" // for customer number change messages
	MOTypeUnknown     MessageObjectType = "unknown"
	MOTypeVideo       MessageObjectType = "video"
	MOTypeReaction    MessageObjectType = "reaction"
	MOTypeUnsupported MessageObjectType = "unsupported"
)

func (m MessageObjectType) String() string {
	return string(m)
}

func (m MessageObjectType) IsValid() bool {
	switch m {
	case MOTypeAudio, MOTypeButton, MOTypeDocument, MOTypeText, MOTypeImage, MOTypeInteractive, MOTypeSticker, MOTypeSystem, MOTypeUnknown, MOTypeVideo, MOTypeReaction, MOTypeContacts:
		return true
	}
	return false
}

type CommonMessageObject interface {
	GetType() string
	GetID() string
	GetFrom() string
	GetTo() string
	GetTimestamp() string
	StructName() string
}

// MessageObject - The messages array of objects is nested within the value object
// and is triggered when a customer updates their profile information or a customer
// sends a message to the business that is subscribed to the webhook.
type MessageObject struct {
	// When the messages type is set to audio, including voice messages, this
	// object is included in the messages object:
	Audio *MessageObjectAudio `json:"audio,omitempty"`
	// When the messages type field is set to button, this object is included in
	// the messages object:
	Button *MessageObjectButton `json:"button,omitempty"`
	// When the messages type field is set to button, this object is included
	// in the messages object. The context for a message that was forwarded or
	// in an inbound reply from the customer.
	Context *MessageObjectContext `json:"context,omitempty"`
	// When messages type is set to document, this object is included in the
	// messages object.
	Document *MessageObjectDocument `json:"document,omitempty"`
	Errors   []ErrorObject          `json:"errors,omitempty"`
	// The customer's phone number who sent the message to the business.
	From string `json:"from"`
	// The ID for the message that was received by the business. You could use
	// messages endpoint to mark this specific message as read.
	ID string `json:"id"`
	// A webhook is triggered when a customer's phone number or profile
	// information has been updated. See messages system identity
	Identity *MessageObjectIdentity `json:"identity,omitempty"`
	// When messages type is set to image, this object is included in the messages object.
	Image *MessageObjectImage `json:"image,omitempty"`
	// When a customer selected a button or list reply, this object is included
	// in the messages object.
	Interactive *MessageObjectInteractive `json:"interactive,omitempty"`
	Contacts    []fbgraph.ContactObject   `json:"contacts,omitempty"`
	// A customer clicked an ad that redirects them to WhatsApp, this object
	// is included in the messages object.
	Referral *wtypes.MessageObjectReferral `json:"referral,omitempty"`
	// When messages type is set to sticker, this object is included in the
	// messages object.
	Sticker *MessageObjectSticker `json:"sticker,omitempty"`
	// When messages type is set to system, a customer has updated their phone
	// number or profile information, this object is included in the messages object.
	System *MessageObjectSystem `json:"system,omitempty"`
	// When messages type is set to location.
	Location *MessageObjectLocation `json:"location,omitempty"`
	// When messages type is set to text, this object is included. This object
	// includes the following field:
	Text *MessageObjectText `json:"text,omitempty"`
	// The time when the customer sent the message to the business.
	Timestamp string                 `json:"timestamp"`
	Type      MessageObjectType      `json:"type"`
	Video     *MessageObjectVideo    `json:"video,omitempty"`
	Reaction  *MessageObjectReaction `json:"reaction,omitempty"`
}

func (m MessageObject) GetType() string {
	return string(m.Type)
}

func (m MessageObject) GetID() string {
	return m.ID
}

func (m MessageObject) GetFrom() string {
	return m.From
}

func (m MessageObject) GetTo() string {
	return ""
}

func (m MessageObject) StructName() string {
	return "whapi.MessageObject"
}

func (m MessageObject) GetTimestamp() string {
	return m.Timestamp
}

type MessageObjectAudio struct {
	ID       string `json:"id"`        // ID for the audio file
	MimeType string `json:"mime_type"` // Mime type of the audio file
}

type MessageObjectButton struct {
	// The payload for a button set up by the business that a customer clicked
	// as part of an interactive message
	Payload string `json:"payload"`
	// Button text
	Text string `json:"text"`
}

type MessageObjectContext struct {
	// Set to true if the message received by the business has been forwarded
	Forwarded bool `json:"forwarded"`
	// Set to true if the message received by the business has been forwarded
	// more than 5 times.
	FrequentlyForwarded bool `json:"frequently_forwarded"`
	// The WhatsApp ID for the customer who replied to an inbound message
	From string `json:"from"`
	// The message ID for the sent message for an inbound reply
	ID string `json:"id"`
}

type MessageObjectDocument struct {
	// Caption for the document, if provided
	Caption string `json:"caption"`
	// Name for the file on the sender's device
	Filename string `json:"filename"`
	Ha256    string `json:"ha256"`
	// Mime type of the document file
	MimeType string `json:"mime_type"`
	// ID for the document
	ID string `json:"id"`
}

type MessageObjectIdentity struct {
	// State of acknowledgment for the messages system customer_identity_changed
	Acknowledged string `json:"acknowledged"`
	// The time when the WhatsApp Business Management API detected the customer may have changed their profile information
	CreatedTimestamp string `json:"created_timestamp"`
	// The ID for the messages system customer_identity_changed
	Hash string `json:"hash"`
}

type MessageObjectImage struct {
	// Caption for the image, if provided
	Caption string `json:"caption"`
	// Image hash
	Sha256 string `json:"sha256"`
	// ID for the image
	ID string `json:"id"`
	// Mime type for the image
	MimeType string `json:"mime_type"`
}

type MessageObjectInteractive struct {
	Type string `json:"type"` // button_reply or list_reply
	// Sent when a customer clicks a button
	ButtonReply *ButtonReply `json:"button_reply,omitempty"`
	// Sent when a customer selects an item from a list
	ListReply *ListReply `json:"list_reply,omitempty"`
}

type ButtonReply struct {
	// Unique ID of a button
	ID string `json:"id"`
	// Title of a button
	Title string `json:"title"`
}

type ListReply struct {
	// Unique ID of the selected list item
	ID string `json:"id"`
	// Title of the selected list item
	Title string `json:"title"`
	// Description of the selected row
	Description string `json:"description"`
}

type MessageObjectSticker struct {
	// image/webp
	MimeType string `json:"mime_type"`
	// Hash for the sticker
	Sha256 string `json:"sha256"`
	// ID for the sticker
	ID string `json:"id"`
}

type MessageObjectSystem struct {
	// Type of system update.
	Type SystemMessageType `json:"type"`
	// Describes the change to the customer's identity or phone number
	Body string `json:"body"`
	// Hash for the identity fetched from server
	Identity string `json:"identity"`
	// New WhatsApp ID for the customer when their phone number is updated. Available on webhook versions V11 and below
	NewWaId string `json:"new_wa_id"`
	// New WhatsApp ID for the customer when their phone number is updated. Available on webhook versions V12 and above
	WaId string `json:"wa_id"`
	// The WhatsApp ID for the customer prior to the update
	Customer string `json:"customer"`
}

type MessageObjectLocation struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
}

type SystemMessageType string

const (
	// A customer changed their phone number
	SysMsgTypeCustomerChangedNumber SystemMessageType = "customer_changed_number"
	// A customer changed their profile information
	SysMsgTypeCustomerIdentityChanged SystemMessageType = "customer_identity_changed"
)

type MessageObjectText struct {
	// The text of the message.
	Body string `json:"body"`
}

type MessageObjectVideo struct {
	// The caption for the video, if provided
	Caption string `json:"caption"`
	// The name for the file on the sender's device
	Filename string `json:"filename"`
	// The hash for the video
	Sha256 string `json:"sha256"`
	// The ID for the video
	ID string `json:"id"`
	// The mime type for the video file
	MimeType string `json:"mime_type"`
}

type MessageObjectReaction struct {
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

type StatusObject struct {
	// Information about the conversation.
	Conversation *StatusConversationObject `json:"conversation,omitempty"`
	// The ID for the message that the business that is subscribed to the webhooks
	// sent to a customer
	ID string `json:"id"`
	// An object containing billing information.
	Pricing *StatusPricingObject `json:"pricing,omitempty"`
	// The customer's WhatsApp ID. A business can respond to a customer using
	// this ID. This ID may not match the customer's phone number, which is
	// returned by the API as input when sending a message to the customer.
	RecipientID string                      `json:"recipient_id"`
	Status      MessageStatus               `json:"status"`
	Timestamp   string                      `json:"timestamp,omitempty"`
	Errors      []wsapi.FBStatusObjectError `json:"errors,omitempty"`
}

func (s StatusObject) JSONPrintErrors() string {
	if len(s.Errors) == 1 {
		eout, err := json.MarshalIndent(s.Errors[0], "", "  ")
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal error")
			return fmt.Sprint(s.Errors[0])
		}
		return string(eout)
	}
	eouts := struct {
		Errors []wsapi.FBStatusObjectError `json:"errors"`
	}{Errors: s.Errors}
	eout, err := json.MarshalIndent(eouts, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal errors")
		return fmt.Sprint(s.Errors)
	}
	return string(eout)
}

func (s StatusObject) JSONReasonError() *wsapi.SentMessageFailedReason {
	if len(s.Errors) == 1 {
		return &wsapi.SentMessageFailedReason{
			Code:           s.Errors[0].Code,
			Title:          s.Errors[0].Title,
			Href:           s.Errors[0].Href,
			LocalizedError: GetLocalizedError(s.Errors[0].Code),
		}
	}

	for i, e := range s.Errors {
		s.Errors[i].LocalizedError = GetLocalizedError(e.Code)
	}

	return &wsapi.SentMessageFailedReason{
		Errors: s.Errors,
	}
}

//TODO: obtain all fields from https://developers.facebook.com/docs/graph-api/webhooks/reference/whatsapp_business_account/#history

type ChunkPhase int

const (
	ChunkPhaseDayOne ChunkPhase = iota
	ChunkPhaseDayOneTo90
	ChunkPhaseDay90To180
)

type HistoryObject struct {
	Metadata struct {
		Phase      ChunkPhase `json:"phase"`
		ChunkOrder int        `json:"chunk_order"` // Indicates chunk number, which you can use to order sets of webhooks sequentially.
		Progress   int        `json:"progress"`    // min 0, max 100; Indicates percentage total of synchronization progress.
	} `json:"metadata"`
	Threads []HistoryThreadObject `json:"threads"`
	Errors  []HistoryErrorObject  `json:"errors,omitempty"`
}

type HistoryErrorObject struct {
	Code      int    `json:"code"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	ErrorData struct {
		Details string `json:"details"`
	} `json:"error_data"`
}

type HistoryThreadObject struct {
	ID       string            `json:"id"` // WHATSAPP_USER_PHONE_NUMBER - The WhatsApp user's phone number.
	Messages []MessageEHObject `json:"messages"`
}

//TODO: obtain all fields from https://developers.facebook.com/docs/graph-api/webhooks/reference/whatsapp_business_account/#smb_message_echoes

type MessageEHObject struct {
	From           string                  `json:"from"`      // BUSINESS_OR_WHATSAPP_USER_PHONE_NUMBER
	To             string                  `json:"to"`        // only included if SMB message echo
	ID             string                  `json:"id"`        // WHATSAPP_MESSAGE_ID
	Timestamp      string                  `json:"timestamp"` // DEVICE_TIMESTAMP
	Type           string                  `json:"type"`      // MESSAGE_TYPE
	Text           *fbgraph.TextObject     `json:"text,omitempty"`
	Image          *fbgraph.MediaObject    `json:"image,omitempty"`
	Audio          *fbgraph.MediaObject    `json:"audio,omitempty"`
	Document       *fbgraph.MediaObject    `json:"document,omitempty"`
	Video          *fbgraph.MediaObject    `json:"video,omitempty"`
	Sticker        *fbgraph.MediaObject    `json:"sticker,omitempty"`
	Contacts       []fbgraph.ContactObject `json:"contacts,omitempty"`
	Context        *fbgraph.MessageContext `json:"context,omitempty"`
	HistoryContext struct {
		Status string `json:"status"` // MESSAGE_STATUS - DELIVERED ERROR PENDING PLAYED READ SENT
	} `json:"history_context,omitempty"`
}

func (m MessageEHObject) GetType() string {
	return string(m.Type)
}

func (m MessageEHObject) GetID() string {
	return m.ID
}

func (m MessageEHObject) GetFrom() string {
	return m.From
}

func (m MessageEHObject) GetTo() string {
	return m.To
}

func (m MessageEHObject) GetTimestamp() string {
	return m.Timestamp
}

func (m MessageEHObject) StructName() string {
	return "whapi.MessageEHObject"
}

//TODO: obtain all fields from https://developers.facebook.com/docs/graph-api/webhooks/reference/whatsapp_business_account/#smb_app_state_sync

type UserPreferencesObject struct {
	WAID     string `json:"wa_id"`
	Detail   string `json:"detail"`
	Category string `json:"category"`
	Value    string `json:"value"`
	// Timestamp string `json:"timestamp"`
}

type StateSyncObjectAction string

const (
	StateSyncObjectActionAdd    StateSyncObjectAction = "add"
	StateSyncObjectActionRemove StateSyncObjectAction = "remove"
)

type StateSyncObject struct {
	Type    string `json:"type"`
	Contact struct {
		FullName    string `json:"full_name"`
		FirstName   string `json:"first_name"`
		PhoneNumber string `json:"phone_number"`
	} `json:"contact"`
	Action   StateSyncObjectAction `json:"action"`
	Metadata struct {
		Timestamp string  `json:"timestamp"`
		Version   float64 `json:"version"`
	} `json:"metadata"`
}

type MessageStatus string

const (
	// A webhook is triggered when a message received by a business has been delivered.
	MessageStatusDelivered MessageStatus = "delivered"
	// A webhook is triggered when a message received by a business has been read.
	MessageStatusRead MessageStatus = "read"
	// A webhook is triggered when a business receives a message from a customer.
	MessageStatusSent MessageStatus = "sent"
)

// WhatsApp defines a conversation as a 24-hour session of messaging between a
// person and a business. There is no limit on the number of messages that can
// be exchanged in the fixed 24-hour window. The 24-hour conversation session
// begins when:
//
//	A business-initiated message is delivered to a customer
//	A business’ reply to a customer message is delivered
//
// The 24-hour conversation session is different from the 24-hour customer support
// window. The customer support window is a rolling window that is refreshed when
// a customer-initiated message is delivered to a business. Within the customer
// support window businesses can send free-form messages. Any business-initiated
// message sent more than 24 hours after the last customer message must be a template message.
type StatusConversationObject struct {
	// Represents the ID of the conversation the given status notification belongs to.
	ID string `json:"id"`
	// Indicates who initiated the conversation
	Origin struct {
		Type OriginType `json:"type"`
	}
	// Date when the conversation expires. This field is only present for messages
	// with a `status` set to `sent`.
	ExpirationTimestamp string `json:"expiration_timestamp"`
}

type StatusPricingObject struct {
	// Indicates the conversation category:
	Category PricingCategory `json:"category"`
	// Type of pricing model used by the business. Current supported value is CBP
	PricingModel string `json:"pricing_model"`
}

type PricingCategory string

const (
	// Indicates an authentication conversation.
	PricingCategoryAuthentication PricingCategory = "authentication"
	// Indicates an authentication-international conversation.
	// https://developers.facebook.com/docs/whatsapp/pricing/authentication-international-rates#authentication-international-rates
	PricingCategoryAuthenticationInternational PricingCategory = "authentication-international"
	// Indicates an marketing conversation.
	PricingCategoryMarketing PricingCategory = "marketing"
	// Indicates a utility conversation.
	PricingCategoryUtility PricingCategory = "utility"
	// Indicates an service conversation.
	PricingCategoryService PricingCategory = "service"
	// Indicates a free entry point conversation.
	// https://developers.facebook.com/docs/whatsapp/pricing#free-entry-point-conversations
	PricingCategoryReferralConversion PricingCategory = "referral_conversion"
)

// Indicates conversation category. This can also be referred to as a conversation entry point
type OriginType string

const (
	// Indicates the conversation was opened by a business sending template
	// categorized as AUTHENTICATION to the customer.
	// This applies any time it has been more than 24 hours since the last customer message.
	OriginTypeAuthentication OriginType = "authentication"
	// Indicates the conversation was opened by a business sending template
	// categorized as MARKETING to the customer. This applies any
	// time it has been more than 24 hours since the last customer message.
	OriginTypeMarketing OriginType = "marketing"
	// Indicates the conversation was opened by a business sending template
	// categorized as UTILITY to the customer. This applies any time it
	// has been more than 24 hours since the last customer message.
	OriginTypeUtility OriginType = "utility"
	// Indicates that the conversation opened by a business replying
	// to a customer within a customer service window.
	OriginTypeService OriginType = "service"
	// Indicates a free entry point conversation.
	OriginTypeReferralConversion OriginType = "referral_conversion"
)

type CallObject struct {
	// The WhatsApp call ID
	ID string `json:"id"`
	// The WhatsApp user's phone number (callee)
	To        string              `json:"to"`
	From      string              `json:"from"`
	Event     string              `json:"event"`
	Timestamp string              `json:"timestamp"`
	Session   *CallSessionObject  `json:"session,omitempty"`
	Direction CallObjectDirection `json:"direction,omitempty"`
}

type CallObjectDirection string

const (
	CallObjectDirectionUserInitiated     CallObjectDirection = "USER_INITIATED"
	CallObjectDirectionBusinessInitiated CallObjectDirection = "BUSINESS_INITIATED"
)

type CallSessionObject struct {
	SDPType string `json:"sdp_type"`
	// RFC 8866 SDP
	SDP string `json:"sdp"`
}

type ErrorObject struct {
	Code      ErrorCode `json:"code"`
	Title     string    `json:"title"` // ErrorTitle
	Message   string    `json:"message,omitempty"`
	ErrorData struct {
		Details string `json:"details,omitempty"`
	} `json:"error_data,omitempty"`
}

type ErrorCode int

const (
	// Se não houver um subcódigo, isso indica que o status de login ou o token
	// de acesso expirou, foi revogado ou é inválido. Se houver um subcódigo, consulte-o.
	// Solução possível: obtenha um novo token de acesso
	ECodeAuthException ErrorCode = 0
	// Possivelmente, um problema temporário devido ao tempo de inatividade.
	// Aguarde um momento e refaça a operação.
	// Acesse o painel de status da WhatsApp Business API e verifique se não há
	// erros de ortografia na sua chamada de API.
	ECodeUnknownAPI ErrorCode = 1
	// Um problema temporário devido ao tempo de inatividade. Aguarde um momento
	// e refaça a operação.
	// Acesse o painel de status da WhatsApp Business API.
	ECodeAPIService ErrorCode = 2
	// Indica um problema que envolve recursos ou permissões.
	// Visite o ponto de extremidade referência para garantir que esteja incluindo a
	// permissão necessária na sua chamada.
	ECodeAPIMethod ErrorCode = 3
	// Um problema temporário devido à limitação. Aguarde um momento e refaça a
	// operação ou verifique o volume de solicitações de API.
	ECodeTooManyAPICalls ErrorCode = 4
	// A permissão não foi concedida ou foi removida. Saiba como lidar com
	// permissões ausentes.
	ECodePermissionDenied ErrorCode = 10
	// Talvez o parâmetro não esteja disponível ou pode estar escrito incorretamente.
	// Solução possível: visite o ponto de extremidade referência para verificar se o parâmetro existe.
	ECodeInvalidParameter ErrorCode = 100
)

// TODO: add error codes below

// 190 – Token de acesso expirado

// O token de acesso expirou.
// Solução possível: obtenha um novo token de acesso.

// 200-299 – Permissão de API (múltiplos valores dependem da permissão)

// A permissão não foi concedida ou foi removida. Saiba como lidar com permissões ausentes.

// 200 (subcódigo = 2494049) — Telefone não permitido

// Esse telefone não pode ser integrado à API de Nuvem.

// 368 – Bloqueio temporário devido a violações das políticas

// Aguarde um momento e refaça a operação. Saiba mais sobre a nossa Aplicação de políticas.

// 506 – Publicação duplicada

// Publicações duplicadas não podem ser publicadas consecutivamente.
// Solução possível: altere o conteúdo da publicação e tente novamente.

// 80007 – Problemas de limite de volume

// Você atingiu a limitação de volume da plataforma.

// 131052 – Erro ao baixar mídia

// Falha ao baixar a mídia do remetente.

// 131042 – Elegibilidade da empresa (problema de pagamento)

// Falha ao enviar mensagem devido a um ou mais erros relacionados à forma de pagamento.
// Para resolver esse erro, verifique se a sua conta do WhatsApp Business tem os seguintes itens configurados: fuso horário, moeda e forma de pagamento. Depois, verifique se a sua conta tem solicitações MessagingFor pendentes ou recusadas. Saiba mais sobre a configuração de cobrança.

// 131043 – Mensagem expirada

// Mensagem não enviada durante seu TTL (tempo de vida).

// 130429 – Limitação de volume atingida

// Falha ao enviar a mensagem porque o número de telefone fez envios demais em um curto período.

// 131045 – Certificado não assinado

// Falha ao enviar a mensagem porque ocorreu um erro relacionado ao certificado.

// 131016 – Serviço sobrecarregado

// O serviço está sobrecarregado. Aguarde um momento e refaça a operação. Em caso de problemas, visite o painel de status da Plataforma do WhatsApp Business.

// 131047 – Mensagem de reengajamento

// Falha ao enviar a mensagem porque mais de 24 horas se passaram desde a última vez que o cliente respondeu a este número.
// Solução possível: envie uma mensagem iniciada pela empresa usando um modelo de mensagem para responder.

// 131048 – Limite de taxa de spam atingido

// Falha ao enviar a mensagem porque há um limite de envios que podem ser feitos deste número de telefone. É possível que muitas mensagens anteriores tenham sido bloqueadas ou marcadas como spam.
// Verifique seu status de qualidade no Gerenciador do WhatsApp e visite a documentação Limitações de volume com base em qualidade para saber mais.

// 131000 – Erro genérico

// Falha ao enviar a mensagem devido a um erro desconhecido.

// 131001 – Mensagem muito longa

// O tamanho da mensagem é superior a 4.096 caracteres.

// 131002 – Tipo de destinatário inválido

// O valor recipient_type só pode ser individual.

// 131005 – Acesso negado

// O número já está registrado no WhatsApp. Consulte Como migrar um número de telefone.

// 131006 – Recurso não encontrado

// Arquivo ou recurso não encontrado.

// 131008 – Parâmetro obrigatório ausente

// Um parâmetro obrigatório está ausente. Acesse a referência de ponto de extremidade para saber mais sobre os requisitos de parâmetro.

// 131009 – Valor do parâmetro inválido

// O valor inserido para um parâmetro não é do tipo correto ou há outro problema. Acesse a referência de ponto de extremidade para saber mais sobre os requisitos de parâmetro. Esse erro também é retornado se o telefone do destinatário não for um número de telefone válido do WhatsApp ou se o telefone do remetente não for um número de telefone registrado no WhatsApp.

// 131021 – Usuário incorreto

// Essa mensagem será recebida quando você enviar uma mensagem para si mesmo.

// 131031 – A conta do remetente está bloqueada

// Sua conta foi bloqueada e não poderá enviar mensagens devido a uma violação da política de integridade.
// Consulte a Aplicação da Política da Plataforma do WhatsApp Business para saber mais.

// 131051 – Tipo de mensagem não aceito

// No momento, esse tipo de mensagem não é aceito.

// 131055 — Método não permitido

// O método que você está tentando usar não é permitido.

// 132000 – Incompatibilidade na contagem de parâmetros do modelo

// O número de parâmetros fornecidos não corresponde à quantidade esperada.

// 132001 – Modelo inexistente

// O modelo não existe no idioma especificado ou não foi aprovado.
// Verifique se o modelo foi aprovado e se o nome do modelo e a localidade do idioma estão corretos. Saiba mais.

// 132005 – Texto hidratado do modelo longo demais

// O texto traduzido é longo demais
// Verifique se o modelo foi traduzido corretamente. Saiba mais.

// 132007 – Política de caracteres do formato do modelo violada

// A política de caracteres do formato foi violada
// Saiba mais sobre as diretrizes de modelo de mensagem.

// 132012 – Incompatibilidade no formato do parâmetro do modelo

// O formato do parâmetro não corresponde ao do modelo criado.
// Verifique se o número de respostas do cliente retornado por um modelo é o número exigido por seu modelo. Saiba mais.

// 133000 – Exclusão incompleta

// Uma operação anterior de exclusão falhou. Antes de se registrar novamente, é necessária uma nova tentativa de exclusão.

// 133001 – Erro ao descriptografar

// Por um motivo desconhecido, a descriptografia do blob de backup falhou.

// 133002 – Erro ao descriptografar o blob de backup

// Falha ao descriptografar o blob de backup devido a formato inválido ou senha incorreta.

// 133003 – Erro ao descriptografar o token de recuperação

// Falha ao descriptografar o token de recuperação devido a formato inválido ou senha incorreta.

// 133004 – Servidor temporariamente indisponível

// O servidor de registro está temporariamente indisponível.

// 133005 – PIN de segurança incompatível

// O PIN está incorreto.

// Verifique se você está usando o PIN correto e tente de novo. Saiba mais.

// 133006 – Token de recuperação incorreto

// O token de recuperação usado para a migração está obsoleto.

// 133007 – Conta bloqueada

// A conta foi bloqueada pelo servidor de registro.

// 133008 – Muitas tentativas de PIN

// Foram feitas muitas tentativas de PIN nesta conta.

// 133009 – PIN inserido muito rapidamente

// A solicitação de registro foi feita com muita rapidez.
