package msgdriver

type ReadMessage struct {
	MessageID        string `json:"message_id,omitempty"`
	WabamanMessageID uint64 `json:"wabaman_message_id,omitempty"`
}
