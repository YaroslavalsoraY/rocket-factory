package converter

import "notification/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembled, error)
}
