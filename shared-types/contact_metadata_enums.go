package types

import "encoding/json"

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
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	if interned, ok := inquiryStatusIntern[s]; ok {
		return interned, nil
	}
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
