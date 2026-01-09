package order

import (
	"context"
	"order/internal/model"
	"order/internal/repository/converter"
)

func (r *Repository) GetOrder(ctx context.Context, uuid string) (model.OrderInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.storage[uuid]
	if !ok {
		return model.OrderInfo{}, model.ErrOrderNotFound
	}

	return *converter.FromRepoToServiceOrder(order), nil
}
