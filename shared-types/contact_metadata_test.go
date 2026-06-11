package types

import (
	"encoding/json"
	"testing"
	"time"
)

func ptr[T any](v T) *T { return &v }

func TestContactMetadata_UnmarshalOnlyKnownFields(t *testing.T) {
	input := `{
		"inquiry_id": "inq-123",
		"inquiry_status": "open",
		"chatbot_disabled": true,
		"account_id": "acc-456"
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	if cm.InquiryID == nil || *cm.InquiryID != "inq-123" {
		t.Errorf("InquiryID = %v, want inq-123", cm.InquiryID)
	}
	if cm.InquiryStatus == nil || *cm.InquiryStatus != "open" {
		t.Errorf("InquiryStatus = %v, want open", cm.InquiryStatus)
	}
	if cm.ChatbotDisabled != MetadataBoolTrue {
		t.Errorf("ChatbotDisabled pointer comparison failed")
	}
	if cm.AccountID == nil || *cm.AccountID != "acc-456" {
		t.Errorf("AccountID = %v, want acc-456", cm.AccountID)
	}
	if cm.OtherFields != nil {
		t.Errorf("OtherFields = %v, want nil", cm.OtherFields)
	}
}

func TestContactMetadata_UnmarshalOnlyUnknownFields(t *testing.T) {
	input := `{"custom_field": "value", "priority": 5}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	if cm.InquiryID != nil {
		t.Errorf("InquiryID = %v, want nil", cm.InquiryID)
	}
	if len(cm.OtherFields) != 2 {
		t.Fatalf("OtherFields has %d entries, want 2", len(cm.OtherFields))
	}
	if cm.OtherFields["custom_field"] != "value" {
		t.Errorf("OtherFields[custom_field] = %v, want value", cm.OtherFields["custom_field"])
	}
	if cm.OtherFields["priority"] != float64(5) {
		t.Errorf("OtherFields[priority] = %v, want 5", cm.OtherFields["priority"])
	}
}

