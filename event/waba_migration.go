package event

import "encoding/json"

// TemplateMigration maps one template across a WABA currency migration.
//
// Meta clones templates into the destination WABA with NEW graph ids, so every
// stored reference to a template by graph id is stale the moment a migration
// completes. NewGraphID is empty when the template did not make it across --
// Meta silently skips templates in outdated formats, and reports that only by
// their absence.
type TemplateMigration struct {
	Name       string `json:"name"`
	Language   string `json:"language"`
	OldGraphID string `json:"old_graph_id"`
	NewGraphID string `json:"new_graph_id,omitempty"`
}

// Migrated reports whether the template exists on the destination WABA.
func (t TemplateMigration) Migrated() bool {
	return t.NewGraphID != ""
}

// WABAMigratedEvent announces that a WABA was cloned into a new one with a
// different billing currency, and that its phone numbers now live on the clone.
//
// One event per migration, never one per template: a consumer has to apply the
// whole remap or none of it, and a per-template stream would leave the consumer
// half-migrated if it died midway. The event is emitted after ms-wabaman has
// already reconciled its own rows, so it describes a completed fact.
//
// Consumers must re-point anything keyed by the old WABA id or by a template
// graph id. Templates listed in Unmigrated do not exist on the destination WABA
// and cannot be sent -- they should be made unselectable rather than left to
// fail at send time.
type WABAMigratedEvent struct {
	// MigrationID is ms-wabaman's record id, for correlation and idempotency.
	MigrationID uint64 `json:"migration_id"`
	// MetaMigrationID is Meta's own handle ("mig_...").
	MetaMigrationID string `json:"meta_migration_id"`

	SourceWABAID string `json:"source_waba_id"`
	DestWABAID   string `json:"dest_waba_id"`
	Currency     string `json:"currency"`

	// BranchIDs are every branch whose phone moved to the destination WABA.
	// Consumers keying branches by WABA id must re-point these.
	BranchIDs []string `json:"branch_ids,omitempty"`
	// StoreIDs are the stores owning those branches, for consumers scoped by
	// store rather than branch.
	StoreIDs []string `json:"store_ids,omitempty"`

	// Templates that exist on the destination WABA under a new graph id.
	Templates []TemplateMigration `json:"templates,omitempty"`
	// Unmigrated templates: present on the source WABA, absent on the
	// destination, and therefore unsendable until re-authored.
	Unmigrated []TemplateMigration `json:"unmigrated,omitempty"`
}

func (e WABAMigratedEvent) ToJSON() string {
	d, _ := json.Marshal(e)
	return string(d)
}

func (e *WABAMigratedEvent) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), e)
}
