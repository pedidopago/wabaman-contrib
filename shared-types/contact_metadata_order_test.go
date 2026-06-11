package types

import (
	"encoding/json"
	"testing"
)

func TestContactMetadataOrder_UnmarshalKnownStatus(t *testing.T) {
	input := `{"status": "created", "prescription": {"is_medical": true, "retained": false}}`

	var o ContactMetadataOrder
	if err := json.Unmarshal([]byte(input), &o); err != nil {
		t.Fatal(err)
	}

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

func TestContactMetadataOrder_StatusInterning(t *testing.T) {
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
		var o ContactMetadataOrder
		if err := json.Unmarshal([]byte(tc.json), &o); err != nil {
			t.Fatal(err)
		}
		if o.Status != tc.sentinel {
			t.Errorf("pointer comparison failed for %s", *tc.sentinel)
		}
	}
}

func TestContactMetadataOrder_UnknownStatus(t *testing.T) {
	input := `{"status": "some_future_status"}`

	var o ContactMetadataOrder
	if err := json.Unmarshal([]byte(input), &o); err != nil {
		t.Fatal(err)
	}

	if o.Status == nil || *o.Status != "some_future_status" {
		t.Errorf("Status = %v, want some_future_status", o.Status)
	}
}

func TestContactMetadataOrder_Overflow(t *testing.T) {
	input := `{"status": "created", "tracking_code": "BR123", "weight": 1.5}`

	var o ContactMetadataOrder
	if err := json.Unmarshal([]byte(input), &o); err != nil {
		t.Fatal(err)
	}

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

func TestContactMetadataOrder_MarshalRoundTrip(t *testing.T) {
	input := `{
		"status": "approved",
		"prescription": {"is_medical": true, "retained": false},
		"custom_field": "value"
	}`

	var o ContactMetadataOrder
	if err := json.Unmarshal([]byte(input), &o); err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(o)
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
}

func TestContactMetadataOrder_MarshalEmpty(t *testing.T) {
	o := ContactMetadataOrder{}

	data, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "{}" {
		t.Errorf("empty marshal = %s, want {}", string(data))
	}
}

func TestContactMetadata_NestedOrder(t *testing.T) {
	input := `{
		"inquiry_id": "inq-1",
		"inquiry_status": "AVAILABLE",
		"order": {
			"status": "created",
			"prescription": {"is_medical": false, "retained": true},
			"delivery_date": "2026-06-15"
		}
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

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
}

func TestEqualOrderStatus(t *testing.T) {
	if !EqualOrderStatus(OrderStatusCreated, OrderStatusCreated) {
		t.Error("same pointer should be equal")
	}

	s := "created"
	if !EqualOrderStatus(&s, OrderStatusCreated) {
		t.Error("same value different pointer should be equal")
	}

	if EqualOrderStatus(OrderStatusCreated, OrderStatusCancelled) {
		t.Error("different values should not be equal")
	}

	if !EqualOrderStatus(nil, nil) {
		t.Error("nil == nil should be true")
	}
	if EqualOrderStatus(nil, OrderStatusCreated) {
		t.Error("nil != non-nil should be false")
	}
}
