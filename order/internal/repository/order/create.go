package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"order/internal/model"
	"order/internal/repository/converter"
)

func (r *Repository) CreateOrder(ctx context.Context, order model.OrderInfo) error {
	repoOrder := converter.FromServiceToRepoOrder(order)

	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).Columns("id", "user_uuid", "part_uuids", "total_price", "status").
		Values(repoOrder.OrderUUID, repoOrder.UserUUID, repoOrder.PartUuids, repoOrder.TotalPrice, repoOrder.Status)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
