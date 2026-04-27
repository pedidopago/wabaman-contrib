package wsapi

import (
	"encoding/json"
	"time"
)

// IncomingCallFromClient is broadcast by the server when a WhatsApp client initiates a voice call.
type IncomingCallFromClient struct {
	CallID             string `json:"call_id"`
	PhoneID            uint   `json:"phone_id"`
	BranchID           string `json:"branch_id"`
	WABAContactID      string `json:"waba_contact_id"`
	ContactPhoneNumber string `json:"contact_phone_number"`
	ContactID          uint64 `json:"contact_id"`
	ContactName        string `json:"contact_name"`
	IsCallTakeover     bool   `json:"is_call_takeover,omitempty"`
	OriginalAgentID    string `json:"original_agent_id,omitempty"`
}

// ReconnectCall is sent by a browser client to reconnect to an active call after a disconnection.
// AgentID and AgentName are filled by the server, not the client.
type ReconnectCall struct {
	CallID    string `json:"call_id"`
	PhoneID   uint   `json:"phone_id"`
	BranchID  string `json:"branch_id"`
	OfferSDP  string `json:"offer_sdp"`
	AgentID   string `json:"agent_id"`   // filled by the Wabaman server, not the client
	AgentName string `json:"agent_name"` // filled by the Wabaman server, not the client
}

// ActiveCallNotification is sent by the server to a reconnecting agent when it detects
// the agent has an active call in progress.
type ActiveCallNotification struct {
	CallID             string `json:"call_id"`
	PhoneID            uint   `json:"phone_id"`
	BranchID           string `json:"branch_id"`
	ContactID          uint64 `json:"contact_id"`
	ContactName        string `json:"contact_name"`
	ContactPhoneNumber string `json:"contact_phone_number"`
	WABAContactID      string `json:"waba_contact_id"`
}

// SetupCallFromBrowser is sent by a browser client to initiate WebRTC call setup.
// AgentID and AgentName are filled by the server, not the client.
type SetupCallFromBrowser struct {
	CallID    string `json:"call_id"`
	PhoneID   uint   `json:"phone_id"`
	BranchID  string `json:"branch_id"`
	OfferSDP  string `json:"offer_sdp"`
	AgentID   string `json:"agent_id"`   // filled by the Wabaman server, not the client
	AgentName string `json:"agent_name"` // filled by the Wabaman server, not the client
}

// TerminateCall is sent to end an active call. It can be sent by either the browser client
// or the server.
type TerminateCall struct {
	CallID string `json:"call_id"`
}

// CallConsumed is broadcast by the server when an agent picks up an incoming call,
// signaling other clients that the call is no longer available.
type CallConsumed struct {
	PhoneID     uint   `json:"phone_id"`
	BranchID    string `json:"branch_id"`
	CallID      string `json:"call_id"`
	AgentID     string `json:"agent_id"`
	AgentName   string `json:"agent_name"`
	ContactID   uint64 `json:"contact_id"`
	ContactName string `json:"contact_name"`
}

// AcceptCall is sent by a browser client to accept an incoming call.
type AcceptCall struct {
	CallID string `json:"call_id"`
}

// RejectCall is sent by a browser client to reject an incoming call.
type RejectCall struct {
	CallID string `json:"call_id"`
}

