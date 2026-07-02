package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"notification/internal/model"
	events_v1 "shared/pkg/proto/events/v1"
)

type orderAssembledDecoder struct{}

func NewOrderAssembledDecoder() *orderAssembledDecoder {
	return &orderAssembledDecoder{}
}

func (oad *orderAssembledDecoder) Decode(data []byte) (model.ShipAssembled, error) {
	var pb events_v1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembled{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembled{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