func TestContactMetadata_UnmarshalMixed(t *testing.T) {
	input := `{
		"inquiry_id": "inq-1",
		"inquiry_agent_name": "Agent Smith",
		"custom_tag": "vip",
		"score": 99.5
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	if cm.InquiryID == nil || *cm.InquiryID != "inq-1" {
		t.Errorf("InquiryID = %v, want inq-1", cm.InquiryID)
	}
	if cm.InquiryAgentName == nil || *cm.InquiryAgentName != "Agent Smith" {
		t.Errorf("InquiryAgentName = %v, want Agent Smith", cm.InquiryAgentName)
	}
	if cm.OtherFields["custom_tag"] != "vip" {
		t.Errorf("OtherFields[custom_tag] = %v, want vip", cm.OtherFields["custom_tag"])
	}
	if cm.OtherFields["score"] != 99.5 {
		t.Errorf("OtherFields[score] = %v, want 99.5", cm.OtherFields["score"])
	}
	if _, exists := cm.OtherFields["inquiry_id"]; exists {
		t.Error("OtherFields should not contain known key inquiry_id")
	}
}

func TestContactMetadata_UnmarshalTimestamps(t *testing.T) {
	input := `{
		"inquiry_created_at": "2025-01-15T10:30:00Z",
		"inquiry_expire_date": "2025-02-15T10:30:00Z"
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	want := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	if cm.InquiryCreatedAt == nil || !cm.InquiryCreatedAt.Equal(want) {
		t.Errorf("InquiryCreatedAt = %v, want %v", cm.InquiryCreatedAt, want)
	}

	wantExpire := time.Date(2025, 2, 15, 10, 30, 0, 0, time.UTC)
	if cm.InquiryExpireDate == nil || !cm.InquiryExpireDate.Equal(wantExpire) {
		t.Errorf("InquiryExpireDate = %v, want %v", cm.InquiryExpireDate, wantExpire)
	}
}

func TestContactMetadata_UnmarshalNestedObject(t *testing.T) {
	input := `{
		"inquiry_id": "inq-1",
		"nested": {"key": "val", "num": 42}
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	nested, ok := cm.OtherFields["nested"].(map[string]any)
	if !ok {
		t.Fatalf("OtherFields[nested] = %T, want map[string]any", cm.OtherFields["nested"])
	}
	if nested["key"] != "val" {
		t.Errorf("nested.key = %v, want val", nested["key"])
	}
}

func TestContactMetadata_UnmarshalEmpty(t *testing.T) {
	var cm ContactMetadata
	if err := json.Unmarshal([]byte(`{}`), &cm); err != nil {
		t.Fatal(err)
	}

	if cm.InquiryID != nil {
		t.Error("InquiryID should be nil")
	}
	if cm.OtherFields != nil {
		t.Errorf("OtherFields = %v, want nil", cm.OtherFields)
	}
}

func TestContactMetadata_MarshalOnlyKnownFields(t *testing.T) {
	cm := ContactMetadata{
		InquiryID:       ptr("inq-1"),
		InquiryStatus:   InquiryStatusPending,
		ChatbotDisabled: MetadataBoolTrue,
	}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if m["inquiry_id"] != "inq-1" {
		t.Errorf("inquiry_id = %v, want inq-1", m["inquiry_id"])
	}
	if m["inquiry_status"] != "PENDING" {
		t.Errorf("inquiry_status = %v, want PENDING", m["inquiry_status"])
	}
	if m["chatbot_disabled"] != true {
		t.Errorf("chatbot_disabled = %v, want true", m["chatbot_disabled"])
	}
	if _, exists := m["inquiry_agent_id"]; exists {
		t.Error("nil fields should be omitted")
	}
}

func TestContactMetadata_MarshalOnlyOtherFields(t *testing.T) {
	cm := ContactMetadata{
		OtherFields: map[string]any{
			"custom": "value",
			"count":  42,
		},
	}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if m["custom"] != "value" {
		t.Errorf("custom = %v, want value", m["custom"])
	}
	if m["count"] != float64(42) {
		t.Errorf("count = %v, want 42", m["count"])
	}
}

func TestContactMetadata_MarshalMixed(t *testing.T) {
	cm := ContactMetadata{
		InquiryID: ptr("inq-1"),
		AccountID: ptr("acc-2"),
		OtherFields: map[string]any{
			"custom_tag": "vip",
		},
	}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if m["inquiry_id"] != "inq-1" {
		t.Errorf("inquiry_id = %v, want inq-1", m["inquiry_id"])
	}
	if m["account_id"] != "acc-2" {
		t.Errorf("account_id = %v, want acc-2", m["account_id"])
	}
	if m["custom_tag"] != "vip" {
		t.Errorf("custom_tag = %v, want vip", m["custom_tag"])
	}
}

func TestContactMetadata_MarshalStructFieldWinsOverOtherFields(t *testing.T) {
	cm := ContactMetadata{
		InquiryID: ptr("correct"),
		OtherFields: map[string]any{
			"inquiry_id": "wrong",
		},
	}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if m["inquiry_id"] != "correct" {
		t.Errorf("inquiry_id = %v, want correct (struct field should win)", m["inquiry_id"])
	}
}

func TestContactMetadata_MarshalNilFieldNoFallbackToOtherFields(t *testing.T) {
	cm := ContactMetadata{
		OtherFields: map[string]any{
			"inquiry_id": "should-be-ignored",
			"custom":     "should-appear",
		},
	}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if _, exists := m["inquiry_id"]; exists {
		t.Error("known key from OtherFields should NOT appear when struct field is nil (no fallback)")
	}
	if m["custom"] != "should-appear" {
		t.Errorf("custom = %v, want should-appear", m["custom"])
	}
}

func TestContactMetadata_MarshalEmpty(t *testing.T) {
	cm := ContactMetadata{}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "{}" {
		t.Errorf("empty marshal = %s, want {}", string(data))
	}
}

func TestContactMetadata_RoundTrip(t *testing.T) {
	input := `{
		"inquiry_id": "inq-1",
		"inquiry_status": "closed",
		"chatbot_disabled": false,
		"inquiry_created_at": "2025-06-01T12:00:00Z",
		"custom_field": "hello",
		"nested": {"a": 1}
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var roundTripped map[string]any
	if err := json.Unmarshal(data, &roundTripped); err != nil {
		t.Fatal(err)
	}

	var original map[string]any
	if err := json.Unmarshal([]byte(input), &original); err != nil {
		t.Fatal(err)
	}

	for key, wantVal := range original {
		gotVal, exists := roundTripped[key]
		if !exists {
			t.Errorf("round-trip missing key %q", key)
			continue
		}
		wantJSON, _ := json.Marshal(wantVal)
		gotJSON, _ := json.Marshal(gotVal)
		if string(wantJSON) != string(gotJSON) {
			t.Errorf("round-trip key %q: got %s, want %s", key, gotJSON, wantJSON)
		}
	}

	for key := range roundTripped {
		if _, exists := original[key]; !exists {
			t.Errorf("round-trip has extra key %q", key)
		}
	}
}

func TestContactMetadata_UnmarshalOverflowNeverDuplicatesKnownFields(t *testing.T) {
	input := `{
		"inquiry_id": "inq-1",
		"inquiry_status": "open",
		"inquiry_display_id": "D-100",
		"inquiry_agent_id": "agent-1",
		"inquiry_agent_name": "Smith",
		"inquiry_ai_evaluation": "positive",
		"inquiry_can_bind_display_id": true,
		"inquiry_created_at": "2025-01-01T00:00:00Z",
		"inquiry_expire_date": "2025-12-31T23:59:59Z",
		"inquiry_seller_agent_id": "seller-1",
		"inquiry_seller_agent_name": "Seller",
		"inquiry_has_pendencies": true,
		"inquiry_sell_opportunity_collected": false,
		"inquiry_quotated_at": "2025-06-01T00:00:00Z",
		"inquiry_done_at": "2025-06-02T00:00:00Z",
		"inquiry_last_status_update": "2025-06-01T12:00:00Z",
		"inquiry_inclusor_agent_id": "inc-1",
		"inquiry_inclusor_agent_name": "Inclusor",
		"inquiry_specialist_agent_id": "spec-1",
		"inquiry_specialist_agent_name": "Specialist",
"account_id": "acc-1",
		"chatbot_disabled": false,
		"extra": "overflow"
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	if len(cm.OtherFields) != 1 {
		t.Errorf("OtherFields has %d entries, want 1: %v", len(cm.OtherFields), cm.OtherFields)
	}
	if cm.OtherFields["extra"] != "overflow" {
		t.Errorf("OtherFields[extra] = %v, want overflow", cm.OtherFields["extra"])
	}
}

func TestContactMetadata_RealWorldPayload(t *testing.T) {
	input := `{
		"chatbot_initial_contact": "2026-04-10T22:07:21Z",
		"chatbot_is_pre_registration": false,
		"chatbot_last_state": "chatbot_create_inquiry",
		"chatbot_name_updated": "Luis Zlochevsky at 2026-05-25T23:38:22Z (runListenPubSubCustomerUpdated)",
		"chatbot_registration_date": "2026-04-10T22:07:33.205684934Z",
		"has_pendencies": false,
		"initial_contact_channel": "whatsapp",
		"initial_contact_date": "2026-04-10T22:07:33.20390498Z",
		"inquiry_agent_id": "01EM2CNQFZFSV97CMA1SSNB67F",
		"inquiry_agent_name": "Luis Zlochevsky",
		"inquiry_ai_evaluation": "OTHERS",
		"inquiry_can_bind_display_id": false,
		"inquiry_created_at": "2026-05-25T23:38:22.258Z",
		"inquiry_display_id": "31049213",
		"inquiry_expire_date": "2026-06-17T20:10:37.576Z",
		"inquiry_id": "01KSGR1XNJMFQZMCGZHCQH2HGS",
		"inquiry_is_chat_open": true,
		"inquiry_is_marketplace": false,
		"inquiry_last_status_update": "2026-06-02T20:10:37.495Z",
		"inquiry_quotated_at": "2026-05-30T02:40:45.877Z",
		"inquiry_quotations_price_net_sum": 630000,
		"inquiry_sell_opportunity_collected": false,
		"inquiry_sell_opportunity_collected_last_note_id": "82370024",
		"inquiry_seller_agent_id": "01EM2CNQFZFSV97CMA1SSNB67F",
		"inquiry_seller_agent_name": "Luis Zlochevsky",
		"inquiry_status": "AVAILABLE",
		"last_order_seq": 10,
		"seller_name": "Luis Zlochevsky"
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	// Known fields
	if cm.InquiryID == nil || *cm.InquiryID != "01KSGR1XNJMFQZMCGZHCQH2HGS" {
		t.Errorf("InquiryID = %v", cm.InquiryID)
	}
	if cm.InquiryStatus == nil || *cm.InquiryStatus != "AVAILABLE" {
		t.Errorf("InquiryStatus = %v", cm.InquiryStatus)
	}
	if cm.InquiryDisplayID == nil || *cm.InquiryDisplayID != "31049213" {
		t.Errorf("InquiryDisplayID = %v", cm.InquiryDisplayID)
	}
	if cm.InquiryAgentID == nil || *cm.InquiryAgentID != "01EM2CNQFZFSV97CMA1SSNB67F" {
		t.Errorf("InquiryAgentID = %v", cm.InquiryAgentID)
	}
	if cm.InquiryAgentName == nil || *cm.InquiryAgentName != "Luis Zlochevsky" {
		t.Errorf("InquiryAgentName = %v", cm.InquiryAgentName)
	}
	if cm.InquiryAIEvaluation == nil || *cm.InquiryAIEvaluation != "OTHERS" {
		t.Errorf("InquiryAIEvaluation = %v", cm.InquiryAIEvaluation)
	}
	if cm.InquiryCanBindDisplayID != MetadataBoolFalse {
		t.Error("InquiryCanBindDisplayID pointer comparison failed")
	}
	if cm.InquiryCreatedAt == nil {
		t.Error("InquiryCreatedAt is nil")
	}
	if cm.InquiryExpireDate == nil {
		t.Error("InquiryExpireDate is nil")
	}
	if cm.InquirySellerAgentID == nil || *cm.InquirySellerAgentID != "01EM2CNQFZFSV97CMA1SSNB67F" {
		t.Errorf("InquirySellerAgentID = %v", cm.InquirySellerAgentID)
	}
	if cm.InquirySellerAgentName == nil || *cm.InquirySellerAgentName != "Luis Zlochevsky" {
		t.Errorf("InquirySellerAgentName = %v", cm.InquirySellerAgentName)
	}
	if cm.InquirySellOpportunityCollected != MetadataBoolFalse {
		t.Error("InquirySellOpportunityCollected pointer comparison failed")
	}
	if cm.InquiryQuotatedAt == nil {
		t.Error("InquiryQuotatedAt is nil")
	}
	if cm.InquiryLastStatusUpdate == nil {
		t.Error("InquiryLastStatusUpdate is nil")
	}
	if cm.ChatbotInitialContact == nil {
		t.Error("ChatbotInitialContact is nil")
	}
	if cm.ChatbotIsPreRegistration != MetadataBoolFalse {
		t.Error("ChatbotIsPreRegistration pointer comparison failed")
	}
	if cm.ChatbotLastState != ChatbotLastStateCreateInquiry {
		t.Error("ChatbotLastState pointer comparison failed")
	}
	if cm.ChatbotRegistrationDate == nil {
		t.Error("ChatbotRegistrationDate is nil")
	}
	if cm.InitialContactChannel != InitialContactChannelWhatsApp {
		t.Error("InitialContactChannel pointer comparison failed")
	}
	if cm.InitialContactDate == nil {
		t.Error("InitialContactDate is nil")
	}
	if cm.InquiryIsChatOpen != MetadataBoolTrue {
		t.Error("InquiryIsChatOpen pointer comparison failed")
	}
	if cm.InquiryIsMarketplace != MetadataBoolFalse {
		t.Error("InquiryIsMarketplace pointer comparison failed")
	}

	if cm.LastOrderSeq == nil || *cm.LastOrderSeq != 10 {
		t.Errorf("LastOrderSeq = %v", cm.LastOrderSeq)
	}

	// Overflow fields — all the non-known keys
	expectedOverflow := []string{
		"chatbot_name_updated",
		"has_pendencies",
		"inquiry_quotations_price_net_sum",
		"inquiry_sell_opportunity_collected_last_note_id",
		"seller_name",
	}
	if len(cm.OtherFields) != len(expectedOverflow) {
		t.Errorf("OtherFields has %d entries, want %d: %v", len(cm.OtherFields), len(expectedOverflow), cm.OtherFields)
	}
	for _, key := range expectedOverflow {
		if _, ok := cm.OtherFields[key]; !ok {
			t.Errorf("OtherFields missing key %q", key)
		}
	}

	// Spot-check some overflow values
	if cm.OtherFields["inquiry_quotations_price_net_sum"] != float64(630000) {
		t.Errorf("inquiry_quotations_price_net_sum = %v", cm.OtherFields["inquiry_quotations_price_net_sum"])
	}
	if cm.OtherFields["seller_name"] != "Luis Zlochevsky" {
		t.Errorf("seller_name = %v", cm.OtherFields["seller_name"])
	}

	// Round-trip: marshal back and compare
	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var roundTripped map[string]any
	if err := json.Unmarshal(data, &roundTripped); err != nil {
		t.Fatal(err)
	}

	var original map[string]any
	if err := json.Unmarshal([]byte(input), &original); err != nil {
		t.Fatal(err)
	}

	if len(roundTripped) != len(original) {
		t.Errorf("round-trip has %d keys, original has %d", len(roundTripped), len(original))
	}

	for key, wantVal := range original {
		gotVal, exists := roundTripped[key]
		if !exists {
			t.Errorf("round-trip missing key %q", key)
			continue
		}
		wantJSON, _ := json.Marshal(wantVal)
		gotJSON, _ := json.Marshal(gotVal)
		if string(wantJSON) != string(gotJSON) {
			t.Errorf("round-trip key %q: got %s, want %s", key, gotJSON, wantJSON)
		}
	}
}

func TestContactMetadata_SkipWelcome(t *testing.T) {
	for _, tc := range []struct {
		name string
		json string
		want MetadataBool
	}{
		{"bool true", `{"skip_welcome": true}`, MetadataBoolTrue},
		{"bool false", `{"skip_welcome": false}`, MetadataBoolFalse},
		{"string true", `{"skip_welcome": "true"}`, MetadataBoolTrue},
		{"string false", `{"skip_welcome": "false"}`, MetadataBoolFalse},
		{"null", `{"skip_welcome": null}`, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var cm ContactMetadata
			if err := json.Unmarshal([]byte(tc.json), &cm); err != nil {
				t.Fatal(err)
			}
			if cm.SkipWelcome != tc.want {
				t.Errorf("SkipWelcome = %v, want %v", cm.SkipWelcome, tc.want)
			}
			if _, exists := cm.OtherFields["skip_welcome"]; exists {
				t.Error("skip_welcome should be a known field, not in OtherFields")
			}
		})
	}
}

func TestContactMetadata_SkipWelcomeMarshal(t *testing.T) {
	cm := ContactMetadata{SkipWelcome: MetadataBoolTrue}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if m["skip_welcome"] != true {
		t.Errorf("skip_welcome = %v, want true", m["skip_welcome"])
	}
}

func TestContactMetadata_MarketingDisabledReason(t *testing.T) {
	input := `{"marketing_disabled_reason": {"error": 131, "meta_reason": "user_opt_out", "date": "2026-06-05", "extra": "keep"}}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	if cm.MarketingDisabledReason == nil {
		t.Fatal("MarketingDisabledReason is nil")
	}
	if cm.MarketingDisabledReason.Error == nil || *cm.MarketingDisabledReason.Error != 131 {
		t.Errorf("MarketingDisabledReason.Error = %v, want 131", cm.MarketingDisabledReason.Error)
	}
	if cm.MarketingDisabledReason.MetaReason == nil || *cm.MarketingDisabledReason.MetaReason != "user_opt_out" {
		t.Errorf("MarketingDisabledReason.MetaReason = %v, want user_opt_out", cm.MarketingDisabledReason.MetaReason)
	}
	if cm.MarketingDisabledReason.Date == nil || *cm.MarketingDisabledReason.Date != "2026-06-05" {
		t.Errorf("MarketingDisabledReason.Date = %v, want 2026-06-05", cm.MarketingDisabledReason.Date)
	}
	if cm.MarketingDisabledReason.OtherFields["extra"] != "keep" {
		t.Errorf("MarketingDisabledReason.OtherFields[extra] = %v, want keep", cm.MarketingDisabledReason.OtherFields["extra"])
	}
	if _, exists := cm.OtherFields["marketing_disabled_reason"]; exists {
		t.Error("marketing_disabled_reason should be a known field, not in OtherFields")
	}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	mdr, ok := m["marketing_disabled_reason"].(map[string]any)
	if !ok {
		t.Fatalf("marketing_disabled_reason = %T, want object", m["marketing_disabled_reason"])
	}
	if mdr["error"] != float64(131) {
		t.Errorf("error = %v, want 131", mdr["error"])
	}
	if mdr["meta_reason"] != "user_opt_out" {
		t.Errorf("meta_reason = %v, want user_opt_out", mdr["meta_reason"])
	}
	if mdr["extra"] != "keep" {
		t.Errorf("extra = %v, want keep", mdr["extra"])
	}
}

func TestContactMetadata_InquiryStatusInterning(t *testing.T) {
	input := `{"inquiry_status": "AVAILABLE"}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	if cm.InquiryStatus != InquiryStatusAvailable {
		t.Error("pointer comparison failed: InquiryStatus should be the interned InquiryStatusAvailable pointer")
	}

	for _, tc := range []struct {
		json     string
		sentinel *string
	}{
		{`{"inquiry_status": "CANCELLED"}`, InquiryStatusCancelled},
		{`{"inquiry_status": "DONE"}`, InquiryStatusDone},
		{`{"inquiry_status": "EXPIRED"}`, InquiryStatusExpired},
		{`{"inquiry_status": "IN_ATTENDANCE"}`, InquiryStatusInAttendance},
		{`{"inquiry_status": "ORDERED"}`, InquiryStatusOrdered},
		{`{"inquiry_status": "PENDING"}`, InquiryStatusPending},
		{`{"inquiry_status": "$del"}`, InquiryStatusDel},
	} {
		var cm2 ContactMetadata
		if err := json.Unmarshal([]byte(tc.json), &cm2); err != nil {
			t.Fatal(err)
		}
		if cm2.InquiryStatus != tc.sentinel {
			t.Errorf("pointer comparison failed for %s", *tc.sentinel)
		}
	}
}

func TestContactMetadata_InquiryStatusUnknownValue(t *testing.T) {
	input := `{"inquiry_status": "SOME_FUTURE_STATUS"}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	if cm.InquiryStatus == nil {
		t.Fatal("InquiryStatus should not be nil")
	}
	if *cm.InquiryStatus != "SOME_FUTURE_STATUS" {
		t.Errorf("InquiryStatus = %v, want SOME_FUTURE_STATUS", *cm.InquiryStatus)
	}
	if cm.InquiryStatus == InquiryStatusAvailable {
		t.Error("unknown value should not match any sentinel pointer")
	}
}

func TestContactMetadata_InquiryStatusMarshal(t *testing.T) {
	cm := ContactMetadata{InquiryStatus: InquiryStatusAvailable}

	data, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if m["inquiry_status"] != "AVAILABLE" {
		t.Errorf("inquiry_status = %v, want AVAILABLE", m["inquiry_status"])
	}
}

func TestEqualInquiryStatus(t *testing.T) {
	if !EqualInquiryStatus(InquiryStatusAvailable, InquiryStatusAvailable) {
		t.Error("same pointer should be equal")
	}

	s := "AVAILABLE"
	if !EqualInquiryStatus(&s, InquiryStatusAvailable) {
		t.Error("same value different pointer should be equal")
	}

	if EqualInquiryStatus(InquiryStatusAvailable, InquiryStatusCancelled) {
		t.Error("different values should not be equal")
	}

	if !EqualInquiryStatus(nil, nil) {
		t.Error("nil == nil should be true")
	}
	if EqualInquiryStatus(nil, InquiryStatusAvailable) {
		t.Error("nil != non-nil should be false")
	}
	if EqualInquiryStatus(InquiryStatusAvailable, nil) {
		t.Error("non-nil != nil should be false")
	}
}
