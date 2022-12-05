package fbgraph

import (
	"encoding/json"
	"strings"
)

type Slice[T any] []T

// json marshal/unmarshal

func (s *Slice[T]) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}
	if strings.HasPrefix(string(data), "[") {
		return json.Unmarshal(data, (*[]T)(s))
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*s = []T{v}
	return nil
}

func (s Slice[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal([]T(s))
}
