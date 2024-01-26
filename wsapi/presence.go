package wsapi

type PresenceViewClient struct {
	AgentID   string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"` // optional, but it is required if the agent_id is empty
	ClientID  uint64 `json:"client_id,omitempty"`
}

type PresenceTypingToClient struct {
	AgentID   string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"` // optional, but it is required if the agent_id is empty
	ClientID  uint64 `json:"client_id,omitempty"`
}

type PresenceRequest struct {
	ClientID uint64 `json:"client_id,omitempty"`
}

type PresenceAgent struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"` // unix seconds since epoch
}

type PresenceResponse struct {
	ClientID uint64          `json:"client_id,omitempty"`
	Viewing  []PresenceAgent `json:"viewing,omitempty"`
	Typing   []PresenceAgent `json:"typing,omitempty"`
}
