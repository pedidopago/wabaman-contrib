package wsapi

// PresenceViewClient is sent by a browser client to indicate that an agent is viewing
// a contact's conversation.
type PresenceViewClient struct {
	AgentID   string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"` // optional, but it is required if the agent_id is empty
	ClientID  uint64 `json:"client_id,omitempty"`
}

// PresenceTypingToClient is sent by a browser client to indicate that an agent is typing
// a reply to a contact.
type PresenceTypingToClient struct {
	AgentID   string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"` // optional, but it is required if the agent_id is empty
	ClientID  uint64 `json:"client_id,omitempty"`
}

// PresenceRequest is sent by a browser client to request the current presence state
// (viewing/typing agents) for a contact.
type PresenceRequest struct {
	ClientID uint64 `json:"client_id,omitempty"`
}

// PresenceAgent identifies an agent in a presence response along with the timestamp
// of their last activity.
type PresenceAgent struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"` // unix seconds since epoch
}

// PresenceResponse is sent by the server in reply to a [PresenceRequest], listing
// the agents currently viewing or typing in a contact's conversation.
type PresenceResponse struct {
	ClientID uint64          `json:"client_id,omitempty"`
	Viewing  []PresenceAgent `json:"viewing,omitempty"`
	Typing   []PresenceAgent `json:"typing,omitempty"`
}
