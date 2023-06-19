package event

import (
	"encoding/json"

	"github.com/pedidopago/wabaman-contrib/fbgraph"
)

type TemplateEventKind string

const (
	TemplateEventApproved        TemplateEventKind = "APPROVED"
	TemplateEventDisabled        TemplateEventKind = "DISABLED"
	TemplateEventInAppeal        TemplateEventKind = "IN_APPEAL"
	TemplateEventPending         TemplateEventKind = "PENDING"
	TemplateEventReinstated      TemplateEventKind = "REINSTATED"
	TemplateEventRejected        TemplateEventKind = "REJECTED"
	TemplateEventFlagged         TemplateEventKind = "FLAGGED"
	TemplateEventPendingDeletion TemplateEventKind = "PENDING_DELETION"
)

type TemplateEvent struct {
	StoreID    string                   `json:"store_id"`
	BranchID   string                   `json:"branch_id,omitempty"`
	PhoneID    uint                     `json:"phone_id,omitempty"`
	Event      TemplateEventKind        `json:"event,omitempty"`
	TemplateID string                   `json:"template_id,omitempty"`
	Template   *fbgraph.MessageTemplate `json:"template,omitempty"`
}

func (e TemplateEvent) ToJSON() string {
	d, _ := json.Marshal(e)
	return string(d)
}

func (e *TemplateEvent) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), e)
}
