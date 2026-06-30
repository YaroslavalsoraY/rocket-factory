package decoder

import (
	"fmt"

	"assembly/internal/model"
	"google.golang.org/protobuf/proto"
	events_v1 "shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaid, error) {
	var pb events_v1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaid{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaid{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   pb.PaymentMethod,
		TransactionUUID: pb.TransactionUuid,
	}, nil
}
