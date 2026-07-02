package order_assembled_consumer

import (
	"context"
	def "notification/internal/service"

	"go.uber.org/zap"
	kafkaConverter "notification/internal/converter/kafka"
	"platform/pkg/kafka"
	"platform/pkg/logger"
)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
	telegramService        def.TelegramService
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssembledDecoder, telegramService def.TelegramService) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		telegramService:        telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderAssembledConsumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}

// TODO: Implement the logic for handling the order assembled event, сделать отправку уведомления в телеграм