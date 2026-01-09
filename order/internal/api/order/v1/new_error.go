package v1

import (
	"context"
	"net/http"
	order_v1 "shared/pkg/openapi/order/v1"
)

func (h *OrderAPI) NewError(ctx context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code: order_v1.NewOptInt(http.StatusInternalServerError),
			Message: order_v1.NewOptString(err.Error()),
		},
	}
}