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
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  []string{gofakeit.UUID()},
			TotalPrice: gofakeit.Float64Range(1, 99999),
			Status:     model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(testOrder, nil)
	s.orderRepository.On("UpdateOrderStatus", s.ctx, mock.MatchedBy(func(update model.OrderUpdateInfo) bool {
		return update.Status != nil && *update.Status == model.OrderStatusCANCELLED
	}), orderUUID).Return(nil)

	err := s.service.CancelOrder(s.ctx, orderUUID, userUUID)

	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderAnotherUser() {
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

	err := s.service.CancelOrder(s.ctx, orderUUID, userUUID)

	s.ErrorIs(err, model.ErrOrderNotFound)
	s.orderRepository.AssertNotCalled(s.T(), "UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestCancelOrderError() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		testOrder = model.OrderInfo{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  []string{gofakeit.UUID()},
			TotalPrice: gofakeit.Float64Range(1, 99999),
			Status:     model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(testOrder, model.ErrOrderNotFound)

	err := s.service.CancelOrder(s.ctx, orderUUID, userUUID)

	s.ErrorIs(err, model.ErrOrderNotFound)
}
