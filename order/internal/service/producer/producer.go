package producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"order/internal/model"
	"platform/pkg/kafka"
	"platform/pkg/logger"
	eventsV1 "shared/pkg/proto/events/v1"
)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error {
	msg := &eventsV1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaid", zap.Error(err))
		return err
	}

	err = p.orderPaidProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderPaid", zap.Error(err))
		return err
	}

	return nil
}
