package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"
	kafka "platform/pkg/kafka/consumer"
	"platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order assembled event "+event.EventUUID+" was caught")

	err = s.telegramService.SendShipAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send ship assembled notification", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Ship assembled notification sent successfully for event "+event.EventUUID)

	return nil
}
