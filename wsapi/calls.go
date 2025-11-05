package wsapi

type IncomingCallFromClient struct {
	CallID        string `json:"call_id"`
	PhoneID       uint   `json:"phone_id"`
	BranchID      string `json:"branch_id"`
	WABAContactID string `json:"waba_contact_id"`
	ContactID     uint64 `json:"contact_id"`
	ContactName   string `json:"contact_name"`
}

type ConnectToCall struct {
	CallID   string `json:"call_id"`
	PhoneID  uint   `json:"phone_id"`
	BranchID string `json:"branch_id"`
	SDPType  string `json:"sdp_type"`
	SDP      string `json:"sdp"`
}

type TerminateCall struct {
	CallID string `json:"call_id"`
}

type CallConsumed struct {
	CallID    string `json:"call_id"`
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
	SDPType   string `json:"sdp_type"` // must be answer
	SDP       string `json:"sdp"`
}

type AcceptCall struct {
	CallID string `json:"call_id"`
}

type RejectCall struct {
	CallID string `json:"call_id"`
}
