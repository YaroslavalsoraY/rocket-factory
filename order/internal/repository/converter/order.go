package converter

import (
	"database/sql"

	"order/internal/model"
	repoModel "order/internal/repository/model"
)

func FromServiceToRepoOrder(order model.OrderInfo) *repoModel.OrderInfo {
	var transactionUUIDValid, paymentMethodValid bool
	if order.TransactionalUUID != "" {
		transactionUUIDValid = true
	}
	if order.PaymentMethod != "" {
		paymentMethodValid = true
	}

	return &repoModel.OrderInfo{
		OrderUUID:         order.OrderUUID,
		UserUUID:          order.UserUUID,
		PartUuids:         order.PartUuids,
		TotalPrice:        order.TotalPrice,
		TransactionalUUID: sql.NullString{String: order.TransactionalUUID, Valid: transactionUUIDValid},
		PaymentMethod:     sql.NullString{String: string(order.PaymentMethod), Valid: paymentMethodValid},
		Status:            repoModel.OrderStatus(order.Status),
	}
}

func FromRepoToServiceOrder(order repoModel.OrderInfo) *model.OrderInfo {
	return &model.OrderInfo{
		OrderUUID:         order.OrderUUID,
		UserUUID:          order.UserUUID,
		PartUuids:         order.PartUuids,
		TotalPrice:        order.TotalPrice,
		TransactionalUUID: order.TransactionalUUID.String,
		PaymentMethod:     model.PaymentMethod(order.PaymentMethod.String),
		Status:            model.OrderStatus(order.Status),
	}
}

func FromServiceToRepoUpdateOrder(order model.OrderUpdateInfo) *repoModel.OrderUpdateInfo {
	return &repoModel.OrderUpdateInfo{
		UserUUID:          order.UserUUID,
		PartUuids:         order.PartUuids,
		TotalPrice:        order.TotalPrice,
		TransactionalUUID: order.TransactionalUUID,
		PaymentMethod:     (*repoModel.PaymentMethod)(order.PaymentMethod),
		Status:            (*repoModel.OrderStatus)(order.Status),
	}
}
