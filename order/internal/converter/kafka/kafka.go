package kafka

import "order/internal/model"

type ShipAssembledDecoder interface {
	Decode([]byte) (model.ShipAssembled, error)
}