// ICECandidate represents a WebRTC ICE candidate used during call negotiation.
// Fields mirror the RTCIceCandidate Web API.
type ICECandidate struct {
	// A string containing the IP address of the candidate.
	Address string `json:"address,omitempty"`
	// A string representing the transport address for the candidate that can be used for connectivity checks.
	// The format of this address is a candidate-attribute as defined in RFC 5245.
	// This string is empty ("") if the RTCIceCandidate is an "end of candidates" indicator.
	Candidate string `json:"candidate"`
	// A string which indicates whether the candidate is an RTP or an RTCP candidate;
	// its value is either rtp or rtcp, and is derived from the "component-id"
	// field in the candidate a-line string.
	Component string `json:"component,omitempty"`
	// Returns a string containing a unique identifier that is the same for any candidates of the same type,
	// share the same base (the address from which the ICE agent sent the candidate), and come from the same
	// STUN server. This is used to help optimize ICE performance while prioritizing and correlating candidates
	// that appear on multiple RTCIceTransport objects.
	Foundation string `json:"foundation,omitempty"`
	// An integer value indicating the candidate's port number.
	Port int `json:"port,omitempty"`
	// A long integer value indicating the candidate's priority.
	Priority int `json:"priority,omitempty"`
	// A string indicating whether the candidate's protocol is "tcp" or "udp".
	Protocol string `json:"protocol,omitempty"`
	// If the candidate is derived from another candidate, relatedAddress is a string containing that
	// host candidate's IP address. For host candidates, this value is null.
	RelatedAddress string `json:"relatedAddress,omitempty"`
	// For a candidate that is derived from another, such as a relay or reflexive candidate, the relatedPort
	// is a number indicating the port number of the candidate from which this candidate is derived. For host
	// candidates, the relatedPort property is null.
	RelatedPort int `json:"relatedPort,omitempty"`
	// A string specifying the candidate's media stream identification tag which uniquely identifies the media
	// stream within the component with which the candidate is associated, or null if no such association exists.
	SdpMid *string `json:"sdpMid"`
	// If not null, sdpMLineIndex indicates the zero-based index number of the media description
	// (as defined in RFC 4566) in the SDP with which the candidate is associated.
	SDPMLineIndex *uint16 `json:"sdpMLineIndex"`
	// If protocol is "tcp", tcpType represents the type of TCP candidate. Otherwise, tcpType is null.
	TcpType string `json:"tcpType,omitempty"`
	// A string indicating the type of candidate as one of the strings listed on RTCIceCandidate.type.
	Type string `json:"type,omitempty"`
	// A string containing a randomly-generated username fragment ("ice-ufrag") which ICE uses for message
	// integrity along with a randomly-generated password ("ice-pwd"). You can use this string to verify
	// generations of ICE generation; each generation of the same ICE process will use the same usernameFragment,
	// even across ICE restarts.
	UsernameFragment *string `json:"usernameFragment"`
}

// SendBrowserCandidate is a bidirectional message carrying an ICE candidate
// between the browser client and the server during WebRTC negotiation.
type SendBrowserCandidate struct {
	PhoneID   uint            `json:"phone_id"`
	CallID    string          `json:"call_id"`
	Candidate json.RawMessage `json:"candidate"`
}

// CallStarted is broadcast by the server when a call has been successfully established.
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

// CallEnded is broadcast by the server when a call has ended.
type CallEnded struct {
	CallID             string    `json:"call_id"`
	PhoneID            uint      `json:"phone_id"`
	BranchID           string    `json:"branch_id"`
	WABAContactID      string    `json:"waba_contact_id"`
	ContactPhoneNumber string    `json:"contact_phone_number,omitempty"`
	ContactID          uint64    `json:"contact_id"`
	ContactName        string    `json:"contact_name"`
	EndTime            time.Time `json:"end_time"`
}

// CallOnAnswerSDP is sent by the server to deliver the SDP answer during WebRTC negotiation.
type CallOnAnswerSDP struct {
	CallID   string `json:"call_id"`
	PhoneID  uint   `json:"phone_id"`
	BranchID string `json:"branch_id"`
	SDP      string `json:"sdp"`
}

// CallStartTimer is broadcast by the server to signal clients to start the call duration timer.
type CallStartTimer struct {
	CallID   string `json:"call_id"`
	PhoneID  uint   `json:"phone_id"`
	BranchID string `json:"branch_id"`
}

// TerminateCallOrigin indicates which side initiated the call termination.
type TerminateCallOrigin int

const (
	TerminateCallOriginWhatsApp TerminateCallOrigin = iota + 1
	TerminateCallOriginBrowser
)
