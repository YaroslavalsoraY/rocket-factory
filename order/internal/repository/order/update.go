package order

import (
	"context"

	"order/internal/model"
	"order/internal/repository/converter"
	repoModel "order/internal/repository/model"
)

func (r *Repository) UpdateOrder(ctx context.Context, update model.OrderUpdateInfo, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	updateInfo := converter.FromServiceToRepoUpdateOrder(update)
	order, ok := r.storage[uuid]

	if !ok {
		return model.ErrOrderNotFound
	}

	if order.Status == repoModel.OrderStatusPAID {
		return model.ErrAlreadyPaid
	}

	if updateInfo.UserUUID != nil {
		order.UserUUID = *updateInfo.UserUUID
	}

	if len(updateInfo.PartUuids) != 0 {
		order.PartUuids = append(order.PartUuids, updateInfo.PartUuids...)
	}

	if updateInfo.PaymentMethod != nil {
		order.PaymentMethod = *updateInfo.PaymentMethod
	}

	if updateInfo.Status != nil {
		order.Status = *updateInfo.Status
	}

	if updateInfo.TotalPrice != nil {
		order.TotalPrice = *updateInfo.TotalPrice
	}

	if updateInfo.TransactionalUUID != nil {
		order.TransactionalUUID = *updateInfo.TransactionalUUID
	}

	r.storage[uuid] = order

	return nil
}
