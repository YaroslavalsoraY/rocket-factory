package order

import (
	"order/internal/client/grpc"
	"order/internal/repository"
)

type Service struct {
	OrderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient grpc.PaymentClient
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) *Service {
	return &Service{
		OrderRepository: orderRepository,

		inventoryClient: inventoryClient,
		paymentClient: paymentClient,
	}
}