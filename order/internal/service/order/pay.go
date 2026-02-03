package order

import (
	"context"

	"order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, uuid string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.GetOrder(ctx, uuid)
	if err != nil {
		return "", err
	}

	if order.Status == model.OrderStatusPAID {
		return "", model.ErrAlreadyPaid
	}

	transactionID, err := s.paymentClient.PayOrder(ctx, uuid, order.UserUUID, string(paymentMethod))
	if err != nil {
		return "", nil
	}

	status := model.OrderStatusPAID
	newInfo := model.OrderUpdateInfo{
		TransactionalUUID: &transactionID,
		PaymentMethod:     &paymentMethod,
		Status:            &status,
	}

	err = s.OrderRepository.PayOrder(ctx, newInfo, uuid)
	if err != nil {
		return "", err
	}

	return transactionID, nil
}
