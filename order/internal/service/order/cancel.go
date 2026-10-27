package order

import (
	"context"

	"order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID, userUUID string) error {
	order, err := s.GetOrder(ctx, orderUUID, userUUID)
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

	err = s.OrderRepository.UpdateOrderStatus(ctx, newInfo, orderUUID)
	if err != nil {
		return err
	}

	return nil
}
