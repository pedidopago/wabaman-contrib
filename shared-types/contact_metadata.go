package types

import (
	"encoding/json"
	"time"
)

//go:generate zajson -heap -types ContactMetadata,ContactMetadataOrder,OrderPrescription,ContactMetadataPrescription,PrescriptionCustomer,PrescriptionCustomerAddress

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
	OtherFields                     map[string]any               `json:"-" zajson:"-,remain"`
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

	for k, v := range raw {
		var err error
		switch k {
		case "inquiry_id":
			err = json.Unmarshal(v, &cm.InquiryID)
		case "inquiry_status":
			cm.InquiryStatus, err = internInquiryStatus(v)
		case "inquiry_display_id":
			err = json.Unmarshal(v, &cm.InquiryDisplayID)
		case "inquiry_agent_id":
			err = json.Unmarshal(v, &cm.InquiryAgentID)
		case "inquiry_agent_name":
			err = json.Unmarshal(v, &cm.InquiryAgentName)
		case "inquiry_ai_evaluation":
			err = json.Unmarshal(v, &cm.InquiryAIEvaluation)
		case "inquiry_can_bind_display_id":
			cm.InquiryCanBindDisplayID, err = internMetadataBool(v)
		case "inquiry_created_at":
			err = json.Unmarshal(v, &cm.InquiryCreatedAt)
		case "inquiry_expire_date":
			err = json.Unmarshal(v, &cm.InquiryExpireDate)
		case "inquiry_seller_agent_id":
			err = json.Unmarshal(v, &cm.InquirySellerAgentID)
		case "inquiry_seller_agent_name":
			err = json.Unmarshal(v, &cm.InquirySellerAgentName)
		case "inquiry_has_pendencies":
			cm.InquiryHasPendencies, err = internMetadataBool(v)
		case "inquiry_sell_opportunity_collected":
			cm.InquirySellOpportunityCollected, err = internMetadataBool(v)
		case "inquiry_quotated_at":
			err = json.Unmarshal(v, &cm.InquiryQuotatedAt)
		case "inquiry_done_at":
			err = json.Unmarshal(v, &cm.InquiryDoneAt)
		case "inquiry_last_status_update":
			err = json.Unmarshal(v, &cm.InquiryLastStatusUpdate)
		case "inquiry_is_chat_open":
			cm.InquiryIsChatOpen, err = internMetadataBool(v)
		case "inquiry_is_marketplace":
			cm.InquiryIsMarketplace, err = internMetadataBool(v)
		case "inquiry_inclusor_agent_id":
			err = json.Unmarshal(v, &cm.InquiryInclusorAgentID)
		case "inquiry_inclusor_agent_name":
			err = json.Unmarshal(v, &cm.InquiryInclusorAgentName)
		case "inquiry_specialist_agent_id":
			err = json.Unmarshal(v, &cm.InquirySpecialistAgentID)
		case "inquiry_specialist_agent_name":
			err = json.Unmarshal(v, &cm.InquirySpecialistAgentName)
		case "account_id":
			err = json.Unmarshal(v, &cm.AccountID)
		case "chatbot_disabled":
			cm.ChatbotDisabled, err = internMetadataBool(v)
		case "chatbot_initial_contact":
			err = json.Unmarshal(v, &cm.ChatbotInitialContact)
		case "chatbot_is_pre_registration":
			cm.ChatbotIsPreRegistration, err = internMetadataBool(v)
		case "chatbot_last_state":
			cm.ChatbotLastState, err = internChatbotLastState(v)
		case "chatbot_registration_date":
			err = json.Unmarshal(v, &cm.ChatbotRegistrationDate)
		case "initial_contact_channel":
			cm.InitialContactChannel, err = internInitialContactChannel(v)
		case "initial_contact_date":
			err = json.Unmarshal(v, &cm.InitialContactDate)
		case "customer_name":
			err = json.Unmarshal(v, &cm.CustomerName)
		case "customer_document_country":
			cm.CustomerDocumentCountry, err = internCustomerDocumentCountry(v)
		case "customer_document_type":
			err = json.Unmarshal(v, &cm.CustomerDocumentType)
		case "customer_document":
			err = json.Unmarshal(v, &cm.CustomerDocument)
		case "last_order_seq":
			err = json.Unmarshal(v, &cm.LastOrderSeq)
		case "last_coupon_offered":
			err = json.Unmarshal(v, &cm.LastCouponOffered)
		case "order":
			err = json.Unmarshal(v, &cm.Order)
		case "prescription":
			err = json.Unmarshal(v, &cm.Prescription)
		default:
			if cm.OtherFields == nil {
				cm.OtherFields = make(map[string]any)
			}
			var decoded any
			if err := json.Unmarshal(v, &decoded); err != nil {
				return err
			}
			cm.OtherFields[k] = decoded
		}
		if err != nil {
			return err
		}
	}

	return nil
}
