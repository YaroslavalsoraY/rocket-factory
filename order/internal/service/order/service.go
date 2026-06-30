package order

import (
	"order/internal/client/grpc"
	"order/internal/repository"
	def "order/internal/service"
)

type service struct {
	OrderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	orderPaidProducerService def.ProducerService
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
	producerService def.ProducerService,
) *service {
	return &service{
		OrderRepository: orderRepository,

		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,

		orderPaidProducerService: producerService,
	}
}
