package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gabstv/mmm"
	"github.com/pedidopago/zajson"
)

const testArenaSize = 64 * 1024

func zmarshal(t *testing.T, v interface{ ZMarshalJSON(w *zajson.Writer) error }) []byte {
	t.Helper()
	arena := mmm.NewArena(testArenaSize)
	defer mmm.DestroyArena(&arena)
	w := zajson.NewWriter(arena)
	if err := v.ZMarshalJSON(w); err != nil {
		t.Fatalf("ZMarshalJSON: %v", err)
	}
	out := make([]byte, len(w.Bytes()))
	copy(out, w.Bytes())
	return out
}

func zunmarshal[T interface {
	*E
	ZUnmarshalJSON(r *zajson.Reader) error
}, E any](t *testing.T, data []byte) *E {
	t.Helper()
	arena := mmm.NewArena(testArenaSize)
	defer mmm.DestroyArena(&arena)
	r := zajson.NewReader(arena, data)
	v := new(E)
	if err := T(v).ZUnmarshalJSON(r); err != nil {
		t.Fatalf("ZUnmarshalJSON: %v", err)
	}
	return v
}

// zRoundTrip marshals v via Z, unmarshals into a new instance via Z, then
// JSON-marshals both and compares. It returns the Z-unmarshaled copy.
func zRoundTrip[T interface {
	*E
	ZMarshalJSON(w *zajson.Writer) error
	ZUnmarshalJSON(r *zajson.Reader) error
}, E any](t *testing.T, v *E) *E {
	t.Helper()
	data := zmarshal(t, T(v))
	got := zunmarshal[T](t, data)
	return got
}

func assertJSONEqual(t *testing.T, label string, a, b any) {
	t.Helper()
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if string(ja) != string(jb) {
		t.Errorf("%s: got %s, want %s", label, ja, jb)
	}
}

// ---------------------------------------------------------------------------
// ContactMetadata Z round-trip tests
// ---------------------------------------------------------------------------

func TestZ_ContactMetadata_UnmarshalKnownFields(t *testing.T) {
	input := []byte(`{
		"inquiry_id": "inq-123",
		"inquiry_status": "open",
		"chatbot_disabled": true,
		"account_id": "acc-456"
	}`)

	cm := zunmarshal[*ContactMetadata](t, input)

	if cm.InquiryID == nil || *cm.InquiryID != "inq-123" {
		t.Errorf("InquiryID = %v, want inq-123", cm.InquiryID)
	}
	if cm.InquiryStatus == nil || *cm.InquiryStatus != "open" {
		t.Errorf("InquiryStatus = %v, want open", cm.InquiryStatus)
	}
	if cm.ChatbotDisabled != MetadataBoolTrue {
		t.Error("ChatbotDisabled pointer comparison failed")
	}
	if cm.AccountID == nil || *cm.AccountID != "acc-456" {
		t.Errorf("AccountID = %v, want acc-456", cm.AccountID)
	}
	if cm.OtherFields != nil {
		t.Errorf("OtherFields = %v, want nil", cm.OtherFields)
	}
}

func TestZ_ContactMetadata_UnmarshalOnlyUnknownFields(t *testing.T) {
	input := []byte(`{"custom_field": "value", "priority": 5}`)

	cm := zunmarshal[*ContactMetadata](t, input)

	if cm.InquiryID != nil {
		t.Errorf("InquiryID = %v, want nil", cm.InquiryID)
	}
	if len(cm.OtherFields) != 2 {
		t.Fatalf("OtherFields has %d entries, want 2", len(cm.OtherFields))
	}
	if cm.OtherFields["custom_field"] != "value" {
		t.Errorf("OtherFields[custom_field] = %v, want value", cm.OtherFields["custom_field"])
	}
	if fmt.Sprint(cm.OtherFields["priority"]) != "5" {
		t.Errorf("OtherFields[priority] = %v (%T), want 5", cm.OtherFields["priority"], cm.OtherFields["priority"])
	}
}

