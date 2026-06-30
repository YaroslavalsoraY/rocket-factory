package order_consumer

import (
	"context"

	kafkaConverter "order/internal/converter/kafka"
	"order/internal/repository"
	"platform/pkg/kafka"
	"platform/pkg/logger"

	"go.uber.org/zap"
)

type service struct {
	shipAssembledConsumer kafka.Consumer
	shipAssembledDecoder  kafkaConverter.ShipAssembledDecoder
	orderRepository       repository.OrderRepository
}

func NewService(shipAssembledConsumer kafka.Consumer, shipAssembledDecoder kafkaConverter.ShipAssembledDecoder, orderRepository repository.OrderRepository) *service {
	return &service{
		shipAssembledConsumer: shipAssembledConsumer,
		shipAssembledDecoder:  shipAssembledDecoder,
		orderRepository: orderRepository,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderPaidConsumer service")

	err := s.shipAssembledConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
