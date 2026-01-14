package repository

import (
	"context"

	"order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.OrderInfo) error
	GetOrder(ctx context.Context, uuid string) (model.OrderInfo, error)
	UpdateOrder(ctx context.Context, update model.OrderUpdateInfo, uuid string) error
}
