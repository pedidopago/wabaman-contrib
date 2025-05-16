package types

import (
	"encoding/json"
	"errors"
	"log/slog"
	"runtime"
)

// CachedMetadata is not thread safe
type CachedMetadata struct {
	raw        json.RawMessage
	rawIsDirty bool
	parsed     map[string]any
}

func NewCachedMetadataFromJSONBytes(data []byte) *CachedMetadata {
	var m CachedMetadata
	m.raw = make([]byte, len(data))
	copy(m.raw[:], data)

	return &m
}

func NewCachedMetadataPtrFromJSONBytes(data []byte) *CachedMetadata {
	var m CachedMetadata
	m.raw = make([]byte, len(data))
	copy(m.raw[:], data)

	return &m
}

func (m *CachedMetadata) IsEmpty() bool {
	if m == nil {
		return true
	}

	return len(m.raw) == 0 && m.parsed == nil
}

func (m *CachedMetadata) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("{}"), nil
	}

	if len(m.raw) == 0 && m.parsed == nil {
		return []byte("{}"), nil
	}

	if m.rawIsDirty && m.parsed != nil {
		return json.Marshal(m.parsed)
	}

	if len(m.raw) == 0 && m.parsed != nil {
		return json.Marshal(m.parsed)
	}

	return []byte(m.raw), nil
}

func (m *CachedMetadata) UnmarshalJSON(data []byte) error {
	if m == nil || data == nil {
		return nil
	}

	if string(data) == "null" {
		return nil
	}

	m.raw = make([]byte, len(data))
	copy(m.raw[:], data)

	return nil
}

func (m *CachedMetadata) Get(key string) any {
	if m == nil {
		return nil
	}

	if m.parsed != nil {
		return m.parsed[key]
	}

	if len(m.raw) == 0 {
		return nil
	}

	parsed := make(map[string]any)
	err := json.Unmarshal(m.raw, &parsed)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		rawstr := ""
		if len(m.raw) > 0 {
			rawstr = string(m.raw)
		}
		slog.Error("failed to unmarshal metadata (CachedMetadata)", slog.String("error", err.Error()), slog.String("file", file), slog.Int("line", line), slog.String("raw", rawstr))
		return nil
	}
	m.parsed = parsed

	return parsed
}

func (m *CachedMetadata) Set(key string, value any) {
	if m == nil {
		return
	}

	if m.parsed == nil && len(m.raw) > 0 {
		m.parsed = make(map[string]any)
		err := json.Unmarshal(m.raw, &m.parsed)
		if err != nil {
			_, file, line, _ := runtime.Caller(1)
			slog.Error("failed to unmarshal metadata (CachedMetadata)", slog.String("error", err.Error()), slog.String("file", file), slog.Int("line", line))
			return
		}
	}

	m.rawIsDirty = true

	m.parsed[key] = value
}

func (m *CachedMetadata) UnmarshalTo(target any) error {
	if m == nil {
		return nil
	}

	if target == nil {
		return errors.New("target is nil")
	}

	if len(m.raw) != 0 && !m.rawIsDirty {
		return json.Unmarshal(m.raw, target)
	}

	if len(m.parsed) > 0 {
		ms, err := json.Marshal(m.parsed)
		if err != nil {
			return err
		}
		return json.Unmarshal(ms, target)
	}

	return nil
}
