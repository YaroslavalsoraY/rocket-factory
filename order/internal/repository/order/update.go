package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"order/internal/model"
	"order/internal/repository/converter"
)

// func (r *Repository) UpdateOrder(ctx context.Context, update model.OrderUpdateInfo, uuid string) error {

// 	updateInfo := converter.FromServiceToRepoUpdateOrder(update)
// 	order, ok := r.storage[uuid]

// 	if !ok {
// 		return model.ErrOrderNotFound
// 	}

// 	if order.Status == repoModel.OrderStatusPAID {
// 		return model.ErrAlreadyPaid
// 	}

// 	if updateInfo.UserUUID != nil {
// 		order.UserUUID = *updateInfo.UserUUID
// 	}

// 	if len(updateInfo.PartUuids) != 0 {
// 		order.PartUuids = append(order.PartUuids, updateInfo.PartUuids...)
// 	}

// 	if updateInfo.PaymentMethod != nil {
// 		order.PaymentMethod = *updateInfo.PaymentMethod
// 	}

// 	if updateInfo.Status != nil {
// 		order.Status = *updateInfo.Status
// 	}

// 	if updateInfo.TotalPrice != nil {
// 		order.TotalPrice = *updateInfo.TotalPrice
// 	}

// 	if updateInfo.TransactionalUUID != nil {
// 		order.TransactionalUUID = *updateInfo.TransactionalUUID
// 	}

// 	return nil
// }

func (r *Repository) PayOrder(ctx context.Context, update model.OrderUpdateInfo, uuid string) error {
	newInfo := converter.FromServiceToRepoUpdateOrder(update)

	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("transactional_uuid", *newInfo.TransactionalUUID).
		Set("payment_method", *newInfo.PaymentMethod).
		Set("status", *newInfo.Status).
		Where(sq.Eq{"id": uuid})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CancelOrder(ctx context.Context, update model.OrderUpdateInfo, uuid string) error {
	newInfo := converter.FromServiceToRepoUpdateOrder(update)

	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("status", *newInfo.Status).
		Where(sq.Eq{"id": uuid})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
