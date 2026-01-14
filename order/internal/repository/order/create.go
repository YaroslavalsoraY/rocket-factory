package order

import (
	"context"

	"order/internal/model"
	"order/internal/repository/converter"
)

func (r *Repository) CreateOrder(ctx context.Context, order model.OrderInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	repoOrder := converter.FromServiceToRepoOrder(order)

	r.storage[repoOrder.OrderUUID] = *repoOrder

	return nil
}