func TestZ_ContactMetadata_UnmarshalMixed(t *testing.T) {
	input := []byte(`{
		"inquiry_id": "inq-1",
		"inquiry_agent_name": "Agent Smith",
		"custom_tag": "vip",
		"score": 99.5
	}`)

	cm := zunmarshal[*ContactMetadata](t, input)

	if cm.InquiryID == nil || *cm.InquiryID != "inq-1" {
		t.Errorf("InquiryID = %v, want inq-1", cm.InquiryID)
	}
	if cm.InquiryAgentName == nil || *cm.InquiryAgentName != "Agent Smith" {
		t.Errorf("InquiryAgentName = %v, want Agent Smith", cm.InquiryAgentName)
	}
	if cm.OtherFields["custom_tag"] != "vip" {
		t.Errorf("OtherFields[custom_tag] = %v, want vip", cm.OtherFields["custom_tag"])
	}
	if fmt.Sprint(cm.OtherFields["score"]) != "99.5" {
		t.Errorf("OtherFields[score] = %v (%T), want 99.5", cm.OtherFields["score"], cm.OtherFields["score"])
	}
}

func TestZ_ContactMetadata_UnmarshalTimestamps(t *testing.T) {
	input := []byte(`{
		"inquiry_created_at": "2025-01-15T10:30:00Z",
		"inquiry_expire_date": "2025-02-15T10:30:00Z"
	}`)

	cm := zunmarshal[*ContactMetadata](t, input)

	want := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	if cm.InquiryCreatedAt == nil || !cm.InquiryCreatedAt.Equal(want) {
		t.Errorf("InquiryCreatedAt = %v, want %v", cm.InquiryCreatedAt, want)
	}

	wantExpire := time.Date(2025, 2, 15, 10, 30, 0, 0, time.UTC)
	if cm.InquiryExpireDate == nil || !cm.InquiryExpireDate.Equal(wantExpire) {
		t.Errorf("InquiryExpireDate = %v, want %v", cm.InquiryExpireDate, wantExpire)
	}
}

func TestZ_ContactMetadata_UnmarshalEmpty(t *testing.T) {
	cm := zunmarshal[*ContactMetadata](t, []byte(`{}`))

	if cm.InquiryID != nil {
		t.Error("InquiryID should be nil")
	}
	if cm.OtherFields != nil {
		t.Errorf("OtherFields = %v, want nil", cm.OtherFields)
	}
}

func TestZ_ContactMetadata_MarshalKnownFields(t *testing.T) {
	cm := ContactMetadata{
		InquiryID:       ptr("inq-1"),
		InquiryStatus:   InquiryStatusPending,
		ChatbotDisabled: MetadataBoolTrue,
	}

	data := zmarshal(t, &cm)

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Z output is not valid JSON: %v\nraw: %s", err, data)
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

func TestZ_ContactMetadata_MarshalEmpty(t *testing.T) {
	cm := ContactMetadata{}
	data := zmarshal(t, &cm)

	if string(data) != "{}" {
		t.Errorf("empty marshal = %s, want {}", data)
	}
}

func TestZ_ContactMetadata_MarshalOnlyOtherFields(t *testing.T) {
	cm := ContactMetadata{
		OtherFields: map[string]any{
			"custom": "value",
			"count":  42,
		},
	}

	data := zmarshal(t, &cm)

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Z output is not valid JSON: %v\nraw: %s", err, data)
	}

	if m["custom"] != "value" {
		t.Errorf("custom = %v, want value", m["custom"])
	}
	if m["count"] != float64(42) {
		t.Errorf("count = %v, want 42", m["count"])
	}
}

func TestZ_ContactMetadata_InquiryStatusInterning(t *testing.T) {
	for _, tc := range []struct {
		json     string
		sentinel *string
	}{
		{`{"inquiry_status": "AVAILABLE"}`, InquiryStatusAvailable},
		{`{"inquiry_status": "CANCELLED"}`, InquiryStatusCancelled},
		{`{"inquiry_status": "DONE"}`, InquiryStatusDone},
		{`{"inquiry_status": "EXPIRED"}`, InquiryStatusExpired},
		{`{"inquiry_status": "IN_ATTENDANCE"}`, InquiryStatusInAttendance},
		{`{"inquiry_status": "ORDERED"}`, InquiryStatusOrdered},
		{`{"inquiry_status": "PENDING"}`, InquiryStatusPending},
		{`{"inquiry_status": "$del"}`, InquiryStatusDel},
	} {
		cm := zunmarshal[*ContactMetadata](t, []byte(tc.json))
		if cm.InquiryStatus != tc.sentinel {
			t.Errorf("pointer comparison failed for %s", *tc.sentinel)
		}
	}
}

