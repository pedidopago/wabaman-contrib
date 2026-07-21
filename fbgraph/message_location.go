package fbgraph

import (
	"bytes"
	"encoding/json"
)

type LocationObject struct {
	// Endereço da localização. (Opcional)
	Address string `json:"address,omitzero"`
	// Latitude da localização em graus decimais. (Obrigatório)
	Latitude string `json:"latitude"`
	// Longitude da localização em graus decimais. (Obrigatório)
	Longitude string `json:"longitude"`
	// Nome da localização. (Opcional)
	Name string `json:"name,omitzero"`
}

// UnmarshalJSON accepts latitude/longitude as either JSON strings (the send
// API contract) or JSON numbers (what Meta emits on inbound location
// webhooks). Both are normalized to the string form this type stores, so a
// single type can be decoded from both directions.
//
// ponytail: temporary shim for the send/receive type mismatch. Once inbound
// (number) and outbound (string) location payloads use distinct types, drop
// this method and let the struct tags decode directly.
func (l *LocationObject) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Address   string          `json:"address"`
		Latitude  json.RawMessage `json:"latitude"`
		Longitude json.RawMessage `json:"longitude"`
		Name      string          `json:"name"`
	}
	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}
	l.Address = shadow.Address
	l.Latitude = jsonScalarToString(shadow.Latitude)
	l.Longitude = jsonScalarToString(shadow.Longitude)
	l.Name = shadow.Name
	return nil
}

// jsonScalarToString renders a JSON scalar as a plain string: a quoted string
// is unquoted, a bare number is returned verbatim, null/absent becomes "".
func jsonScalarToString(raw json.RawMessage) string {
	s := string(bytes.TrimSpace(raw))
	if s == "" || s == "null" {
		return ""
	}
	if s[0] == '"' {
		var out string
		if err := json.Unmarshal(raw, &out); err == nil {
			return out
		}
	}
	return s
}
