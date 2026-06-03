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
	OtherFields       map[string]any        `json:"-" zajson:"-,remain"`
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

	for k, v := range raw {
		var err error
		switch k {
		case "id":
			err = json.Unmarshal(v, &p.ID)
		case "display_id":
			err = json.Unmarshal(v, &p.DisplayID)
		case "doctor_id":
			err = json.Unmarshal(v, &p.DoctorID)
		case "customer":
			err = json.Unmarshal(v, &p.Customer)
		case "file_url":
			err = json.Unmarshal(v, &p.FileURL)
		case "preseller_id":
			err = json.Unmarshal(v, &p.PresellerID)
		case "preseller_name":
			err = json.Unmarshal(v, &p.PresellerName)
		case "observation":
			err = json.Unmarshal(v, &p.Observation)
		case "observation_author":
			err = json.Unmarshal(v, &p.ObservationAuthor)
		default:
			if p.OtherFields == nil {
				p.OtherFields = make(map[string]any)
			}
			var decoded any
			if err := json.Unmarshal(v, &decoded); err != nil {
				return err
			}
			p.OtherFields[k] = decoded
		}
		if err != nil {
			return err
		}
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
	OtherFields      map[string]any               `json:"-" zajson:"-,remain"`
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

	for k, v := range raw {
		var err error
		switch k {
		case "doctor_customer_id":
			err = json.Unmarshal(v, &c.DoctorCustomerID)
		case "account_id":
			err = json.Unmarshal(v, &c.AccountID)
		case "customer_id":
			err = json.Unmarshal(v, &c.CustomerID)
		case "name":
			err = json.Unmarshal(v, &c.Name)
		case "cpf":
			err = json.Unmarshal(v, &c.CPF)
		case "phone":
			err = json.Unmarshal(v, &c.Phone)
		case "email":
			err = json.Unmarshal(v, &c.Email)
		case "birthdate":
			err = json.Unmarshal(v, &c.Birthdate)
		case "address":
			err = json.Unmarshal(v, &c.Address)
		default:
			if c.OtherFields == nil {
				c.OtherFields = make(map[string]any)
			}
			var decoded any
			if err := json.Unmarshal(v, &decoded); err != nil {
				return err
			}
			c.OtherFields[k] = decoded
		}
		if err != nil {
			return err
		}
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
	OtherFields map[string]any `json:"-" zajson:"-,remain"`
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

	for k, v := range raw {
		var err error
		switch k {
		case "country":
			err = json.Unmarshal(v, &a.Country)
		case "street":
			err = json.Unmarshal(v, &a.Street)
		case "number":
			err = json.Unmarshal(v, &a.Number)
		case "complement":
			err = json.Unmarshal(v, &a.Complement)
		case "city":
			err = json.Unmarshal(v, &a.City)
		case "uf":
			err = json.Unmarshal(v, &a.UF)
		case "district":
			err = json.Unmarshal(v, &a.District)
		case "zip":
			err = json.Unmarshal(v, &a.Zip)
		default:
			if a.OtherFields == nil {
				a.OtherFields = make(map[string]any)
			}
			var decoded any
			if err := json.Unmarshal(v, &decoded); err != nil {
				return err
			}
			a.OtherFields[k] = decoded
		}
		if err != nil {
			return err
		}
	}

	return nil
}
