package types

import (
	"encoding/json"
	"log/slog"
)

type SafeMetadata json.RawMessage

func (m SafeMetadata) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(json.RawMessage(m))
	if err != nil {
		slog.Error("SafeMetadata.MarshalJSON: failed to marshal json.RawMessage", slog.Any("err", err), slog.Any("data", string(m)))
		return nil, nil
	}
	return d, nil
}

func (m *SafeMetadata) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*json.RawMessage)(m))
}
