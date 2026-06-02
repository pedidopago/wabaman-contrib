package types

import "encoding/json"

type ContactMetadataPrescription struct {
	ID                *string               `json:"id,omitzero"`
	DisplayID         *string               `json:"display_id,omitzero"`
	DoctorID          *string               `json:"doctor_id,omitzero"`
	Customer          *PrescriptionCustomer `json:"customer,omitzero"`
	FileURL           *string               `json:"file_url,omitzero"`
	PresellerID       *string               `json:"preseller_id,omitzero"`
	PresellerName     *string               `json:"preseller_name,omitzero"`
	Observation       *string               `json:"observation,omitzero"`
	ObservationAuthor *string               `json:"observation_author,omitzero"`
	OtherFields       map[string]any        `json:"-"`
}

var contactMetadataPrescriptionKnownKeys = map[string]struct{}{
	"id":                 {},
	"display_id":         {},
	"doctor_id":          {},
	"customer":           {},
	"file_url":           {},
	"preseller_id":       {},
	"preseller_name":     {},
	"observation":        {},
	"observation_author": {},
}

func (p ContactMetadataPrescription) MarshalJSON() ([]byte, error) {
	type alias ContactMetadataPrescription
	knownBytes, err := json.Marshal(alias(p))
	if err != nil {
		return nil, err
	}

	if len(p.OtherFields) == 0 {
		return knownBytes, nil
	}

	var merged map[string]json.RawMessage
	if err := json.Unmarshal(knownBytes, &merged); err != nil {
		return nil, err
	}

	for k, v := range p.OtherFields {
		if _, known := contactMetadataPrescriptionKnownKeys[k]; known {
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

func (p *ContactMetadataPrescription) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	known := map[string]any{
		"id":                 &p.ID,
		"display_id":         &p.DisplayID,
		"doctor_id":          &p.DoctorID,
		"customer":           &p.Customer,
		"file_url":           &p.FileURL,
		"preseller_id":       &p.PresellerID,
		"preseller_name":     &p.PresellerName,
		"observation":        &p.Observation,
		"observation_author": &p.ObservationAuthor,
	}

	for key, dst := range known {
		if v, ok := raw[key]; ok {
			if err := json.Unmarshal(v, dst); err != nil {
				return err
			}
		}
	}

	for k, v := range raw {
		if _, isKnown := contactMetadataPrescriptionKnownKeys[k]; isKnown {
			continue
		}
		if p.OtherFields == nil {
			p.OtherFields = make(map[string]any)
		}
		var decoded any
		if err := json.Unmarshal(v, &decoded); err != nil {
			return err
		}
		p.OtherFields[k] = decoded
	}

	return nil
}

type PrescriptionCustomer struct {
	DoctorCustomerID *string                      `json:"doctor_customer_id,omitzero"`
	AccountID        *string                      `json:"account_id,omitzero"`
	CustomerID       *string                      `json:"customer_id,omitzero"`
	Name             *string                      `json:"name,omitzero"`
	CPF              *string                      `json:"cpf,omitzero"`
	Phone            *string                      `json:"phone,omitzero"`
	Email            *string                      `json:"email,omitzero"`
	Birthdate        *string                      `json:"birthdate,omitzero"`
	Address          *PrescriptionCustomerAddress `json:"address,omitzero"`
	OtherFields      map[string]any               `json:"-"`
}

var prescriptionCustomerKnownKeys = map[string]struct{}{
	"doctor_customer_id": {},
	"account_id":         {},
	"customer_id":        {},
	"name":               {},
	"cpf":                {},
	"phone":              {},
	"email":              {},
	"birthdate":          {},
	"address":            {},
}

func (c PrescriptionCustomer) MarshalJSON() ([]byte, error) {
	type alias PrescriptionCustomer
	knownBytes, err := json.Marshal(alias(c))
	if err != nil {
		return nil, err
	}

	if len(c.OtherFields) == 0 {
		return knownBytes, nil
	}

	var merged map[string]json.RawMessage
	if err := json.Unmarshal(knownBytes, &merged); err != nil {
		return nil, err
	}

	for k, v := range c.OtherFields {
		if _, known := prescriptionCustomerKnownKeys[k]; known {
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

func (c *PrescriptionCustomer) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	known := map[string]any{
		"doctor_customer_id": &c.DoctorCustomerID,
		"account_id":         &c.AccountID,
		"customer_id":        &c.CustomerID,
		"name":               &c.Name,
		"cpf":                &c.CPF,
		"phone":              &c.Phone,
		"email":              &c.Email,
		"birthdate":          &c.Birthdate,
		"address":            &c.Address,
	}

	for key, dst := range known {
		if v, ok := raw[key]; ok {
			if err := json.Unmarshal(v, dst); err != nil {
				return err
			}
		}
	}

	for k, v := range raw {
		if _, isKnown := prescriptionCustomerKnownKeys[k]; isKnown {
			continue
		}
		if c.OtherFields == nil {
			c.OtherFields = make(map[string]any)
		}
		var decoded any
		if err := json.Unmarshal(v, &decoded); err != nil {
			return err
		}
		c.OtherFields[k] = decoded
	}

	return nil
}

type PrescriptionCustomerAddress struct {
	Country     *string        `json:"country,omitzero"`
	Street      *string        `json:"street,omitzero"`
	Number      *string        `json:"number,omitzero"`
	Complement  *string        `json:"complement,omitzero"`
	City        *string        `json:"city,omitzero"`
	UF          *string        `json:"uf,omitzero"`
	District    *string        `json:"district,omitzero"`
	Zip         *string        `json:"zip,omitzero"`
	OtherFields map[string]any `json:"-"`
}

var prescriptionCustomerAddressKnownKeys = map[string]struct{}{
	"country":    {},
	"street":     {},
	"number":     {},
	"complement": {},
	"city":       {},
	"uf":         {},
	"district":   {},
	"zip":        {},
}

func (a PrescriptionCustomerAddress) MarshalJSON() ([]byte, error) {
	type alias PrescriptionCustomerAddress
	knownBytes, err := json.Marshal(alias(a))
	if err != nil {
		return nil, err
	}

	if len(a.OtherFields) == 0 {
		return knownBytes, nil
	}

	var merged map[string]json.RawMessage
	if err := json.Unmarshal(knownBytes, &merged); err != nil {
		return nil, err
	}

	for k, v := range a.OtherFields {
		if _, known := prescriptionCustomerAddressKnownKeys[k]; known {
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

func (a *PrescriptionCustomerAddress) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	known := map[string]any{
		"country":    &a.Country,
		"street":     &a.Street,
		"number":     &a.Number,
		"complement": &a.Complement,
		"city":       &a.City,
		"uf":         &a.UF,
		"district":   &a.District,
		"zip":        &a.Zip,
	}

	for key, dst := range known {
		if v, ok := raw[key]; ok {
			if err := json.Unmarshal(v, dst); err != nil {
				return err
			}
		}
	}

	for k, v := range raw {
		if _, isKnown := prescriptionCustomerAddressKnownKeys[k]; isKnown {
			continue
		}
		if a.OtherFields == nil {
			a.OtherFields = make(map[string]any)
		}
		var decoded any
		if err := json.Unmarshal(v, &decoded); err != nil {
			return err
		}
		a.OtherFields[k] = decoded
	}

	return nil
}
