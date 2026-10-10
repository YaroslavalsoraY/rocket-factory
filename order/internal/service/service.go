package service

import (
	"context"

	"order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partsUUIDs []string) (string, float32, error)
	GetOrder(ctx context.Context, orderUuid string, userUuid string) (model.OrderInfo, error)
	CancelOrder(ctx context.Context, orderUUID string, userUUID string) error
	PayOrder(ctx context.Context, orderUUID string, userUUID string, paymentMethod model.PaymentMethod) (string, error)
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error
}
