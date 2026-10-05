package event

import (
	"encoding/json"
	"time"
)

// ContactNumberChangedEvent announces that a customer moved their WhatsApp
// account to a new phone number, and ms-wabaman re-keyed the contact to it.
//
// The contact keeps its id: only its number changed. ms-wabaman keeps the
// previous number as an alias of the contact, so its history stays reachable,
// but anything a consumer stored by OldNumber for this contact should now be
// keyed by NewNumber. The old number may later be reassigned by the carrier to
// someone else, so it must not keep being used to reach this customer.
//
// Emitted after ms-wabaman committed the change. When another contact already
// held NewNumber on the same phone, that contact was first merged into this one
// and a ContactReplacedEvent was published for it before this event.
//
// Safe to reprocess: applying it is
// `UPDATE ... SET number = new WHERE number = old`, a no-op the second time.
type ContactNumberChangedEvent struct {
	// ContactID is the contact whose number changed. Never zero.
	ContactID uint64 `json:"contact_id"`
	// PhoneID is the business phone the contact belongs to.
	PhoneID uint64 `json:"phone_id"`
	// BranchID of the business phone, for consumers keyed by branch.
	BranchID string `json:"branch_id,omitempty"`
	// CustomerID linked to the contact, for consumers keyed by customer rather
	// than by contact. Empty when the contact has none.
	CustomerID string `json:"customer_id,omitempty"`

	// OldNumber and NewNumber are WhatsApp IDs as ms-wabaman stores them in the
	// contact's waba_contact_id (digits only, normalized).
	OldNumber string `json:"old_number"`
	NewNumber string `json:"new_number"`

	// UserID is the contact's business-scoped user ID (BSUID), when known.
	UserID string `json:"user_id,omitempty"`

	// OccurredAt is when Meta reported the change.
	OccurredAt time.Time `json:"occurred_at"`
}

func (e ContactNumberChangedEvent) ToJSON() string {
	d, _ := json.Marshal(e)
	return string(d)
}

func (e *ContactNumberChangedEvent) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), e)
}
