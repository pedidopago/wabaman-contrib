package types

import (
	"encoding/json"
)

type SafeMetadata json.RawMessage

func (m SafeMetadata) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}

	return []byte(m), nil
}

func (m *SafeMetadata) UnmarshalJSON(data []byte) error {
	if m == nil || data == nil {
		return nil
	}

	if string(data) == "null" {
		return nil
	}

	*m = make([]byte, len(data))
	copy((*m)[:], data)

	return nil
}
