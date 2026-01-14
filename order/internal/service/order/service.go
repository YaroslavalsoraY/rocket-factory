package order

import (
	"order/internal/client/grpc"
	"order/internal/repository"
)

type service struct {
	OrderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) *service {
	return &service{
		OrderRepository: orderRepository,

		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
