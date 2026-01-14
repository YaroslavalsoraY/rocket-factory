package order

import (
	"context"

	"order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	status := model.OrderStatusCANCELLED
	newInfo := model.OrderUpdateInfo{
		Status: &status,
	}

	err := s.OrderRepository.UpdateOrder(ctx, newInfo, uuid)
	if err != nil {
		return err
	}

	return nil
}
