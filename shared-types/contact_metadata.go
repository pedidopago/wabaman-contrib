package types

import (
	"encoding/json"
	"time"
)

type ContactMetadata struct {
	InquiryID                       *string                      `json:"inquiry_id,omitzero"`
	InquiryStatus                   InquiryStatus                `json:"inquiry_status,omitzero"`
	InquiryDisplayID                *string                      `json:"inquiry_display_id,omitzero"`
	InquiryAgentID                  *string                      `json:"inquiry_agent_id,omitzero"`
	InquiryAgentName                *string                      `json:"inquiry_agent_name,omitzero"`
	InquiryAIEvaluation             *string                      `json:"inquiry_ai_evaluation,omitzero"`
	InquiryCanBindDisplayID         MetadataBool                 `json:"inquiry_can_bind_display_id,omitzero"`
	InquiryCreatedAt                *time.Time                   `json:"inquiry_created_at,omitzero"`
	InquiryExpireDate               *time.Time                   `json:"inquiry_expire_date,omitzero"`
	InquirySellerAgentID            *string                      `json:"inquiry_seller_agent_id,omitzero"`
	InquirySellerAgentName          *string                      `json:"inquiry_seller_agent_name,omitzero"`
	InquiryHasPendencies            MetadataBool                 `json:"inquiry_has_pendencies,omitzero"`
	InquirySellOpportunityCollected MetadataBool                 `json:"inquiry_sell_opportunity_collected,omitzero"`
	InquiryQuotatedAt               *time.Time                   `json:"inquiry_quotated_at,omitzero"`
	InquiryDoneAt                   *time.Time                   `json:"inquiry_done_at,omitzero"`
	InquiryLastStatusUpdate         *time.Time                   `json:"inquiry_last_status_update,omitzero"`
	InquiryIsChatOpen               MetadataBool                 `json:"inquiry_is_chat_open,omitzero"`
	InquiryIsMarketplace            MetadataBool                 `json:"inquiry_is_marketplace,omitzero"`
	InquiryInclusorAgentID          *string                      `json:"inquiry_inclusor_agent_id,omitzero"`
	InquiryInclusorAgentName        *string                      `json:"inquiry_inclusor_agent_name,omitzero"`
	InquirySpecialistAgentID        *string                      `json:"inquiry_specialist_agent_id,omitzero"`
	InquirySpecialistAgentName      *string                      `json:"inquiry_specialist_agent_name,omitzero"`
	AccountID                       *string                      `json:"account_id,omitzero"`
	ChatbotDisabled                 MetadataBool                 `json:"chatbot_disabled,omitzero"`
	ChatbotInitialContact           *time.Time                   `json:"chatbot_initial_contact,omitzero"`
	ChatbotIsPreRegistration        MetadataBool                 `json:"chatbot_is_pre_registration,omitzero"`
	ChatbotLastState                ChatbotLastState             `json:"chatbot_last_state,omitzero"`
	ChatbotRegistrationDate         *time.Time                   `json:"chatbot_registration_date,omitzero"`
	InitialContactChannel           InitialContactChannel        `json:"initial_contact_channel,omitzero"`
	InitialContactDate              *time.Time                   `json:"initial_contact_date,omitzero"`
	CustomerName                    *string                      `json:"customer_name,omitzero"`
	CustomerDocumentCountry         CustomerDocumentCountry      `json:"customer_document_country,omitzero"`
	CustomerDocumentType            *string                      `json:"customer_document_type,omitzero"`
	CustomerDocument                *string                      `json:"customer_document,omitzero"`
	LastOrderSeq                    *int                         `json:"last_order_seq,omitzero"`
	LastCouponOffered               *string                      `json:"last_coupon_offered,omitzero"`
	Order                           *ContactMetadataOrder        `json:"order,omitzero"`
	Prescription                    *ContactMetadataPrescription `json:"prescription,omitzero"`
	OtherFields                     map[string]any               `json:"-"`
}

var contactMetadataKnownKeys = map[string]struct{}{
	"inquiry_id":                         {},
	"inquiry_status":                     {},
	"inquiry_display_id":                 {},
	"inquiry_agent_id":                   {},
	"inquiry_agent_name":                 {},
	"inquiry_ai_evaluation":              {},
	"inquiry_can_bind_display_id":        {},
	"inquiry_created_at":                 {},
	"inquiry_expire_date":                {},
	"inquiry_seller_agent_id":            {},
	"inquiry_seller_agent_name":          {},
	"inquiry_has_pendencies":             {},
	"inquiry_sell_opportunity_collected": {},
	"inquiry_quotated_at":                {},
	"inquiry_done_at":                    {},
	"inquiry_last_status_update":         {},
	"inquiry_is_chat_open":               {},
	"inquiry_is_marketplace":             {},
	"inquiry_inclusor_agent_id":          {},
	"inquiry_inclusor_agent_name":        {},
	"inquiry_specialist_agent_id":        {},
	"inquiry_specialist_agent_name":      {},
	"account_id":                         {},
	"chatbot_disabled":                   {},
	"chatbot_initial_contact":            {},
	"chatbot_is_pre_registration":        {},
	"chatbot_last_state":                 {},
	"chatbot_registration_date":          {},
	"initial_contact_channel":            {},
	"initial_contact_date":               {},
	"customer_name":                      {},
	"customer_document_country":          {},
	"customer_document_type":             {},
	"customer_document":                  {},
	"last_order_seq":                     {},
	"last_coupon_offered":                {},
	"order":                              {},
	"prescription":                       {},
}

