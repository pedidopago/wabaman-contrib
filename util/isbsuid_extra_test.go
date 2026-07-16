package util

import "testing"

func TestIsBSUIDTightened(t *testing.T) {
	cases := map[string]bool{
		"US.13491208655302741918":     true,  // regular BSUID
		"BR.abc123":                   true,  // regular BSUID, other country
		"US.ENT.11815799212886844830": true,  // parent BSUID
		"5511987654321":               false, // phone
		"":                            false, // empty
		"US":                          false, // too short
		"jo.doe@example.com":          false, // email (has @)
		"[email":                      false, // stray char before dot no longer passes
		`a\.b`:                        false, // backslash not a letter
		"ab.cd":                       true,  // shape-valid (accepted; harmless in practice)
	}
	for in, want := range cases {
		if got := IsBSUID(in); got != want {
			t.Errorf("IsBSUID(%q) = %v, want %v", in, got, want)
		}
	}
}
