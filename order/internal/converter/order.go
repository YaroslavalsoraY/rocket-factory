package converter

import (
	"order/internal/model"
	order_v1 "shared/pkg/openapi/order/v1"
)

func FromHttpToServiceOrder(order order_v1.OrderDto) *model.OrderInfo {
	return &model.OrderInfo{
		OrderUUID:         order.OrderUUID,
		UserUUID:          order.UserUUID,
		PartUuids:         order.PartUuids,
		TotalPrice:        order.TotalPrice,
		TransactionalUUID: order.TransactionalUUID,
		PaymentMethod:     model.PaymentMethod(order.PaymentMethod),
		Status:            model.OrderStatus(order.Status),
	}
}

func FromServiceToHttpOrder(order model.OrderInfo) *order_v1.OrderDto {
	return &order_v1.OrderDto{
		OrderUUID:         order.OrderUUID,
		UserUUID:          order.UserUUID,
		PartUuids:         order.PartUuids,
		TotalPrice:        order.TotalPrice,
		TransactionalUUID: order.TransactionalUUID,
		PaymentMethod:     order_v1.PaymentMethod(order.PaymentMethod),
		Status:            order_v1.OrderStatus(order.Status),
	}
}
