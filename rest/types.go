package rest

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pedidopago/go-common/util"
	"github.com/pedidopago/wabaman-contrib/fbgraph"
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
	NewMessageStatusImmediate         NewMessageStatus = "sent_immediately"
	NewMessageStatusQueuedForTemplate NewMessageStatus = "queued_for_template"
	NewMessageStatusUnknown           NewMessageStatus = "unknown"
)

type NewMessageRequest struct {
	BranchID         string                            `json:"branch_id" validate:"required"`
	FromNumber       string                            `json:"from_number"`
	ToNumber         string                            `json:"to_number"`
	Type             MessageType                       `json:"type"`
	Text             *fbgraph.TextObject               `json:"text,omitempty"`
	Template         *fbgraph.TemplateObject           `json:"template,omitempty"`
	Interactive      *fbgraph.InteractiveMessageObject `json:"interactive,omitempty"`
	Image            *fbgraph.MediaObject              `json:"image,omitempty"`
	Audio            *fbgraph.MediaObject              `json:"audio,omitempty"`
	Document         *fbgraph.MediaObject              `json:"document,omitempty"`
	Video            *fbgraph.MediaObject              `json:"video,omitempty"`
	Sticker          *fbgraph.MediaObject              `json:"sticker,omitempty"`
	FallbackTemplate string                            `json:"fallback_template,omitempty"`
	SkipWelcome      bool                              `json:"skip_welcome,omitempty"`
	Origin           string                            `json:"origin,omitempty"`
	ReadMessages     bool                              `json:"read_messages,omitempty"`
	Verbose          bool                              `json:"verbose,omitempty"`
	ContactMetadata  map[string]any                    `json:"contact_metadata,omitempty"`
	AgentID          string                            `json:"agent_id,omitempty"`
	AgentName        string                            `json:"agent_name,omitempty"`
}

type NewMessageResponse struct {
	MessageID         string           `json:"message_id"`
	SendMessageStatus NewMessageStatus `json:"send_message_status"`
}

type NewMediaResponse struct {
	MediaID string `json:"media_id"`
}

type UpdateContactRequest struct {
	CustomerID      string         `json:"customer_id,omitempty"`
	CustomerName    string         `json:"customer_name,omitempty"`
	WABAProfileName string         `json:"waba_profile_name,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
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
	ID                           uint64    `json:"id,omitempty"`                              // "id": 1,
	PhoneID                      uint      `json:"phone_id,omitempty"`                        // "phone_id": 1,
	WabaContactID                string    `json:"waba_contact_id,omitempty"`                 // "waba_contact_id": "5511941011935",
	WabaProfileName              string    `json:"waba_profile_name,omitempty"`               // "waba_profile_name": "Gabriel Ochsenhofer",
	LastActivity                 time.Time `json:"last_activity,omitempty"`                   // "last_activity": "2022-10-03T19:32:19Z",
	LastMessageReceivedID        uint64    `json:"last_message_received_id,omitempty"`        // "last_message_received_id": 145,
	LastMessageSentId            uint64    `json:"last_message_sent_id,omitempty"`            // "last_message_sent_id": 918,
	LastMessageReceivedTimestamp time.Time `json:"last_message_received_timestamp,omitempty"` // "last_message_received_timestamp": "2022-10-03T19:32:19Z",
	LastMessageSentTimestamp     time.Time `json:"last_message_sent_timestamp,omitempty"`     // "last_message_sent_timestamp": "2022-10-04T02:30:23Z",
	CustomerID                   string    `json:"customer_id,omitempty"`                     // "customer_id": "01F5E1TNWH1TCTGJ1VW71X1NA8",
	CreatedAt                    time.Time `json:"created_at,omitempty"`                      // "created_at": "2022-09-02T14:04:08Z",
	UpdatedAt                    time.Time `json:"updated_at,omitempty"`                      // "updated_at": "2022-10-04T02:30:26Z",

	Metadata map[string]any `json:"metadata,omitempty"`
	//TODO: add fields below
	// HostMessages string `json:"host_messages,omitempty"` // "host_messages": null,
	// ClientMessages string `json:"client_messages,omitempty"` // "client_messages": null
}

type GetBusinessesRequest struct {
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
	ID                        uint      `json:"id"`
	BusinessID                uint      `json:"business_id"`
	WhatsAppID                string    `json:"whatsapp_id"`
	WhatsAppBusinessAccountID string    `json:"whatsapp_business_account_id"`
	PhoneNumber               string    `json:"phone_number"`
	BranchID                  string    `json:"branch_id"`
	TemplateNamespace         string    `json:"template_namespace,omitempty"`
	DefaultTplHeaderImage     string    `json:"default_tpl_header_image,omitempty"`
	DefaultTplHeaderVideo     string    `json:"default_tpl_header_video,omitempty"`
	DefaultReheatTemplate     string    `json:"default_reheat_template,omitempty"`
	FBAppID                   string    `json:"fb_app_id"`
	AlertEmail                string    `json:"alert_email,omitempty"`
	AlertDiscord              string    `json:"alert_discord,omitempty"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}
