package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"order/internal/model"
	"order/internal/repository/converter"
	repoModel "order/internal/repository/model"
)

func (r *Repository) GetOrder(ctx context.Context, uuid string) (model.OrderInfo, error) {
	builderSelect := sq.Select("id", "user_uuid", "part_uuids", "total_price", "transactional_uuid", "payment_method", "status").
		From("orders").PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": uuid})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return model.OrderInfo{}, err
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var order repoModel.OrderInfo

	err = row.Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&order.PartUuids,
		&order.TotalPrice,
		&order.TransactionalUUID,
		&order.PaymentMethod,
		&order.Status,
	)
	if err != nil {
		return model.OrderInfo{}, err
	}

	return *converter.FromRepoToServiceOrder(order), nil
}
