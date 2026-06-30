package service

import (
	"context"

	"assembly/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error
}
