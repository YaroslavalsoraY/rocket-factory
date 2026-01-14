package order

import (
	"order/internal/model"

	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestGetOrder() {
	var (
		uuid = gofakeit.UUID()
		testOrder = model.OrderInfo{
			OrderUUID:          uuid,
			UserUUID:          gofakeit.UUID(),
			PartUuids:   []string{gofakeit.UUID()},
			TotalPrice:         gofakeit.Float64Range(1, 99999),
			TransactionalUUID: gofakeit.UUID(),
			PaymentMethod: model.PaymentMethodPAYMENTMETHODCREDITCARD,
			Status: model.OrderStatusPAID,
		}
	)

	s.orderRepository.On("GetOrder", s.ctx, uuid).Return(testOrder, nil)

	order, err := s.service.GetOrder(s.ctx, uuid)

	s.NoError(err)
	s.Equal(testOrder, order)
}