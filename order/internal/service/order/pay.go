package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
	"platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "OrderService.PayOrder")
	defer span.End()

	order, err := s.GetOrder(ctx, orderUUID, userUUID)
	if err != nil {
		return "", err
	}

	if order.Status == model.OrderStatusPAID {
		return "", model.ErrAlreadyPaid
	}

	transactionID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, string(paymentMethod))
	if err != nil {
		return "", err
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

	s.revenue.Add(ctx, order.TotalPrice)

	s.orderPaidProducerService.ProduceOrderPaid(ctx, model.OrderPaid{
		EventUUID:       uuid.New().String(),
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   string(paymentMethod),
		TransactionUUID: transactionID,
	})

	return transactionID, nil
}
