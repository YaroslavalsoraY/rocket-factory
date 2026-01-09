package v1

import (
	"context"
	"order/internal/model"
	order_v1 "shared/pkg/openapi/order/v1"
)

func (a *OrderAPI) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	err := a.OrderService.CancelOrder(ctx, params.OrderUUID)

	if err != nil {
		switch err {
		case model.ErrOrderNotFound:
			return &order_v1.NotFoundError{
				Message: model.ErrOrderNotFound.Error(),
				Code:    404,
			}, nil
		case model.ErrAlreadyPaid:
			return &order_v1.ConflictError{
				Message: model.ErrAlreadyPaid.Error(),
				Code:    409,
			}, nil
		default:
			return &order_v1.InternalServerError{
				Code:    500,
				Message: model.ErrInternal.Error(),
			}, nil
		}
	}

	return &order_v1.CancelOrderNoContent{}, nil
}
