package assebmle_consumer

import (
	"context"
	"math/rand/v2"
	"time"

	"assembly/internal/model"
	"go.uber.org/zap"
	kafka "platform/pkg/kafka/consumer"
	"platform/pkg/logger"
)

const maxBuildTimeSec = 10

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	buildTimeSec := rand.Int64N(maxBuildTimeSec)

	time.Sleep(time.Duration(buildTimeSec * int64(time.Second)))

	err = s.shipAssembledProducer.ProduceShipAssembled(ctx, model.ShipAssembled{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: buildTimeSec,
	})
	if err != nil {
		logger.Error(ctx, "Failed to produce ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
