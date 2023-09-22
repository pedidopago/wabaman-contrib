package util_test

import (
	"testing"

	"github.com/pedidopago/wabaman-contrib/util"
)

func TestCountAndValidateTemplateVariables(t *testing.T) {
	subjects := []struct {
		textContent   string
		expectedCount int
		expectedValid bool
	}{
		{"bolovo {{1}} pegue 2 pague {{2}}", 2, true},
		{"", 0, true},
		{"{", 0, true},
		{"}", 0, true},
		{"{}", 0, true},
		{"{a}", 0, true},
		{"{a}{b}", 0, true},
		{"hello { {{1}}", 1, true},
		{"hello { {{1}} }", 1, true},
		{"hello { {{1}} {{2}} }", 2, true},
		{"hello { {{3}} {{2}} {{1}} }", 0, false},
	}

	for i, item := range subjects {
		count, valid := util.CountAndValidateTemplateVariables(item.textContent)
		if count != item.expectedCount {
			t.Errorf("Test %d: Expected count %d, got %d", i, item.expectedCount, count)
		}
		if valid != item.expectedValid {
			t.Errorf("Test %d: Expected valid %t, got %t", i, item.expectedValid, valid)
		}
	}
}
