package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *ServiceSuite) TestCancelOrder() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		testOrder = model.OrderInfo{
			OrderUUID:         orderUUID,
			UserUUID:          userUUID,
			PartUuids:         []string{gofakeit.UUID()},
			TotalPrice:        gofakeit.Float64Range(1, 99999),
			Status:            model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(testOrder, nil)
	s.orderRepository.On("CancelOrder", s.ctx, mock.Anything, orderUUID).Return(nil)

	err := s.service.CancelOrder(s.ctx, orderUUID)

	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderError() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		testOrder = model.OrderInfo{
			OrderUUID:         orderUUID,
			UserUUID:          userUUID,
			PartUuids:         []string{gofakeit.UUID()},
			TotalPrice:        gofakeit.Float64Range(1, 99999),
			Status:            model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(testOrder, model.ErrOrderNotFound)

	err := s.service.CancelOrder(s.ctx, orderUUID)

	s.ErrorIs(err, model.ErrOrderNotFound)
}
