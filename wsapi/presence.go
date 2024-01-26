package wsapi

type PresenceViewClient struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"` // optional, but it is required if the agent_id is empty
}

type PresenceTypingToClient struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"` // optional, but it is required if the agent_id is empty
}

type PresenceAgentViewingClient struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
}

type PresenceAgentTypingToClient struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
}
