package kafka

import "assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembled, error)
}
