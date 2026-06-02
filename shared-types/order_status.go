package types

import "encoding/json"

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
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	if interned, ok := orderStatusIntern[s]; ok {
		return interned, nil
	}
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
