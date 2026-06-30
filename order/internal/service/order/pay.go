package order

import (
	"context"
	"github.com/google/uuid"
	"order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.GetOrder(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	if order.Status == model.OrderStatusPAID {
		return "", model.ErrAlreadyPaid
	}

	transactionID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, string(paymentMethod))
	if err != nil {
		return "", nil
	}

	status := model.OrderStatusPAID
	newInfo := model.OrderUpdateInfo{
		TransactionalUUID: &transactionID,
		PaymentMethod:     &paymentMethod,
		Status:            &status,
	}

	err = s.OrderRepository.PayOrder(ctx, newInfo, orderUUID)
	if err != nil {
		return "", err
	}

	s.orderPaidProducerService.ProduceOrderPaid(ctx, model.OrderPaid{
		EventUUID: uuid.New().String(),
		OrderUUID: order.OrderUUID,
		UserUUID: order.UserUUID,
		PaymentMethod: string(paymentMethod),
		TransactionUUID: transactionID,
	})

	return transactionID, nil
}
