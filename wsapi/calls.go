package wsapi

import "time"

type IncomingCallFromClient struct {
	CallID             string `json:"call_id"`
	PhoneID            uint   `json:"phone_id"`
	BranchID           string `json:"branch_id"`
	WABAContactID      string `json:"waba_contact_id"`
	ContactPhoneNumber string `json:"contact_phone_number"`
	ContactID          uint64 `json:"contact_id"`
	ContactName        string `json:"contact_name"`
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

type SendBrowserCandidate struct {
	CallID    string `json:"call_id"`
	Candidate struct {
		// A string containing the IP address of the candidate.
		Address string `json:"address"`
		// A string representing the transport address for the candidate that can be used for connectivity checks.
		// The format of this address is a candidate-attribute as defined in RFC 5245.
		// This string is empty ("") if the RTCIceCandidate is an "end of candidates" indicator.
		Candidate string `json:"candidate"`
		// A string which indicates whether the candidate is an RTP or an RTCP candidate;
		// its value is either rtp or rtcp, and is derived from the "component-id"
		// field in the candidate a-line string.
		Component string `json:"component"`
		// Returns a string containing a unique identifier that is the same for any candidates of the same type,
		// share the same base (the address from which the ICE agent sent the candidate), and come from the same
		// STUN server. This is used to help optimize ICE performance while prioritizing and correlating candidates
		// that appear on multiple RTCIceTransport objects.
		Foundation string `json:"foundation"`
		// An integer value indicating the candidate's port number.
		Port int `json:"port"`
		// A long integer value indicating the candidate's priority.
		Priority int `json:"priority"`
		// A string indicating whether the candidate's protocol is "tcp" or "udp".
		Protocol string `json:"protocol"`
		// If the candidate is derived from another candidate, relatedAddress is a string containing that
		// host candidate's IP address. For host candidates, this value is null.
		RelatedAddress string `json:"relatedAddress"`
		// For a candidate that is derived from another, such as a relay or reflexive candidate, the relatedPort
		// is a number indicating the port number of the candidate from which this candidate is derived. For host
		// candidates, the relatedPort property is null.
		RelatedPort int `json:"relatedPort"`
		// A string specifying the candidate's media stream identification tag which uniquely identifies the media
		// stream within the component with which the candidate is associated, or null if no such association exists.
		SdpMid string `json:"sdpMid"`
		// If not null, sdpMLineIndex indicates the zero-based index number of the media description
		// (as defined in RFC 4566) in the SDP with which the candidate is associated.
		SdpMLineIndex *uint16 `json:"sdpMLineIndex"`
		// If protocol is "tcp", tcpType represents the type of TCP candidate. Otherwise, tcpType is null.
		TcpType string `json:"tcpType"`
		// A string indicating the type of candidate as one of the strings listed on RTCIceCandidate.type.
		Type string `json:"type"`
		// A string containing a randomly-generated username fragment ("ice-ufrag") which ICE uses for message
		// integrity along with a randomly-generated password ("ice-pwd"). You can use this string to verify
		// generations of ICE generation; each generation of the same ICE process will use the same usernameFragment,
		// even across ICE restarts.
		UsernameFragment string `json:"usernameFragment"`
	} `json:"candidate"`
}

type CallStarted struct {
	CallID             string    `json:"call_id"`
	PhoneID            uint      `json:"phone_id"`
	BranchID           string    `json:"branch_id"`
	WABAContactID      string    `json:"waba_contact_id"`
	ContactPhoneNumber string    `json:"contact_phone_number"`
	ContactID          uint64    `json:"contact_id"`
	ContactName        string    `json:"contact_name"`
	StartTime          time.Time `json:"start_time"`
}

type CallEnded struct {
	CallID             string    `json:"call_id"`
	PhoneID            uint      `json:"phone_id"`
	BranchID           string    `json:"branch_id"`
	WABAContactID      string    `json:"waba_contact_id"`
	ContactPhoneNumber string    `json:"contact_phone_number"`
	ContactID          uint64    `json:"contact_id"`
	ContactName        string    `json:"contact_name"`
	EndTime            time.Time `json:"end_time"`
}
