package order

import (
	"go.opentelemetry.io/otel/metric"
	"order/internal/client/grpc"
	"order/internal/repository"
	def "order/internal/service"
)

type service struct {
	OrderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	orderPaidProducerService def.ProducerService

	orders     metric.Int64Counter
	revenue    metric.Float64Counter
	assembling metric.Float64Histogram
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
	producerService def.ProducerService,
	ordersCounter metric.Int64Counter,
	ordersRevenue metric.Float64Counter,
) *service {
	return &service{
		OrderRepository: orderRepository,

		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,

		orderPaidProducerService: producerService,

		orders:  ordersCounter,
		revenue: ordersRevenue,
	}
}
