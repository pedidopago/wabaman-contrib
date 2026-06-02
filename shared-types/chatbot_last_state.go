package types

import "encoding/json"

type ChatbotLastState *string

var (
	chatbotLastStateCreateInquiry                  = "chatbot_create_inquiry"
	chatbotLastStatePause                          = "chatbot_pause"
	chatbotLastStateAskForNameOrCreateInquiry      = "chatbot_ask_for_name_or_create_inquiry"
	chatbotLastStateConfirmUnusualName             = "chatbot_confirm_unusual_name"
	chatbotLastStateAskAgeThreshold                = "chatbot_ask_age_threshold"
	chatbotLastStateAskForDocumentsOrCreateInquiry = "chatbot_ask_for_documents_or_create_inquiry"
	chatbotLastStateWelcome                        = "chatbot_welcome"
	chatbotLastStateInquiryWaitForAttachments      = "chatbot_inquiry_wait_for_attachments"
	chatbotLastStateEmpty                          = ""
	chatbotLastStateHelpHuman                      = "chatbot_help_human"
	chatbotLastStateNewInquiry                     = "chatbot_new_inquiry"

	ChatbotLastStateCreateInquiry                  ChatbotLastState = &chatbotLastStateCreateInquiry
	ChatbotLastStatePause                          ChatbotLastState = &chatbotLastStatePause
	ChatbotLastStateAskForNameOrCreateInquiry      ChatbotLastState = &chatbotLastStateAskForNameOrCreateInquiry
	ChatbotLastStateConfirmUnusualName             ChatbotLastState = &chatbotLastStateConfirmUnusualName
	ChatbotLastStateAskAgeThreshold                ChatbotLastState = &chatbotLastStateAskAgeThreshold
	ChatbotLastStateAskForDocumentsOrCreateInquiry ChatbotLastState = &chatbotLastStateAskForDocumentsOrCreateInquiry
	ChatbotLastStateWelcome                        ChatbotLastState = &chatbotLastStateWelcome
	ChatbotLastStateInquiryWaitForAttachments      ChatbotLastState = &chatbotLastStateInquiryWaitForAttachments
	ChatbotLastStateEmpty                          ChatbotLastState = &chatbotLastStateEmpty
	ChatbotLastStateHelpHuman                      ChatbotLastState = &chatbotLastStateHelpHuman
	ChatbotLastStateNewInquiry                     ChatbotLastState = &chatbotLastStateNewInquiry
	ChatbotLastStateDel                            ChatbotLastState = &delStr
)

var chatbotLastStateIntern = map[string]ChatbotLastState{
	"chatbot_create_inquiry":                      ChatbotLastStateCreateInquiry,
	"chatbot_pause":                               ChatbotLastStatePause,
	"chatbot_ask_for_name_or_create_inquiry":      ChatbotLastStateAskForNameOrCreateInquiry,
	"chatbot_confirm_unusual_name":                ChatbotLastStateConfirmUnusualName,
	"chatbot_ask_age_threshold":                   ChatbotLastStateAskAgeThreshold,
	"chatbot_ask_for_documents_or_create_inquiry": ChatbotLastStateAskForDocumentsOrCreateInquiry,
	"chatbot_welcome":                             ChatbotLastStateWelcome,
	"chatbot_inquiry_wait_for_attachments":        ChatbotLastStateInquiryWaitForAttachments,
	"":                                            ChatbotLastStateEmpty,
	"chatbot_help_human":                          ChatbotLastStateHelpHuman,
	"chatbot_new_inquiry":                         ChatbotLastStateNewInquiry,
	"$del":                                        ChatbotLastStateDel,
}

func internChatbotLastState(raw json.RawMessage) (ChatbotLastState, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	if interned, ok := chatbotLastStateIntern[s]; ok {
		return interned, nil
	}
	return &s, nil
}

func EqualChatbotLastState(a, b ChatbotLastState) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
