package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"order/internal/model"
)

func (s *ServiceSuite) TestGetOrder() {
	var (
		testOrder = model.OrderInfo{
			OrderUUID:         gofakeit.UUID(),
			UserUUID:          gofakeit.UUID(),
			PartUuids:         []string{gofakeit.UUID()},
			TotalPrice:        gofakeit.Float64Range(1, 99999),
			TransactionalUUID: gofakeit.UUID(),
			PaymentMethod:     model.PaymentMethodPAYMENTMETHODCREDITCARD,
			Status:            model.OrderStatusPAID,
		}
	)

	s.orderRepository.On("GetOrder", s.ctx, testOrder.OrderUUID).Return(testOrder, nil)

	order, err := s.service.GetOrder(s.ctx, testOrder.OrderUUID, testOrder.UserUUID)

	s.NoError(err)
	s.Equal(testOrder, order)
}

func (s *ServiceSuite) TestGetOrderAnotherUser() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		testOrder = model.OrderInfo{
			OrderUUID: orderUUID,
			UserUUID:  gofakeit.UUID(),
			Status:    model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(testOrder, nil)

	order, err := s.service.GetOrder(s.ctx, orderUUID, userUUID)

	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(order)
}
