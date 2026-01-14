package v1

import "order/internal/service"

type OrderAPI struct {
	OrderService service.OrderService
}

func NewApi(orderService service.OrderService) *OrderAPI {
	return &OrderAPI{
		OrderService: orderService,
	}
}
