package assebmle_consumer

import (
	"context"

	kafkaConverter "assembly/internal/converter/kafka"
	def "assembly/internal/service"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
	"platform/pkg/kafka"
	"platform/pkg/logger"
)

type service struct {
	orderPaidConsumer     kafka.Consumer
	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
	shipAssembledProducer def.ProducerService

	shipAssemblingTime metric.Int64Histogram
}

func NewService(orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	shipAssembledProducer def.ProducerService,
	shipAssemblingTime metric.Int64Histogram,
) *service {
	return &service{
		orderPaidConsumer:     orderPaidConsumer,
		orderPaidDecoder:      orderPaidDecoder,
		shipAssembledProducer: shipAssembledProducer,
		shipAssemblingTime:    shipAssemblingTime,
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
