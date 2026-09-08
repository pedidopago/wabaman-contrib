package rest

import (
	"testing"

	"github.com/google/go-querystring/query"
)

// The contact_phone_number filter exists so a caller holding only a phone can
// find a BSUID contact, whose waba_contact_id is not a phone. Both encoders are
// checked because a field that never reaches the query string is precisely the
// failure this filter was added to fix: GetContactsRequest.BuildQuery is written
// by hand, so a new field is silently dropped unless it is emitted explicitly.

func TestGetContactsRequestBuildQueryEmitsContactPhoneNumbers(t *testing.T) {
	q := GetContactsRequest{
		BranchID:            "01GFPBDTDN6ZY0YEZSH8VFHPA4",
		ContactPhoneNumbers: []string{"5511987654321", "5521912345678"},
	}.BuildQuery()

	got := q["contact_phone_number"]
	if len(got) != 2 || got[0] != "5511987654321" || got[1] != "5521912345678" {
		t.Fatalf("contact_phone_number = %v, want both numbers", got)
	}
}

func TestGetContactsRequestBuildQueryOmitsEmptyContactPhoneNumbers(t *testing.T) {
	q := GetContactsRequest{BranchID: "01GFPBDTDN6ZY0YEZSH8VFHPA4"}.BuildQuery()

	if _, ok := q["contact_phone_number"]; ok {
		t.Fatalf("contact_phone_number present on a request that did not set it: %v", q)
	}
}

func TestGetContactsV2RequestEncodesContactPhoneNumbers(t *testing.T) {
	v, err := query.Values(GetContactsV2Request{
		BranchID:                 "01GFPBDTDN6ZY0YEZSH8VFHPA4",
		ContactPhoneNumbers:      []string{"5511987654321"},
		ExactContactPhoneNumbers: true,
	})
	if err != nil {
		t.Fatalf("query.Values: %v", err)
	}

	if got := v["contact_phone_number"]; len(got) != 1 || got[0] != "5511987654321" {
		t.Errorf("contact_phone_number = %v, want [5511987654321]", got)
	}
	if got := v.Get("exact_contact_phone_numbers"); got != "true" {
		t.Errorf("exact_contact_phone_numbers = %q, want \"true\"", got)
	}
}