func (cm ContactMetadata) MarshalJSON() ([]byte, error) {
	type alias ContactMetadata
	knownBytes, err := json.Marshal(alias(cm))
	if err != nil {
		return nil, err
	}

	if len(cm.OtherFields) == 0 {
		return knownBytes, nil
	}

	var merged map[string]json.RawMessage
	if err := json.Unmarshal(knownBytes, &merged); err != nil {
		return nil, err
	}

	for k, v := range cm.OtherFields {
		if _, known := contactMetadataKnownKeys[k]; known {
			continue
		}
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		merged[k] = raw
	}

	return json.Marshal(merged)
}

func (cm *ContactMetadata) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	known := map[string]any{
		"inquiry_id":                    &cm.InquiryID,
		"inquiry_display_id":            &cm.InquiryDisplayID,
		"inquiry_agent_id":              &cm.InquiryAgentID,
		"inquiry_agent_name":            &cm.InquiryAgentName,
		"inquiry_ai_evaluation":         &cm.InquiryAIEvaluation,
		"inquiry_created_at":            &cm.InquiryCreatedAt,
		"inquiry_expire_date":           &cm.InquiryExpireDate,
		"inquiry_seller_agent_id":       &cm.InquirySellerAgentID,
		"inquiry_seller_agent_name":     &cm.InquirySellerAgentName,
		"inquiry_quotated_at":           &cm.InquiryQuotatedAt,
		"inquiry_done_at":               &cm.InquiryDoneAt,
		"inquiry_last_status_update":    &cm.InquiryLastStatusUpdate,
		"inquiry_inclusor_agent_id":     &cm.InquiryInclusorAgentID,
		"inquiry_inclusor_agent_name":   &cm.InquiryInclusorAgentName,
		"inquiry_specialist_agent_id":   &cm.InquirySpecialistAgentID,
		"inquiry_specialist_agent_name": &cm.InquirySpecialistAgentName,
		"account_id":                    &cm.AccountID,
		"chatbot_initial_contact":       &cm.ChatbotInitialContact,
		"chatbot_registration_date":     &cm.ChatbotRegistrationDate,
		"initial_contact_date":          &cm.InitialContactDate,
		"customer_name":                 &cm.CustomerName,
		"customer_document_type":        &cm.CustomerDocumentType,
		"customer_document":             &cm.CustomerDocument,
		"last_order_seq":                &cm.LastOrderSeq,
		"last_coupon_offered":           &cm.LastCouponOffered,
		"order":                         &cm.Order,
		"prescription":                  &cm.Prescription,
	}

	// Interned enum fields
	if v, ok := raw["inquiry_status"]; ok {
		p, err := internInquiryStatus(v)
		if err != nil {
			return err
		}
		cm.InquiryStatus = p
	}
	if v, ok := raw["chatbot_last_state"]; ok {
		p, err := internChatbotLastState(v)
		if err != nil {
			return err
		}
		cm.ChatbotLastState = p
	}
	for _, bf := range []struct {
		key string
		dst *MetadataBool
	}{
		{"inquiry_can_bind_display_id", &cm.InquiryCanBindDisplayID},
		{"inquiry_has_pendencies", &cm.InquiryHasPendencies},
		{"inquiry_sell_opportunity_collected", &cm.InquirySellOpportunityCollected},
		{"inquiry_is_chat_open", &cm.InquiryIsChatOpen},
		{"inquiry_is_marketplace", &cm.InquiryIsMarketplace},
		{"chatbot_disabled", &cm.ChatbotDisabled},
		{"chatbot_is_pre_registration", &cm.ChatbotIsPreRegistration},
	} {
		if v, ok := raw[bf.key]; ok {
			p, err := internMetadataBool(v)
			if err != nil {
				return err
			}
			*bf.dst = p
		}
	}
	if v, ok := raw["initial_contact_channel"]; ok {
		p, err := internInitialContactChannel(v)
		if err != nil {
			return err
		}
		cm.InitialContactChannel = p
	}
	if v, ok := raw["customer_document_country"]; ok {
		p, err := internCustomerDocumentCountry(v)
		if err != nil {
			return err
		}
		cm.CustomerDocumentCountry = p
	}

	for key, dst := range known {
		if v, ok := raw[key]; ok {
			if err := json.Unmarshal(v, dst); err != nil {
				return err
			}
		}
	}

	for k, v := range raw {
		if _, isKnown := contactMetadataKnownKeys[k]; isKnown {
			continue
		}
		if cm.OtherFields == nil {
			cm.OtherFields = make(map[string]any)
		}
		var decoded any
		if err := json.Unmarshal(v, &decoded); err != nil {
			return err
		}
		cm.OtherFields[k] = decoded
	}

	return nil
}
