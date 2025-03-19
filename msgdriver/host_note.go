package msgdriver

type HostNote struct {
	ContactID     uint64                 `json:"contact_id"`
	WabaContactID string                 `json:"waba_contact_id"`
	PhoneID       uint                   `json:"phone_id"`
	Format        string                 `json:"format"`
	Title         string                 `json:"title"`
	TitleIcon     string                 `json:"title_icon,omitempty"`
	Description   string                 `json:"description,omitempty"`
	AgentID       string                 `json:"agent_id,omitempty"`
	AgentName     string                 `json:"agent_name,omitempty"`
	Origin        string                 `json:"origin,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}
