package msgdriver

type WebhookType string

const (
	WebhookTypeSendMessage     WebhookType = "send_message"
	WebhookTypeNotifyMessageID WebhookType = "notify_message_id"
	WebhookTypeReadMessage     WebhookType = "read_message"
	WebhookTypeQuery           WebhookType = "query"
	WebhookTypeHostNote        WebhookType = "host_note"
)

type WebhookRequest struct {
	Token    string      `json:"token"`
	Business Business    `json:"business"`
	Phone    Phone       `json:"phone"`
	Type     WebhookType `json:"type"`

	// one of:

	SendMessage     *SendMessage     `json:"send_message,omitempty"`
	NotifyMessageID *NotifyMessageID `json:"notify_message_id,omitempty"`
	ReadMessage     *ReadMessage     `json:"read_message,omitempty"`
	Query           *Query           `json:"query,omitempty"`
	HostNote        *HostNote        `json:"host_note,omitempty"`
}

type WebhookResponse struct {
	Type WebhookType `json:"type"`

	// one of:

	SendMessage *SendMessageResult `json:"send_message,omitempty"`
	Query       *QueryResult       `json:"query,omitempty"`
}

// Business is equivalent to a Store (PP).
type Business struct {
	ID      uint   `json:"id"`
	StoreID string `json:"store_id"`
}

// Phone is the phone of the company. This may be translated to a branch if the driver
// does not use the concept of phones.
type Phone struct {
	ID                         uint   `json:"id"`
	BranchID                   string `json:"branch_id"`
	BranchName                 string `json:"branch_name"`
	DefaultTemplateImage       string `json:"default_template_image,omitempty"`
	DefaultTemplateVideo       string `json:"default_template_video,omitempty"`
	DefaultTemplateCompanyName string `json:"default_template_company_name,omitempty"`
	DriverName                 string `json:"driver_name"`
	DriverData                 string `json:"driver_data"` // 1000 char data to use with the driver
	PhoneNumber                string `json:"phone_number"`
	LanguageCode               string `json:"language_code,omitempty"`
}
