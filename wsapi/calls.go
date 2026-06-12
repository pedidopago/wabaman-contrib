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

// CallAnswerIntent is sent by a browser client to the server to indicate the agent is ready to answer the call.
type CallAnswerIntent struct {
	CallID    string `json:"call_id"`
	AgentID   string `json:"agent_id,omitempty"`   // filled by the Wabaman server, not the client
	AgentName string `json:"agent_name,omitempty"` // filled by the Wabaman server, not the client
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
	AgentID   string          `json:"agent_id,omitempty"` // set for multi-agent calls to route to the correct agent's peer connection
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
// StartedAt, when non-zero, tells the client the call was answered at this unix
// timestamp (useful for reconnect/takeover so the timer reflects elapsed time).
type CallStartTimer struct {
	CallID    string `json:"call_id"`
	PhoneID   uint   `json:"phone_id"`
	BranchID  string `json:"branch_id"`
	StartedAt int64  `json:"started_at,omitempty"`
}

// TerminateCallOrigin indicates which side initiated the call termination.
type TerminateCallOrigin int

const (
	TerminateCallOriginWhatsApp TerminateCallOrigin = iota + 1
	TerminateCallOriginBrowser
)

// RequestCallEligibility is sent by a browser client to ask whether a contact can be called
// right now. The server replies with [CallEligibility] for the same (phone_id, contact_id) pair,
// and additionally pushes unsolicited updates via [CallPermissionState] when the underlying
// permission state changes.
type RequestCallEligibility struct {
	PhoneID   uint   `json:"phone_id"`
	ContactID uint64 `json:"contact_id"`
}

// CallPermissionLimit mirrors a single limit entry from Meta's GET /call_permissions endpoint
// (e.g. "start_call" allowed 100 per PT24H), surfaced verbatim so the UI can render usage.
type CallPermissionLimit struct {
	// ISO 8601 duration, e.g. PT24H or P7D.
	TimePeriod   string `json:"time_period"`
	MaxAllowed   int    `json:"max_allowed"`
	CurrentUsage int    `json:"current_usage"`
	// Unix timestamp; present only when the limit window is currently exhausted.
	LimitExpirationTime int64 `json:"limit_expiration_time,omitempty"`
}

// CallEligibility is the server's reply to a [RequestCallEligibility] and also represents
// the can-call-now state at a point in time.
type CallEligibility struct {
	PhoneID   uint   `json:"phone_id"`
	ContactID uint64 `json:"contact_id"`
	CanCall   bool   `json:"can_call"`
	// no_permission | temporary | permanent
	Status string `json:"status"`
	// Unix timestamp; 0 for permanent permissions or when no permission is granted.
	ExpiresAt             int64 `json:"expires_at,omitempty"`
	ConsecutiveUnanswered uint8 `json:"consecutive_unanswered"`
	// One of: NO_PERMISSION | EXPIRED | RATE_LIMIT | INTERNAL. Empty when CanCall is true.
	BlockedReason           string                `json:"blocked_reason,omitempty"`
	StartCallLimits         []CallPermissionLimit `json:"start_call_limits,omitempty"`
	PermissionRequestLimits []CallPermissionLimit `json:"permission_request_limits,omitempty"`
}

// CallPermissionState is the unsolicited server-push variant of [CallEligibility],
// emitted when a call_permission_reply webhook updates the underlying state or after
// a connected call resets the consecutive-unanswered counter.
type CallPermissionState struct {
	CallEligibility
}

// InitiateCallFromBrowser is sent by a browser client to start an outbound (business-initiated)
// call to a contact. OfferSDP is the browser's local-side WebRTC offer toward the gateway
// (ms-wabaman-webrtc); the gateway mints a separate offer for Cloud API.
type InitiateCallFromBrowser struct {
	PhoneID   uint   `json:"phone_id"`
	ContactID uint64 `json:"contact_id"`
	OfferSDP  string `json:"offer_sdp"`
	// Filled by the server from the session, not the client.
	AgentID   string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
}

// CallInitiated confirms an outbound call has been accepted by Cloud API and a call_id assigned.
type CallInitiated struct {
	CallID    string `json:"call_id"`
	PhoneID   uint   `json:"phone_id"`
	ContactID uint64 `json:"contact_id"`
}

// CallInitiateFailed reports an outbound initiate failure. Reasons:
//
//	NO_PERMISSION   — user has not granted permission (or it has expired)
//	RATE_LIMIT      — start_call action limit reached
//	EXPIRED         — temporary permission has expired
//	SDP_INVALID     — Cloud API rejected the SDP offer
//	INTERNAL        — anything else (network, etc.)
type CallInitiateFailed struct {
	PhoneID   uint   `json:"phone_id"`
	ContactID uint64 `json:"contact_id"`
	Reason    string `json:"reason"`
	Message   string `json:"message,omitempty"`
	// Unix timestamp; present for RATE_LIMIT when the window expiration is known.
	RetryAfter int64 `json:"retry_after,omitempty"`
}

// CallStatusUpdated is broadcast when an outbound call moves between lifecycle states.
// Values: initiating | connecting | ringing | in_progress | completed | missed | rejected.
type CallStatusUpdated struct {
	CallID    string `json:"call_id"`
	PhoneID   uint   `json:"phone_id"`
	Status    string `json:"status"`
	ContactID uint64 `json:"contact_id,omitempty"`
	AgentID   string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
}

// SendCallPermissionRequest is the agent's intent to send a call-permission request to a contact.
// The server picks free-form vs template based on whether a 24h customer-service window is open.
type SendCallPermissionRequest struct {
	PhoneID   uint   `json:"phone_id"`
	ContactID uint64 `json:"contact_id"`
	BodyText  string `json:"body_text,omitempty"`
	// Filled by the server from the session, not the client.
	AgentID   string `json:"agent_id,omitempty"`
	AgentName string `json:"agent_name,omitempty"`
}

// SendCallPermissionResponse is sent by the server after handling a [SendCallPermissionRequest].
// On failure, Success is false and Error carries a human-readable reason.
// On success, WAMID is the WhatsApp message ID and Branch is "interactive" or "template".
type SendCallPermissionResponse struct {
	PhoneID   uint   `json:"phone_id"`
	ContactID uint64 `json:"contact_id"`
	Success   bool   `json:"success"`
	// Set when Success is false.
	Error string `json:"error,omitempty"`
	// Set when Success is true: WhatsApp message ID returned by Meta.
	WAMID string `json:"wamid,omitempty"`
	// Set when Success is true: "interactive" (24h window) or "template" (fallback).
	Branch string `json:"branch,omitempty"`
}

// CallInviteAgent is sent by an in-call agent to invite another agent to join.
// Server fills InviterAgentID and InviterAgentName from JWT.
type CallInviteAgent struct {
	CallID              string `json:"call_id"`
	PhoneID             uint   `json:"phone_id"`
	BranchID            string `json:"branch_id"`
	ContactID           uint64 `json:"contact_id"`
	ContactName         string `json:"contact_name"`
	ContactPhoneNumber  string `json:"contact_phone_number"`
	TargetAgentID       string `json:"target_agent_id"`
	InviterAgentID      string `json:"inviter_agent_id"`
	InviterAgentName    string `json:"inviter_agent_name"`
}

// CallInviteAck is sent by the target agent to confirm receipt of the invite (delivery confirmation).
// AgentID is filled by the server from JWT.
type CallInviteAck struct {
	CallID         string `json:"call_id"`
	AgentID        string `json:"agent_id"`         // server fills from JWT
	InviterAgentID string `json:"inviter_agent_id"` // client sets from the invite payload
}

// CallInviteAccepted is sent by the target agent to accept the invite; the agent will start WebRTC join.
// AgentID and AgentName are filled by the server from JWT.
type CallInviteAccepted struct {
	CallID         string `json:"call_id"`
	AgentID        string `json:"agent_id"`         // server fills from JWT
	AgentName      string `json:"agent_name"`       // server fills from JWT
	InviterAgentID string `json:"inviter_agent_id"` // client sets from the invite payload
}

// CallInviteRejected is sent by the target agent to decline the invite.
// AgentID and AgentName are filled by the server from JWT.
type CallInviteRejected struct {
	CallID         string `json:"call_id"`
	AgentID        string `json:"agent_id"`         // server fills from JWT
	AgentName      string `json:"agent_name"`       // server fills from JWT
	InviterAgentID string `json:"inviter_agent_id"` // client sets from the invite payload
}

// LeaveCall is sent by an agent to leave a multi-agent call without terminating it.
type LeaveCall struct {
	CallID string `json:"call_id"`
}

// JoinCall is sent by a secondary agent to join an active multi-agent call.
// AgentID and AgentName are filled by the server from JWT.
type JoinCall struct {
	CallID    string `json:"call_id"`
	PhoneID   uint   `json:"phone_id"`
	OfferSDP  string `json:"offer_sdp"`
	AgentID   string `json:"agent_id"`   // server fills from JWT
	AgentName string `json:"agent_name"` // server fills from JWT
}

// JoinCallAnswer is sent by the server back to the joining agent with the SDP answer.
// On failure AnswerSDP is empty and Reason is set to one of:
// CALL_NOT_FOUND, CALL_FULL, ALREADY_IN_CALL, JOIN_FAILED.
type JoinCallAnswer struct {
	CallID    string `json:"call_id"`
	AnswerSDP string `json:"answer_sdp"`
	Reason    string `json:"reason,omitempty"`
}

// CallAgentJoined is broadcast to all agents in a call when a new agent's WebRTC is live.
type CallAgentJoined struct {
	CallID           string `json:"call_id"`
	PhoneID          uint   `json:"phone_id"`
	BranchID         string `json:"branch_id"`
	ContactID        uint64 `json:"contact_id"`
	AgentID          string `json:"agent_id"`
	AgentName        string `json:"agent_name"`
	ParticipantCount int    `json:"participant_count"`
}

// CallAgentLeft is broadcast to all agents in a call when an agent disconnects.
type CallAgentLeft struct {
	CallID           string `json:"call_id"`
	PhoneID          uint   `json:"phone_id"`
	BranchID         string `json:"branch_id"`
	ContactID        uint64 `json:"contact_id"`
	AgentID          string `json:"agent_id"`
	AgentName        string `json:"agent_name"`
	ParticipantCount int    `json:"participant_count"`
}

// CallInviteFailed is sent by the server to the inviter to explain why the invite was rejected.
// Reason values: CALL_NOT_FOUND, NOT_IN_CALL, ALREADY_IN_CALL, CALL_FULL, AGENT_NOT_FOUND.
type CallInviteFailed struct {
	CallID          string `json:"call_id"`
	TargetAgentID   string `json:"target_agent_id"`
	TargetAgentName string `json:"target_agent_name"`
	Reason          string `json:"reason"`
}
