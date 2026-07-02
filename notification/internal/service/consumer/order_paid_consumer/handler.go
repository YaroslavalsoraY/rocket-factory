package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"
	kafka "platform/pkg/kafka/consumer"
	"platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order paid event "+event.EventUUID+" was caught")

	err = s.telegramService.SendOrderPaidNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send order paid notification", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order paid notification sent successfully for event "+event.EventUUID)

	return nil
}
