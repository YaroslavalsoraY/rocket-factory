package v1

import (
	"context"
	"errors"

	"order/internal/model"
	order_v1 "shared/pkg/openapi/order/v1"
	httpAuth "platform/pkg/middleware/http"
)

func (a *OrderAPI) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	user, ok := httpAuth.GetSessionFromContext(ctx)
    if !ok || user.UserId == "" {
		return &order_v1.UnauthorizedError{
			Code: 401,
			Message: model.ErrUnauthorized.Error(),
		}, nil
	}
	
	err := a.OrderService.CancelOrder(ctx, params.OrderUUID, user.UserId)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Message: model.ErrOrderNotFound.Error(),
				Code:    404,
			}, nil
		}

		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Message: model.ErrOrderNotFound.Error(),
				Code:    404,
			}, nil
		}

		return &order_v1.InternalServerError{
			Code:    500,
			Message: model.ErrInternal.Error(),
		}, nil
	}

	return &order_v1.CancelOrderNoContent{}, nil
}
