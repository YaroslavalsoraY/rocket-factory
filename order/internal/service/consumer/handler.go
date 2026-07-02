package order_consumer

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"order/internal/model"
	kafka "platform/pkg/kafka/consumer"
	"platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Ship was built for "+strconv.Itoa(int(event.BuildTimeSec))+" seconds")

	newStatus := model.OrderStatusCOMPLETED
	s.orderRepository.UpdateOrderStatus(ctx, model.OrderUpdateInfo{
		Status: &newStatus,
	}, event.OrderUUID)

	return nil
}
