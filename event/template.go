package event

import (
	"encoding/json"

	"github.com/pedidopago/wabaman-contrib/fbgraph"
)

type TemplateEvent struct {
	StoreID    string                  `json:"store_id"`
	BranchID   string                  `json:"branch_id,omitempty"`
	PhoneID    uint                    `json:"phone_id,omitempty"`
	Event      string                  `json:"event,omitempty"`
	TemplateID string                  `json:"template_id,omitempty"`
	Template   *fbgraph.TemplateObject `json:"template,omitempty"`
}

func (e TemplateEvent) ToJSON() string {
	d, _ := json.Marshal(e)
	return string(d)
}

func (e *TemplateEvent) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), e)
}
