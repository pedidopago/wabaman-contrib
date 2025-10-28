package wsapi

type IncomingCallFromClient struct {
	CallID        string `json:"call_id"`
	PhoneID       uint   `json:"phone_id"`
	BranchID      string `json:"branch_id"`
	WABAContactID string `json:"waba_contact_id"`
	ContactID     uint64 `json:"contact_id"`
	ContactName   string `json:"contact_name"`
	SDPType       string `json:"sdp_type"`
	SDP           string `json:"sdp"`
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
