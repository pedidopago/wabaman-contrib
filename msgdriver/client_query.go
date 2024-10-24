package msgdriver

import "encoding/json"

type ClientQuery struct {
	Key    string         `json:"key"`
	Params map[string]any `json:"params"`
}

type ClientQueryResult struct {
	Data json.RawMessage `json:"data"`
}
