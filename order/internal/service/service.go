package service

import (
	"context"

	"order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partsUUIDs []string) (string, float32, error)
	GetOrder(ctx context.Context, uuid string) (model.OrderInfo, error)
	CancelOrder(ctx context.Context, uuid string) error
	PayOrder(ctx context.Context, uuid string, paymentMethod model.PaymentMethod) (string, error)
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error
}