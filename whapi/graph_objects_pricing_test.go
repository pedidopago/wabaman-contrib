package whapi

import (
	"encoding/json"
	"testing"
)

func TestStatusPricingObjectUnmarshal(t *testing.T) {
	tests := []struct {
		name         string
		raw          string
		wantCategory PricingCategory
		wantType     PricingType
		wantBillable bool
		wantIsNil    bool
	}{
		{
			name:         "charged regular message",
			raw:          `{"billable":true,"pricing_model":"PMP","type":"regular","category":"marketing"}`,
			wantCategory: PricingCategoryMarketing,
			wantType:     PricingTypeRegular,
			wantBillable: true,
		},
		{
			name:         "utility free inside customer service window",
			raw:          `{"billable":false,"pricing_model":"PMP","type":"free_customer_service","category":"utility"}`,
			wantCategory: PricingCategoryUtility,
			wantType:     PricingTypeFreeCustomerService,
			wantBillable: false,
		},
		{
			name:         "free entry point",
			raw:          `{"billable":false,"pricing_model":"PMP","type":"free_entry_point","category":"service"}`,
			wantCategory: PricingCategoryService,
			wantType:     PricingTypeFreeEntryPoint,
			wantBillable: false,
		},
		{
			// Payloads predating the `billable` field: fall back to `type`.
			name:         "legacy payload without billable",
			raw:          `{"pricing_model":"PMP","type":"regular","category":"utility"}`,
			wantCategory: PricingCategoryUtility,
			wantType:     PricingTypeRegular,
			wantBillable: true,
			wantIsNil:    true,
		},
		{
			// Neither field present: not billable, and the caller can tell it is
			// unknown because Billable is nil and Type is empty.
			name:         "legacy payload without billable or type",
			raw:          `{"pricing_model":"CBP","category":"utility"}`,
			wantCategory: PricingCategoryUtility,
			wantType:     "",
			wantBillable: false,
			wantIsNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p StatusPricingObject
			if err := json.Unmarshal([]byte(tt.raw), &p); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if p.Category != tt.wantCategory {
				t.Errorf("Category = %q, want %q", p.Category, tt.wantCategory)
			}
			if p.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", p.Type, tt.wantType)
			}
			if got := p.IsBillable(); got != tt.wantBillable {
				t.Errorf("IsBillable() = %v, want %v", got, tt.wantBillable)
			}
			if (p.Billable == nil) != tt.wantIsNil {
				t.Errorf("Billable nil = %v, want %v", p.Billable == nil, tt.wantIsNil)
			}
		})
	}
}

func TestStatusPricingObjectIsBillableNil(t *testing.T) {
	var p *StatusPricingObject
	if p.IsBillable() {
		t.Error("nil pricing object must not report billable")
	}
}
