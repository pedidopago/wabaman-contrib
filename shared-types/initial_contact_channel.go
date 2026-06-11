package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pedidopago/zajson"
)

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
	if len(raw) == 4 && raw[0] == 'n' {
		return nil, nil
	}
	if len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' {
		return nil, fmt.Errorf("expected JSON string for InitialContactChannel, got %s", raw)
	}
	unquoted := raw[1 : len(raw)-1]
	if interned, ok := initialContactChannelIntern[string(unquoted)]; ok {
		return interned, nil
	}
	s := string(unquoted)
	return &s, nil
}

func ZReadInitialContactChannel(r *zajson.Reader) (InitialContactChannel, error) {
	if r.PeekNull() {
		if err := r.ReadNull(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	s, err := r.ReadString()
	if err != nil {
		return nil, err
	}
	if interned, ok := initialContactChannelIntern[s]; ok {
		return interned, nil
	}
	cs := strings.Clone(s)
	return &cs, nil
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
