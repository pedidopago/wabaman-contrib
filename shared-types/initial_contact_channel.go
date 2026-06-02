package types

import "encoding/json"

type InitialContactChannel *string

var (
	initialContactChannelWhatsApp         = "whatsapp"
	initialContactChannelNewConversation  = "new_conversation"
	initialContactChannelZuckr            = "zuckr"
	initialContactChannelMarketplace      = "marketplace"
	initialContactChannelWebsite          = "website"
	initialContactChannelMsPrescription   = "ms-prescription"
	initialContactChannelWebsiteEcommerce = "website_ecommerce"
	initialContactChannelEmail            = "email"

	InitialContactChannelWhatsApp         InitialContactChannel = &initialContactChannelWhatsApp
	InitialContactChannelNewConversation  InitialContactChannel = &initialContactChannelNewConversation
	InitialContactChannelZuckr            InitialContactChannel = &initialContactChannelZuckr
	InitialContactChannelMarketplace      InitialContactChannel = &initialContactChannelMarketplace
	InitialContactChannelWebsite          InitialContactChannel = &initialContactChannelWebsite
	InitialContactChannelMsPrescription   InitialContactChannel = &initialContactChannelMsPrescription
	InitialContactChannelWebsiteEcommerce InitialContactChannel = &initialContactChannelWebsiteEcommerce
	InitialContactChannelEmail            InitialContactChannel = &initialContactChannelEmail
	InitialContactChannelDel              InitialContactChannel = &delStr
)

var initialContactChannelIntern = map[string]InitialContactChannel{
	"whatsapp":          InitialContactChannelWhatsApp,
	"new_conversation":  InitialContactChannelNewConversation,
	"zuckr":             InitialContactChannelZuckr,
	"marketplace":       InitialContactChannelMarketplace,
	"website":           InitialContactChannelWebsite,
	"ms-prescription":   InitialContactChannelMsPrescription,
	"website_ecommerce": InitialContactChannelWebsiteEcommerce,
	"email":             InitialContactChannelEmail,
	"$del":              InitialContactChannelDel,
}

func internInitialContactChannel(raw json.RawMessage) (InitialContactChannel, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	if interned, ok := initialContactChannelIntern[s]; ok {
		return interned, nil
	}
	return &s, nil
}

func EqualInitialContactChannel(a, b InitialContactChannel) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
