package msgdriver

import "encoding/json"

type Query struct {
	Key    string         `json:"key"`
	Params map[string]any `json:"params"`
}

type QueryResult struct {
	Data json.RawMessage `json:"data"`
}
