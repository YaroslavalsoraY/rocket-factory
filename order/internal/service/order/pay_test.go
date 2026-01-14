package order

import (
	"order/internal/model"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestPayOrder() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID = gofakeit.UUID()
		testOrder = model.OrderInfo{
			OrderUUID:          orderUUID,
			UserUUID:          userUUID,
			PartUuids:   []string{gofakeit.UUID()},
			TotalPrice:         gofakeit.Float64Range(1, 99999),
			TransactionalUUID: gofakeit.UUID(),
			PaymentMethod: model.PaymentMethodPAYMENTMETHODCREDITCARD,
			Status: model.OrderStatusPAID,
		}
	)
	
	s.orderRepository.On("UpdateOrder", s.ctx, mock.Anything, orderUUID).Return(nil)

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(testOrder, nil)

	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, string(model.PaymentMethodPAYMENTMETHODCARD)).Return(gofakeit.UUID(), nil)

	transactionUUID, err := s.service.PayOrder(s.ctx, orderUUID, model.PaymentMethodPAYMENTMETHODCARD)

	s.NoError(err)
	s.NotEmpty(transactionUUID)
}