package model

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
	OrderUUID string `json:"order_uuid"`
	UserUUID string `json:"user_uuid"`
	PartUuids []string `json:"part_uuids"`
	TotalPrice float64 `json:"total_price"`
	TransactionalUUID string        `json:"transactional_uuid"`
	PaymentMethod     PaymentMethod `json:"payment_method"`
	Status            OrderStatus   `json:"status"`
}

type OrderUpdateInfo struct {
	UserUUID          *string
	PartUuids         []string
	TotalPrice        *float64
	TransactionalUUID *string
	PaymentMethod     *PaymentMethod
	Status            *OrderStatus
}