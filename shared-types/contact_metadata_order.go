package types

import "encoding/json"

type ContactMetadataOrder struct {
	Status       OrderStatus        `json:"status,omitzero"`
	Prescription *OrderPrescription `json:"prescription,omitzero"`
	OtherFields  map[string]any     `json:"-"`
}

type OrderPrescription struct {
	IsMedical   MetadataBool   `json:"is_medical,omitzero"`
	Retained    MetadataBool   `json:"retained,omitzero"`
	OtherFields map[string]any `json:"-"`
}

var prescriptionKnownKeys = map[string]struct{}{
	"is_medical": {},
	"retained":   {},
}

func (p OrderPrescription) MarshalJSON() ([]byte, error) {
	type alias OrderPrescription
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
		if _, known := prescriptionKnownKeys[k]; known {
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

func (p *OrderPrescription) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		var err error
		switch k {
		case "is_medical":
			p.IsMedical, err = internMetadataBool(v)
		case "retained":
			p.Retained, err = internMetadataBool(v)
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

var contactMetadataOrderKnownKeys = map[string]struct{}{
	"status":       {},
	"prescription": {},
}

func (o ContactMetadataOrder) MarshalJSON() ([]byte, error) {
	type alias ContactMetadataOrder
	knownBytes, err := json.Marshal(alias(o))
	if err != nil {
		return nil, err
	}

	if len(o.OtherFields) == 0 {
		return knownBytes, nil
	}

	var merged map[string]json.RawMessage
	if err := json.Unmarshal(knownBytes, &merged); err != nil {
		return nil, err
	}

	for k, v := range o.OtherFields {
		if _, known := contactMetadataOrderKnownKeys[k]; known {
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

func (o *ContactMetadataOrder) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		var err error
		switch k {
		case "status":
			o.Status, err = internOrderStatus(v)
		case "prescription":
			err = json.Unmarshal(v, &o.Prescription)
		default:
			if o.OtherFields == nil {
				o.OtherFields = make(map[string]any)
			}
			var decoded any
			if err := json.Unmarshal(v, &decoded); err != nil {
				return err
			}
			o.OtherFields[k] = decoded
		}
		if err != nil {
			return err
		}
	}

	return nil
}
