package order

import (
	"context"
	"order/internal/model"
)

func (s *Service) GetOrder(ctx context.Context, uuid string) (model.OrderInfo, error) {
	order, err := s.OrderRepository.GetOrder(ctx, uuid)
	if err != nil {
		return model.OrderInfo{}, err
	}

	return order, nil
}