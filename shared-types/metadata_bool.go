package types

import "encoding/json"

type MetadataBool *bool

var (
	boolTrue  = true
	boolFalse = false

	MetadataBoolTrue  MetadataBool = &boolTrue
	MetadataBoolFalse MetadataBool = &boolFalse
)

var metadataBoolIntern = map[bool]MetadataBool{
	true:  MetadataBoolTrue,
	false: MetadataBoolFalse,
}

func internMetadataBool(raw json.RawMessage) (MetadataBool, error) {
	var b bool
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, err
	}
	return metadataBoolIntern[b], nil
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
