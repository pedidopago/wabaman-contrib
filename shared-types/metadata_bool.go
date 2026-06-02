package types

import (
	"encoding/json"
	"fmt"
)

type MetadataBool *bool

var (
	boolTrue  = true
	boolFalse = false

	MetadataBoolTrue  MetadataBool = &boolTrue
	MetadataBoolFalse MetadataBool = &boolFalse
)

func internMetadataBool(raw json.RawMessage) (MetadataBool, error) {
	if len(raw) == 4 && raw[0] == 't' {
		return MetadataBoolTrue, nil
	}
	if len(raw) == 5 && raw[0] == 'f' {
		return MetadataBoolFalse, nil
	}
	return nil, fmt.Errorf("invalid boolean: %s", raw)
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
