package types

import (
	"encoding/json"
	"testing"
)

func TestContactMetadataPrescription_UnmarshalFull(t *testing.T) {
	input := `{
		"id": "rx-001",
		"display_id": "RX-001",
		"doctor_id": "doc-123",
		"customer": {
			"doctor_customer_id": "dc-1",
			"account_id": "acc-1",
			"customer_id": "cust-1",
			"name": "João Silva",
			"cpf": "12345678900",
			"phone": "+5511999999999",
			"email": "joao@example.com",
			"birthdate": "1990-05-15",
			"address": {
				"country": "BR",
				"street": "Rua das Flores",
				"number": "123",
				"complement": "Apto 4",
				"city": "São Paulo",
				"uf": "SP",
				"district": "Centro",
				"zip": "01001-000"
			}
		},
		"file_url": "https://example.com/rx.pdf",
		"preseller_id": "ps-1",
		"preseller_name": "Maria",
		"observation": "Uso contínuo",
		"observation_author": "Dr. Santos"
	}`

	var p ContactMetadataPrescription
	if err := json.Unmarshal([]byte(input), &p); err != nil {
		t.Fatal(err)
	}

	if p.ID == nil || *p.ID != "rx-001" {
		t.Errorf("ID = %v", p.ID)
	}
	if p.DisplayID == nil || *p.DisplayID != "RX-001" {
		t.Errorf("DisplayID = %v", p.DisplayID)
	}
	if p.DoctorID == nil || *p.DoctorID != "doc-123" {
		t.Errorf("DoctorID = %v", p.DoctorID)
	}
	if p.FileURL == nil || *p.FileURL != "https://example.com/rx.pdf" {
		t.Errorf("FileURL = %v", p.FileURL)
	}
	if p.PresellerID == nil || *p.PresellerID != "ps-1" {
		t.Errorf("PresellerID = %v", p.PresellerID)
	}
	if p.PresellerName == nil || *p.PresellerName != "Maria" {
		t.Errorf("PresellerName = %v", p.PresellerName)
	}
	if p.Observation == nil || *p.Observation != "Uso contínuo" {
		t.Errorf("Observation = %v", p.Observation)
	}
	if p.ObservationAuthor == nil || *p.ObservationAuthor != "Dr. Santos" {
		t.Errorf("ObservationAuthor = %v", p.ObservationAuthor)
	}
	if p.OtherFields != nil {
		t.Errorf("OtherFields = %v, want nil", p.OtherFields)
	}

	// Customer
	if p.Customer == nil {
		t.Fatal("Customer is nil")
	}
	c := p.Customer
	if c.DoctorCustomerID == nil || *c.DoctorCustomerID != "dc-1" {
		t.Errorf("Customer.DoctorCustomerID = %v", c.DoctorCustomerID)
	}
	if c.Name == nil || *c.Name != "João Silva" {
		t.Errorf("Customer.Name = %v", c.Name)
	}
	if c.CPF == nil || *c.CPF != "12345678900" {
		t.Errorf("Customer.CPF = %v", c.CPF)
	}
	if c.Birthdate == nil || *c.Birthdate != "1990-05-15" {
		t.Errorf("Customer.Birthdate = %v", c.Birthdate)
	}
	if c.OtherFields != nil {
		t.Errorf("Customer.OtherFields = %v, want nil", c.OtherFields)
	}

	// Address
	if c.Address == nil {
		t.Fatal("Customer.Address is nil")
	}
	a := c.Address
	if a.Country == nil || *a.Country != "BR" {
		t.Errorf("Address.Country = %v", a.Country)
	}
	if a.Street == nil || *a.Street != "Rua das Flores" {
		t.Errorf("Address.Street = %v", a.Street)
	}
	if a.Number == nil || *a.Number != "123" {
		t.Errorf("Address.Number = %v", a.Number)
	}
	if a.Complement == nil || *a.Complement != "Apto 4" {
		t.Errorf("Address.Complement = %v", a.Complement)
	}
	if a.City == nil || *a.City != "São Paulo" {
		t.Errorf("Address.City = %v", a.City)
	}
	if a.UF == nil || *a.UF != "SP" {
		t.Errorf("Address.UF = %v", a.UF)
	}
	if a.District == nil || *a.District != "Centro" {
		t.Errorf("Address.District = %v", a.District)
	}
	if a.Zip == nil || *a.Zip != "01001-000" {
		t.Errorf("Address.Zip = %v", a.Zip)
	}
	if a.OtherFields != nil {
		t.Errorf("Address.OtherFields = %v, want nil", a.OtherFields)
	}
}

func TestContactMetadataPrescription_OverflowAtAllLevels(t *testing.T) {
	input := `{
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
	}`

	var p ContactMetadataPrescription
	if err := json.Unmarshal([]byte(input), &p); err != nil {
		t.Fatal(err)
	}

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

func TestContactMetadataPrescription_RoundTrip(t *testing.T) {
	input := `{
		"id": "rx-001",
		"doctor_id": "doc-1",
		"custom_field": "hello",
		"customer": {
			"name": "Test",
			"cpf": "000",
			"loyalty_tier": "gold",
			"address": {
				"city": "SP",
				"uf": "SP",
				"region_code": "SE"
			}
		},
		"file_url": "https://example.com/rx.pdf"
	}`

	var p ContactMetadataPrescription
	if err := json.Unmarshal([]byte(input), &p); err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(p)
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

func TestContactMetadataPrescription_MarshalEmpty(t *testing.T) {
	p := ContactMetadataPrescription{}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "{}" {
		t.Errorf("empty marshal = %s, want {}", string(data))
	}
}

func TestContactMetadata_NestedPrescription(t *testing.T) {
	input := `{
		"inquiry_id": "inq-1",
		"prescription": {
			"id": "rx-001",
			"customer": {
				"name": "Test",
				"address": {"city": "SP"}
			}
		}
	}`

	var cm ContactMetadata
	if err := json.Unmarshal([]byte(input), &cm); err != nil {
		t.Fatal(err)
	}

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
