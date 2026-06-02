package types

import (
	"encoding/json"
	"fmt"
)

type InquiryStatus *string

// InquiryStatus interned sentinels.
// Use pointer comparison for fast checks: cm.InquiryStatus == InquiryStatusAvailable
// Use EqualInquiryStatus for safe comparison when values may not come from unmarshal.
var (
	inquiryStatusAvailable    = "AVAILABLE"
	inquiryStatusCancelled    = "CANCELLED"
	inquiryStatusDone         = "DONE"
	inquiryStatusExpired      = "EXPIRED"
	inquiryStatusInAttendance = "IN_ATTENDANCE"
	inquiryStatusOrdered      = "ORDERED"
	inquiryStatusPending      = "PENDING"

	InquiryStatusAvailable    InquiryStatus = &inquiryStatusAvailable
	InquiryStatusCancelled    InquiryStatus = &inquiryStatusCancelled
	InquiryStatusDone         InquiryStatus = &inquiryStatusDone
	InquiryStatusExpired      InquiryStatus = &inquiryStatusExpired
	InquiryStatusInAttendance InquiryStatus = &inquiryStatusInAttendance
	InquiryStatusOrdered      InquiryStatus = &inquiryStatusOrdered
	InquiryStatusPending      InquiryStatus = &inquiryStatusPending
	InquiryStatusDel          InquiryStatus = &delStr
)

var inquiryStatusIntern = map[string]InquiryStatus{
	"AVAILABLE":     InquiryStatusAvailable,
	"CANCELLED":     InquiryStatusCancelled,
	"DONE":          InquiryStatusDone,
	"EXPIRED":       InquiryStatusExpired,
	"IN_ATTENDANCE": InquiryStatusInAttendance,
	"ORDERED":       InquiryStatusOrdered,
	"PENDING":       InquiryStatusPending,
	"$del":          InquiryStatusDel,
}

func internInquiryStatus(raw json.RawMessage) (InquiryStatus, error) {
	if len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' {
		return nil, fmt.Errorf("expected JSON string for InquiryStatus, got %s", raw)
	}
	unquoted := raw[1 : len(raw)-1]
	if interned, ok := inquiryStatusIntern[string(unquoted)]; ok {
		return interned, nil
	}
	s := string(unquoted)
	return &s, nil
}

func EqualInquiryStatus(a, b InquiryStatus) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
