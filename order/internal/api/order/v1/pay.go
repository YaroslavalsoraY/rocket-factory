package v1

import (
	"context"
	"errors"

	"order/internal/model"
	order_v1 "shared/pkg/openapi/order/v1"
	httpAuth "platform/pkg/middleware/http"
)

func (a *OrderAPI) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	user, ok := httpAuth.GetSessionFromContext(ctx)
    if !ok || user.UserId == "" {
		return &order_v1.UnauthorizedError{
			Code: 401,
			Message: model.ErrUnauthorized.Error(),
		}, nil
	}
	
	transactionUUID, err := a.OrderService.PayOrder(ctx, params.OrderUUID, user.UserId, model.PaymentMethod(req.PaymentMethod))
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
