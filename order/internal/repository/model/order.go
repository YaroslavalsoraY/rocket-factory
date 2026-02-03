package model

import "database/sql"

type PaymentMethod string

const (
	PaymentMethodPAYMENTMETHODUNKNOWNUNSPECIFIED PaymentMethod = "PAYMENT_METHOD_UNKNOWN_UNSPECIFIED"
	PaymentMethodPAYMENTMETHODCARD               PaymentMethod = "PAYMENT_METHOD_CARD"
	PaymentMethodPAYMENTMETHODSBP                PaymentMethod = "PAYMENT_METHOD_SBP"
	PaymentMethodPAYMENTMETHODCREDITCARD         PaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	PaymentMethodPAYMENTMETHODINVESTORMONEY      PaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type OrderInfo struct {
	OrderUUID         string
	UserUUID          string
	PartUuids         []string
	TotalPrice        float64
	TransactionalUUID sql.NullString
	PaymentMethod     sql.NullString
	Status            OrderStatus
}

type OrderUpdateInfo struct {
	UserUUID          *string
	PartUuids         []string
	TotalPrice        *float64
	TransactionalUUID *string
	PaymentMethod     *PaymentMethod
	Status            *OrderStatus
}
