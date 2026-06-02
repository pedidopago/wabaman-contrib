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

	for _, bf := range []struct {
		key string
		dst *MetadataBool
	}{
		{"is_medical", &p.IsMedical},
		{"retained", &p.Retained},
	} {
		if v, ok := raw[bf.key]; ok {
			val, err := internMetadataBool(v)
			if err != nil {
				return err
			}
			*bf.dst = val
		}
	}

	for k, v := range raw {
		if _, isKnown := prescriptionKnownKeys[k]; isKnown {
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

	known := map[string]any{
		"prescription": &o.Prescription,
	}

	if v, ok := raw["status"]; ok {
		p, err := internOrderStatus(v)
		if err != nil {
			return err
		}
		o.Status = p
	}

	for key, dst := range known {
		if v, ok := raw[key]; ok {
			if err := json.Unmarshal(v, dst); err != nil {
				return err
			}
		}
	}

	for k, v := range raw {
		if _, isKnown := contactMetadataOrderKnownKeys[k]; isKnown {
			continue
		}
		if o.OtherFields == nil {
			o.OtherFields = make(map[string]any)
		}
		var decoded any
		if err := json.Unmarshal(v, &decoded); err != nil {
			return err
		}
		o.OtherFields[k] = decoded
	}

	return nil
}
