package order

import (
	"context"

	"order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUuid, userUuid string) (model.OrderInfo, error) {
	order, err := s.OrderRepository.GetOrder(ctx, orderUuid)
	if err != nil {
		return model.OrderInfo{}, err
	}

	if order.UserUUID != userUuid {
		return model.OrderInfo{}, model.ErrOrderNotFound
	}

	return order, nil
}