func TestZ_ContactMetadata_InquiryStatusUnknown(t *testing.T) {
	cm := zunmarshal[*ContactMetadata](t, []byte(`{"inquiry_status": "SOME_FUTURE_STATUS"}`))

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

func TestZ_ContactMetadata_RoundTrip(t *testing.T) {
	input := `{
		"inquiry_id": "inq-1",
		"inquiry_status": "closed",
		"chatbot_disabled": false,
		"inquiry_created_at": "2025-06-01T12:00:00Z",
		"custom_field": "hello",
		"nested": {"a": 1}
	}`

	// JSON unmarshal → Z marshal → JSON unmarshal → compare
	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	zBytes := zmarshal(t, &cm)

	var roundTripped map[string]any
	if err := json.Unmarshal(zBytes, &roundTripped); err != nil {
		t.Fatalf("Z output is not valid JSON: %v\nraw: %s", err, zBytes)
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

func TestZ_ContactMetadata_FullRoundTrip(t *testing.T) {
	// Z unmarshal → Z marshal → JSON unmarshal → compare with original
	input := []byte(`{
		"inquiry_id": "inq-1",
		"inquiry_status": "AVAILABLE",
		"inquiry_agent_id": "agent-1",
		"inquiry_can_bind_display_id": true,
		"inquiry_created_at": "2025-01-15T10:30:00Z",
		"chatbot_disabled": false,
		"chatbot_last_state": "chatbot_create_inquiry",
		"initial_contact_channel": "whatsapp",
		"customer_name": "João Silva",
		"customer_document_country": "BR",
		"last_order_seq": 5,
		"custom_overflow": "hello"
	}`)

	cm := zunmarshal[*ContactMetadata](t, input)
	zBytes := zmarshal(t, cm)

	var roundTripped map[string]any
	if err := json.Unmarshal(zBytes, &roundTripped); err != nil {
		t.Fatalf("Z output is not valid JSON: %v\nraw: %s", err, zBytes)
	}

	var original map[string]any
	if err := json.Unmarshal(input, &original); err != nil {
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

func TestZ_ContactMetadata_RealWorldPayload(t *testing.T) {
	input := []byte(`{
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
		"inquiry_quotations_price_net_sum": 630000,
		"inquiry_quotated_at": "2026-05-30T02:40:45.877Z",
		"inquiry_sell_opportunity_collected": false,
		"inquiry_sell_opportunity_collected_last_note_id": "82370024",
		"inquiry_seller_agent_id": "01EM2CNQFZFSV97CMA1SSNB67F",
		"inquiry_seller_agent_name": "Luis Zlochevsky",
		"inquiry_status": "AVAILABLE",
		"last_order_seq": 10,
		"seller_name": "Luis Zlochevsky"
	}`)

	cm := zunmarshal[*ContactMetadata](t, input)

	// Known fields
	if cm.InquiryID == nil || *cm.InquiryID != "01KSGR1XNJMFQZMCGZHCQH2HGS" {
		t.Errorf("InquiryID = %v", cm.InquiryID)
	}
	if cm.InquiryStatus != InquiryStatusAvailable {
		t.Error("InquiryStatus pointer comparison failed")
	}
	if cm.InquiryCanBindDisplayID != MetadataBoolFalse {
		t.Error("InquiryCanBindDisplayID pointer comparison failed")
	}
	if cm.InquiryIsChatOpen != MetadataBoolTrue {
		t.Error("InquiryIsChatOpen pointer comparison failed")
	}
	if cm.ChatbotLastState != ChatbotLastStateCreateInquiry {
		t.Error("ChatbotLastState pointer comparison failed")
	}
	if cm.InitialContactChannel != InitialContactChannelWhatsApp {
		t.Error("InitialContactChannel pointer comparison failed")
	}
	if cm.LastOrderSeq == nil || *cm.LastOrderSeq != 10 {
		t.Errorf("LastOrderSeq = %v", cm.LastOrderSeq)
	}

	// Overflow
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

	// Z round-trip
	zBytes := zmarshal(t, cm)
	var roundTripped map[string]any
	if err := json.Unmarshal(zBytes, &roundTripped); err != nil {
		t.Fatalf("Z output is not valid JSON: %v\nraw: %s", err, zBytes)
	}
	var original map[string]any
	if err := json.Unmarshal(input, &original); err != nil {
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

// ---------------------------------------------------------------------------
// ContactMetadataOrder Z tests
// ---------------------------------------------------------------------------

func TestZ_ContactMetadataOrder_UnmarshalKnownStatus(t *testing.T) {
	input := []byte(`{"status": "created", "prescription": {"is_medical": true, "retained": false}}`)

	o := zunmarshal[*ContactMetadataOrder](t, input)

	if o.Status != OrderStatusCreated {
		t.Error("Status pointer comparison failed")
	}
	if o.Prescription == nil {
		t.Fatal("Prescription is nil")
	}
	if o.Prescription.IsMedical != MetadataBoolTrue {
		t.Error("Prescription.IsMedical pointer comparison failed")
	}
	if o.Prescription.Retained != MetadataBoolFalse {
		t.Error("Prescription.Retained pointer comparison failed")
	}
	if o.OtherFields != nil {
		t.Errorf("OtherFields = %v, want nil", o.OtherFields)
	}
}

func TestZ_ContactMetadataOrder_StatusInterning(t *testing.T) {
	for _, tc := range []struct {
		json     string
		sentinel OrderStatus
	}{
		{`{"status": "assembling"}`, OrderStatusAssembling},
		{`{"status": "approved"}`, OrderStatusApproved},
		{`{"status": "created"}`, OrderStatusCreated},
		{`{"status": "finalized"}`, OrderStatusFinalized},
		{`{"status": "prescription_collection"}`, OrderStatusPrescriptionCollection},
		{`{"status": "cancelled"}`, OrderStatusCancelled},
		{`{"status": "$del"}`, OrderStatusDel},
	} {
		o := zunmarshal[*ContactMetadataOrder](t, []byte(tc.json))
		if o.Status != tc.sentinel {
			t.Errorf("pointer comparison failed for %s", *tc.sentinel)
		}
	}
}

func TestZ_ContactMetadataOrder_Overflow(t *testing.T) {
	input := []byte(`{"status": "created", "tracking_code": "BR123", "weight": 1.5}`)

	o := zunmarshal[*ContactMetadataOrder](t, input)

	if o.Status != OrderStatusCreated {
		t.Error("Status pointer comparison failed")
	}
	if len(o.OtherFields) != 2 {
		t.Fatalf("OtherFields has %d entries, want 2", len(o.OtherFields))
	}
	if o.OtherFields["tracking_code"] != "BR123" {
		t.Errorf("tracking_code = %v", o.OtherFields["tracking_code"])
	}
}

func TestZ_ContactMetadataOrder_RoundTrip(t *testing.T) {
	input := []byte(`{
		"status": "approved",
		"prescription": {"is_medical": true, "retained": false},
		"custom_field": "value"
	}`)

	o := zunmarshal[*ContactMetadataOrder](t, input)
	zBytes := zmarshal(t, o)

	var roundTripped, original map[string]any
	if err := json.Unmarshal(zBytes, &roundTripped); err != nil {
		t.Fatalf("Z output is not valid JSON: %v\nraw: %s", err, zBytes)
	}
	if err := json.Unmarshal(input, &original); err != nil {
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
}

// ---------------------------------------------------------------------------
// ContactMetadataPrescription Z tests
// ---------------------------------------------------------------------------

func TestZ_ContactMetadataPrescription_FullRoundTrip(t *testing.T) {
	input := []byte(`{
		"id": "rx-001",
		"display_id": "RX-001",
		"doctor_id": "doc-123",
		"customer": {
			"doctor_customer_id": "dc-1",
			"name": "João Silva",
			"cpf": "12345678900",
			"address": {
				"country": "BR",
				"street": "Rua das Flores",
				"number": "123",
				"city": "São Paulo",
				"uf": "SP",
				"district": "Centro",
				"zip": "01001-000"
			}
		},
		"file_url": "https://example.com/rx.pdf",
		"extra_top": "top_val"
	}`)

	p := zunmarshal[*ContactMetadataPrescription](t, input)

	if p.ID == nil || *p.ID != "rx-001" {
		t.Errorf("ID = %v", p.ID)
	}
	if p.Customer == nil {
		t.Fatal("Customer is nil")
	}
	if p.Customer.Name == nil || *p.Customer.Name != "João Silva" {
		t.Errorf("Customer.Name = %v", p.Customer.Name)
	}
	if p.Customer.Address == nil {
		t.Fatal("Customer.Address is nil")
	}
	if p.Customer.Address.Country == nil || *p.Customer.Address.Country != "BR" {
		t.Errorf("Address.Country = %v", p.Customer.Address.Country)
	}
	if p.OtherFields["extra_top"] != "top_val" {
		t.Errorf("OtherFields[extra_top] = %v", p.OtherFields["extra_top"])
	}

	// Z round-trip
	zBytes := zmarshal(t, p)
	var roundTripped, original map[string]any
	if err := json.Unmarshal(zBytes, &roundTripped); err != nil {
		t.Fatalf("Z output is not valid JSON: %v\nraw: %s", err, zBytes)
	}
	if err := json.Unmarshal(input, &original); err != nil {
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

func TestZ_ContactMetadataPrescription_OverflowAtAllLevels(t *testing.T) {
	input := []byte(`{
		"id": "rx-001",
		"extra_top": "top_val",
		"customer": {
			"name": "Test",
			"extra_customer": "cust_val",
			"address": {
				"city": "SP",
				"extra_addr": "addr_val"
			}
		}
	}`)

	p := zunmarshal[*ContactMetadataPrescription](t, input)

	if p.OtherFields["extra_top"] != "top_val" {
		t.Errorf("OtherFields[extra_top] = %v", p.OtherFields["extra_top"])
	}
	if p.Customer.OtherFields["extra_customer"] != "cust_val" {
		t.Errorf("Customer.OtherFields[extra_customer] = %v", p.Customer.OtherFields["extra_customer"])
	}
	if p.Customer.Address.OtherFields["extra_addr"] != "addr_val" {
		t.Errorf("Address.OtherFields[extra_addr] = %v", p.Customer.Address.OtherFields["extra_addr"])
	}
}

// ---------------------------------------------------------------------------
// Nested struct Z tests via ContactMetadata parent
// ---------------------------------------------------------------------------

func TestZ_ContactMetadata_NestedOrder(t *testing.T) {
	input := []byte(`{
		"inquiry_id": "inq-1",
		"inquiry_status": "AVAILABLE",
		"order": {
			"status": "created",
			"prescription": {"is_medical": false, "retained": true},
			"delivery_date": "2026-06-15"
		}
	}`)

	cm := zunmarshal[*ContactMetadata](t, input)

	if cm.InquiryStatus != InquiryStatusAvailable {
		t.Error("InquiryStatus pointer comparison failed")
	}
	if cm.Order == nil {
		t.Fatal("Order is nil")
	}
	if cm.Order.Status != OrderStatusCreated {
		t.Error("Order.Status pointer comparison failed")
	}
	if cm.Order.Prescription == nil {
		t.Fatal("Order.Prescription is nil")
	}
	if cm.Order.Prescription.IsMedical != MetadataBoolFalse {
		t.Error("Order.Prescription.IsMedical pointer comparison failed")
	}
	if cm.Order.Prescription.Retained != MetadataBoolTrue {
		t.Error("Order.Prescription.Retained pointer comparison failed")
	}
	if cm.Order.OtherFields["delivery_date"] != "2026-06-15" {
		t.Errorf("Order.OtherFields[delivery_date] = %v", cm.Order.OtherFields["delivery_date"])
	}

	// Round-trip
	zBytes := zmarshal(t, cm)
	var roundTripped, original map[string]any
	if err := json.Unmarshal(zBytes, &roundTripped); err != nil {
		t.Fatalf("Z output is not valid JSON: %v\nraw: %s", err, zBytes)
	}
	if err := json.Unmarshal(input, &original); err != nil {
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
}

func TestZ_ContactMetadata_NestedPrescription(t *testing.T) {
	input := []byte(`{
		"inquiry_id": "inq-1",
		"prescription": {
			"id": "rx-001",
			"customer": {
				"name": "Test",
				"address": {"city": "SP"}
			}
		}
	}`)

	cm := zunmarshal[*ContactMetadata](t, input)

	if cm.Prescription == nil {
		t.Fatal("Prescription is nil")
	}
	if cm.Prescription.ID == nil || *cm.Prescription.ID != "rx-001" {
		t.Errorf("Prescription.ID = %v", cm.Prescription.ID)
	}
	if cm.Prescription.Customer == nil {
		t.Fatal("Prescription.Customer is nil")
	}
	if cm.Prescription.Customer.Name == nil || *cm.Prescription.Customer.Name != "Test" {
		t.Errorf("Prescription.Customer.Name = %v", cm.Prescription.Customer.Name)
	}
	if cm.Prescription.Customer.Address == nil {
		t.Fatal("Prescription.Customer.Address is nil")
	}
	if cm.Prescription.Customer.Address.City == nil || *cm.Prescription.Customer.Address.City != "SP" {
		t.Errorf("Prescription.Customer.Address.City = %v", cm.Prescription.Customer.Address.City)
	}
}

// ---------------------------------------------------------------------------
// Cross-path compatibility: JSON unmarshal → Z marshal and vice versa
// ---------------------------------------------------------------------------

func TestZ_CrossPath_JSONUnmarshalThenZMarshal(t *testing.T) {
	input := `{
		"inquiry_id": "inq-1",
		"inquiry_status": "AVAILABLE",
		"chatbot_disabled": true,
		"custom_field": "hello"
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

	zBytes := zmarshal(t, &cm)

	var got map[string]any
	if err := json.Unmarshal(zBytes, &got); err != nil {
		t.Fatalf("Z output is not valid JSON: %v\nraw: %s", err, zBytes)
	}

	var want map[string]any
	if err := json.Unmarshal([]byte(input), &want); err != nil {
		t.Fatal(err)
	}

	for key, wantVal := range want {
		gotVal, exists := got[key]
		if !exists {
			t.Errorf("missing key %q", key)
			continue
		}
		wantJSON, _ := json.Marshal(wantVal)
		gotJSON, _ := json.Marshal(gotVal)
		if string(wantJSON) != string(gotJSON) {
			t.Errorf("key %q: got %s, want %s", key, gotJSON, wantJSON)
		}
	}
}

func TestZ_CrossPath_ZUnmarshalThenJSONMarshal(t *testing.T) {
	input := []byte(`{
		"inquiry_id": "inq-1",
		"inquiry_status": "AVAILABLE",
		"chatbot_disabled": true,
		"custom_field": "hello"
	}`)

	cm := zunmarshal[*ContactMetadata](t, input)

	jsonBytes, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(jsonBytes, &got); err != nil {
		t.Fatal(err)
	}

	var want map[string]any
	if err := json.Unmarshal(input, &want); err != nil {
		t.Fatal(err)
	}

	for key, wantVal := range want {
		gotVal, exists := got[key]
		if !exists {
			t.Errorf("missing key %q", key)
			continue
		}
		wantJSON, _ := json.Marshal(wantVal)
		gotJSON, _ := json.Marshal(gotVal)
		if string(wantJSON) != string(gotJSON) {
			t.Errorf("key %q: got %s, want %s", key, gotJSON, wantJSON)
		}
	}
}

// ---------------------------------------------------------------------------
// MetadataBool multi-format Z test
// ---------------------------------------------------------------------------

func TestZ_MetadataBool_Formats(t *testing.T) {
	for _, tc := range []struct {
		name     string
		json     string
		want     MetadataBool
	}{
		{"bool true", `{"chatbot_disabled": true}`, MetadataBoolTrue},
		{"bool false", `{"chatbot_disabled": false}`, MetadataBoolFalse},
		{"string true", `{"chatbot_disabled": "true"}`, MetadataBoolTrue},
		{"string false", `{"chatbot_disabled": "false"}`, MetadataBoolFalse},
		{"null", `{"chatbot_disabled": null}`, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cm := zunmarshal[*ContactMetadata](t, []byte(tc.json))
			if cm.ChatbotDisabled != tc.want {
				t.Errorf("ChatbotDisabled = %v, want %v", cm.ChatbotDisabled, tc.want)
			}
		})
	}
}
