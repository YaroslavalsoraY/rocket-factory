package order_producer

import (
	"context"

	"assembly/internal/model"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"platform/pkg/kafka"
	"platform/pkg/logger"
	eventsV1 "shared/pkg/proto/events/v1"
)

type service struct {
	shipAssembledProducer kafka.Producer
}

func NewService(shipAssembledProducer kafka.Producer) *service {
	return &service{
		shipAssembledProducer: shipAssembledProducer,
	}
}

func (p *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal ShipAssembled", zap.Error(err))
		return err
	}

	err = p.shipAssembledProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
