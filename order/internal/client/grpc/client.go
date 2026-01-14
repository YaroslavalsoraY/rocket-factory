package grpc

import (
	"context"

	"order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.Filters) ([]*model.PartInfo, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}
