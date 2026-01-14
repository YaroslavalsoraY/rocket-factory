package v1

import (
	"context"
	"errors"

	"order/internal/model"
	order_v1 "shared/pkg/openapi/order/v1"
)

func (a *OrderAPI) CreateNewOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateNewOrderRes, error) {
	UUID, totalPrice, err := a.OrderService.CreateOrder(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: model.ErrPartsNotFound.Error(),
			}, nil
		}

		return &order_v1.InternalServerError{
			Code:    500,
			Message: model.ErrInternal.Error(),
		}, nil
	}

	return &order_v1.CreateOrderResponse{
		OrderUUID:  UUID,
		TotalPrice: totalPrice,
	}, nil
}
