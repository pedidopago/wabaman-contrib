package event

import (
	"encoding/json"

	"github.com/pedidopago/wabaman-contrib/fbgraph"
)

type TemplateEventKind string

const (
	TemplateEventApproved        TemplateEventKind = "APPROVED"
	TemplateEventArchived        TemplateEventKind = "ARCHIVED"
	TemplateEventDeleted         TemplateEventKind = "DELETED"
	TemplateEventDisabled        TemplateEventKind = "DISABLED"
	TemplateEventFlagged         TemplateEventKind = "FLAGGED"
	TemplateEventInAppeal        TemplateEventKind = "IN_APPEAL"
	TemplateEventLimitExceeded   TemplateEventKind = "LIMIT_EXCEEDED"
	TemplateEventLocked          TemplateEventKind = "LOCKED"
	TemplateEventPaused          TemplateEventKind = "PAUSED"
	TemplateEventPending         TemplateEventKind = "PENDING"
	TemplateEventReinstated      TemplateEventKind = "REINSTATED"
	TemplateEventPendingDeletion TemplateEventKind = "PENDING_DELETION"
	TemplateEventRejected        TemplateEventKind = "REJECTED"
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
