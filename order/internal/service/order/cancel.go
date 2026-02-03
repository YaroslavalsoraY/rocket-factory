package order

import (
	"context"

	"order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	order, err := s.OrderRepository.GetOrder(ctx, uuid)
	if err != nil {
		return err
	}

	if order.Status == model.OrderStatusPAID {
		return model.ErrAlreadyPaid
	}

	status := model.OrderStatusCANCELLED
	newInfo := model.OrderUpdateInfo{
		Status: &status,
	}

	err = s.OrderRepository.CancelOrder(ctx, newInfo, uuid)
	if err != nil {
		return err
	}

	return nil
}
