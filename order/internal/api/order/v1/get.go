package v1

import (
	"context"
	"errors"

	"order/internal/converter"
	"order/internal/model"
	httpAuth "platform/pkg/middleware/http"
	order_v1 "shared/pkg/openapi/order/v1"
)

func (a *OrderAPI) GetOrder(ctx context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	user, ok := httpAuth.GetSessionFromContext(ctx)
    if !ok || user.UserId == "" {
		return &order_v1.UnauthorizedError{
			Code: 401,
			Message: model.ErrUnauthorized.Error(),
		}, nil
	}
	
	order, err := a.OrderService.GetOrder(ctx, params.OrderUUID, user.UserId)
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

	return converter.FromServiceToHttpOrder(order), nil
}
