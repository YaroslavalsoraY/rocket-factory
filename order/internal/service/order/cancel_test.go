package order

import (
	"order/internal/model"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCancelOrder() {
	var testUUID = gofakeit.UUID()
	
	s.orderRepository.On("UpdateOrder", s.ctx, mock.Anything, testUUID).Return(nil)

	err := s.service.CancelOrder(s.ctx, testUUID)

	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderError() {
	var testUUID = gofakeit.UUID()
	
	s.orderRepository.On("UpdateOrder", s.ctx, mock.Anything, testUUID).Return(model.ErrOrderNotFound)

	err := s.service.CancelOrder(s.ctx, testUUID)

	s.ErrorIs(err, model.ErrOrderNotFound)
}