package types

import (
	"encoding/json"
	"fmt"

	"github.com/pedidopago/zajson"
)

type MetadataBool *bool

var (
	boolTrue  = true
	boolFalse = false

	MetadataBoolTrue  MetadataBool = &boolTrue
	MetadataBoolFalse MetadataBool = &boolFalse
	MetadataBoolDel   MetadataBool = nil
)

func internMetadataBool(raw json.RawMessage) (MetadataBool, error) {
	if len(raw) == 4 && raw[0] == 'n' {
		return nil, nil
	}
	if len(raw) == 4 && raw[0] == 't' {
		return MetadataBoolTrue, nil
	}
	if len(raw) == 5 && raw[0] == 'f' {
		return MetadataBoolFalse, nil
	}
	if len(raw) >= 2 && raw[0] == '"' {
		if len(raw) == 6 && raw[1] == '$' {
			return nil, nil
		}
		if len(raw) == 6 && raw[1] == 't' {
			return MetadataBoolTrue, nil
		}
		if len(raw) == 7 && raw[1] == 'f' {
			return MetadataBoolFalse, nil
		}
	}
	return nil, fmt.Errorf("invalid boolean: %s", raw)
}

func ZReadMetadataBool(r *zajson.Reader) (MetadataBool, error) {
	switch r.PeekKind() {
	case 'n':
		if err := r.ReadNull(); err != nil {
			return nil, err
		}
		return nil, nil
	case 't':
		if _, err := r.ReadBool(); err != nil {
			return nil, err
		}
		return MetadataBoolTrue, nil
	case 'f':
		if _, err := r.ReadBool(); err != nil {
			return nil, err
		}
		return MetadataBoolFalse, nil
	case '"':
		s, err := r.ReadString()
		if err != nil {
			return nil, err
		}
		switch s {
		case "true":
			return MetadataBoolTrue, nil
		case "false":
			return MetadataBoolFalse, nil
		case "$del":
			return nil, nil
		}
		return nil, fmt.Errorf("invalid boolean string: %q", s)
	default:
		return nil, fmt.Errorf("invalid boolean value, got %c", r.PeekKind())
	}
}

func EqualMetadataBool(a, b MetadataBool) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
