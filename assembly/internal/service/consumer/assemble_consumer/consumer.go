package assebmle_consumer

import (
	"context"

	kafkaConverter "assembly/internal/converter/kafka"
	"platform/pkg/kafka"
	"platform/pkg/logger"

	"go.uber.org/zap"
	def "assembly/internal/service"
)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	shipAssembledProducer def.ProducerService 
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, shipAssembledProducer def.ProducerService ) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		shipAssembledProducer: shipAssembledProducer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
