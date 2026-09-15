package event

import (
	"encoding/json"
	"time"
)

// Reasons a contact was replaced. Both are uniqueness collisions -- never a
// deliberate operator action. A contact is unique on a phone by two keys, its
// WhatsApp id and its BSUID, and a collision on either one replaces it.
const (
	// ContactReplacedIdentityCollision: the contact's WhatsApp identity changed
	// (phone number rewritten by Meta, or a BSUID adopted) to a value another
	// contact on the same phone already held. Arrives one at a time, driven by
	// an inbound webhook or a send.
	ContactReplacedIdentityCollision = "identity_collision"
	// ContactReplacedPhoneDataMigration: a branch moved its phone number, and
	// the contact carried across met a contact already sitting on the
	// destination phone under the same WhatsApp id. Arrives in a burst, one per
	// casualty, at the end of the migration.
	//
	// Carries the same guarantee as an identity collision: the replaced contact
	// is merged into the survivor before it is destroyed, so its messages, held
	// messages, conversation history and tags move with it.
	ContactReplacedPhoneDataMigration = "phone_data_migration"
)

// ContactReplacedEvent announces that OldContactID no longer exists and that
// NewContactID inherited its history, its messages and its FK references.
//
// Both reasons carry that same guarantee: each goes through the same merge
// before the replaced contact is deleted. Reason distinguishes what triggered
// the replacement, not how much survived it.
//
// Consumers must re-point anything keyed by OldContactID. The event is emitted
// after ms-wabaman has already committed the merge, so it describes a completed
// fact -- the old id is gone by the time this is readable, and a consumer that
// tries to load it will find nothing.
//
// There is always a successor: this event is only published when a contact lost
// a uniqueness collision against another contact, which is precisely a contact
// that survives it. It is NOT published when a contact is simply deleted (the
// ghost-contact sweep, or the delete endpoint) -- that is a different fact, with
// the opposite instruction to a consumer: drop the data rather than re-point it.
//
// Safe to reprocess: applying it is
// `UPDATE ... SET contact_id = new WHERE contact_id = old`, which is a no-op the
// second time.
//
// Replacements can chain -- a contact replaced today can itself be replaced
// tomorrow, so an old id can appear as OldContactID in one event and as
// NewContactID in an earlier one. Events are ordered on the queue, so applying
// them in receipt order converges on the final survivor. A consumer that
// processes this queue concurrently gives that up and has to re-resolve the
// chain itself.
type ContactReplacedEvent struct {
	// OldContactID is the contact that was absorbed and deleted. Never zero.
	OldContactID uint64 `json:"old_contact_id"`
	// NewContactID is the contact that survived and inherited everything.
	// Never zero -- a replacement without a successor is not published.
	NewContactID uint64 `json:"new_contact_id"`

	// CustomerID as held by the replaced contact, for consumers keyed by
	// customer rather than by contact. Empty when the contact had none.
	CustomerID string `json:"customer_id,omitempty"`
	// BranchID is taken from the SURVIVOR's phone, not the replaced contact's.
	// In a phone data migration the replaced contact sits on a throwaway
	// temporary phone whose branch id is a generated placeholder that never
	// named a real branch.
	BranchID string `json:"branch_id,omitempty"`

	// Reason is one of the ContactReplaced* constants. Consumers re-point
	// identically either way; it distinguishes a lone event from a migration
	// burst.
	Reason string `json:"reason"`

	// OccurredAt is when the merge committed.
	OccurredAt time.Time `json:"occurred_at"`
}

func (e ContactReplacedEvent) ToJSON() string {
	d, _ := json.Marshal(e)
	return string(d)
}

func (e *ContactReplacedEvent) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), e)
}
