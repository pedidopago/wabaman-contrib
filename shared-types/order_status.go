package types

import (
	"encoding/json"
	"fmt"
)

type OrderStatus *string

var (
	orderStatusAssembling             = "assembling"
	orderStatusApproved               = "approved"
	orderStatusCreated                = "created"
	orderStatusFinalized              = "finalized"
	orderStatusPrescriptionCollection = "prescription_collection"
	orderStatusCancelled              = "cancelled"

	OrderStatusAssembling             OrderStatus = &orderStatusAssembling
	OrderStatusApproved               OrderStatus = &orderStatusApproved
	OrderStatusCreated                OrderStatus = &orderStatusCreated
	OrderStatusFinalized              OrderStatus = &orderStatusFinalized
	OrderStatusPrescriptionCollection OrderStatus = &orderStatusPrescriptionCollection
	OrderStatusCancelled              OrderStatus = &orderStatusCancelled
	OrderStatusDel                    OrderStatus = &delStr
)

var orderStatusIntern = map[string]OrderStatus{
	"assembling":              OrderStatusAssembling,
	"approved":                OrderStatusApproved,
	"created":                 OrderStatusCreated,
	"finalized":               OrderStatusFinalized,
	"prescription_collection": OrderStatusPrescriptionCollection,
	"cancelled":               OrderStatusCancelled,
	"$del":                    OrderStatusDel,
}

func internOrderStatus(raw json.RawMessage) (OrderStatus, error) {
	if len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' {
		return nil, fmt.Errorf("expected JSON string for OrderStatus, got %s", raw)
	}
	unquoted := raw[1 : len(raw)-1]
	if interned, ok := orderStatusIntern[string(unquoted)]; ok {
		return interned, nil
	}
	s := string(unquoted)
	return &s, nil
}

func EqualOrderStatus(a, b OrderStatus) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
