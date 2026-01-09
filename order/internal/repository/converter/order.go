package converter

import (
	"order/internal/model"
	repoModel "order/internal/repository/model"
)

func FromServiceToRepoOrder(order model.OrderInfo) *repoModel.OrderInfo{
	return &repoModel.OrderInfo{
		OrderUUID: order.OrderUUID,
		UserUUID: order.UserUUID,
		PartUuids: order.PartUuids,
		TotalPrice: order.TotalPrice,
		TransactionalUUID: order.TransactionalUUID,
		PaymentMethod: repoModel.PaymentMethod(order.PaymentMethod),
		Status: repoModel.OrderStatus(order.Status),
	}
}

func FromRepoToServiceOrder(order repoModel.OrderInfo) *model.OrderInfo {
	return &model.OrderInfo{
		OrderUUID: order.OrderUUID,
		UserUUID: order.UserUUID,
		PartUuids: order.PartUuids,
		TotalPrice: order.TotalPrice,
		TransactionalUUID: order.TransactionalUUID,
		PaymentMethod: model.PaymentMethod(order.PaymentMethod),
		Status: model.OrderStatus(order.Status),
	}
}

func FromServiceToRepoUpdateOrder(order model.OrderUpdateInfo) *repoModel.OrderUpdateInfo {
	return &repoModel.OrderUpdateInfo{
		UserUUID: order.UserUUID,
		PartUuids: order.PartUuids,
		TotalPrice: order.TotalPrice,
		TransactionalUUID: order.TransactionalUUID,
		PaymentMethod: (*repoModel.PaymentMethod)(order.PaymentMethod),
		Status: (*repoModel.OrderStatus)(order.Status),
	}
}