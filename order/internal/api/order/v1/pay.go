package v1

import (
	"context"
	"errors"

	"order/internal/model"
	order_v1 "shared/pkg/openapi/order/v1"
)

func (a *OrderAPI) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	transactionUUID, err := a.OrderService.PayOrder(ctx, params.OrderUUID, model.PaymentMethod(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: model.ErrOrderNotFound.Error(),
			}, nil
		}

		return &order_v1.InternalServerError{
			Code:    500,
			Message: model.ErrInternal.Error(),
		}, nil
	}

	return &order_v1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
