package service

import (
	"context"

	"notification/internal/model"
)

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, event model.OrderPaid) error
	SendShipAssembledNotification(ctx context.Context, event model.ShipAssembled) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
