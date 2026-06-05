package types

import "encoding/json"

type ContactMetadataMarketingDisabledReason struct {
	Error       *int           `json:"error,omitzero"`
	MetaReason  *string        `json:"meta_reason,omitzero"`
	Date        *string        `json:"date,omitzero"`
	OtherFields map[string]any `json:"-" zajson:"-,remain"`
}

var contactMetadataMarketingDisabledReasonKnownKeys = map[string]struct{}{
	"error":       {},
	"meta_reason": {},
	"date":        {},
}

func (m ContactMetadataMarketingDisabledReason) MarshalJSON() ([]byte, error) {
	type alias ContactMetadataMarketingDisabledReason
	knownBytes, err := json.Marshal(alias(m))
	if err != nil {
		return nil, err
	}

	if len(m.OtherFields) == 0 {
		return knownBytes, nil
	}

	var merged map[string]json.RawMessage
	if err := json.Unmarshal(knownBytes, &merged); err != nil {
		return nil, err
	}

	for k, v := range m.OtherFields {
		if _, known := contactMetadataMarketingDisabledReasonKnownKeys[k]; known {
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

func (m *ContactMetadataMarketingDisabledReason) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		var err error
		switch k {
		case "error":
			err = json.Unmarshal(v, &m.Error)
		case "meta_reason":
			err = json.Unmarshal(v, &m.MetaReason)
		case "date":
			err = json.Unmarshal(v, &m.Date)
		default:
			if m.OtherFields == nil {
				m.OtherFields = make(map[string]any)
			}
			var decoded any
			if err := json.Unmarshal(v, &decoded); err != nil {
				return err
			}
			m.OtherFields[k] = decoded
		}
		if err != nil {
			return err
		}
	}

	return nil
}
