package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *ServiceSuite) TestPayOrder() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		testOrder = model.OrderInfo{
			OrderUUID:         orderUUID,
			UserUUID:          userUUID,
			PartUuids:         []string{gofakeit.UUID()},
			TotalPrice:        gofakeit.Float64Range(1, 99999),
			TransactionalUUID: gofakeit.UUID(),
			PaymentMethod:     model.PaymentMethodPAYMENTMETHODCREDITCARD,
			Status:            model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("PayOrder", mock.Anything, mock.Anything, orderUUID).Return(nil)

	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(testOrder, nil)

	s.paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, string(model.PaymentMethodPAYMENTMETHODCARD)).Return(gofakeit.UUID(), nil)
	s.orderProducer.On("ProduceOrderPaid", mock.Anything, mock.MatchedBy(func(event model.OrderPaid) bool {
		return event.OrderUUID == orderUUID && event.UserUUID == userUUID &&
			event.PaymentMethod == string(model.PaymentMethodPAYMENTMETHODCARD) &&
			event.TransactionUUID != "" && event.EventUUID != ""
	})).Return(nil)

	transactionUUID, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, model.PaymentMethodPAYMENTMETHODCARD)

	s.NoError(err)
	s.NotEmpty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrderAnotherUser() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
		testOrder = model.OrderInfo{
			OrderUUID: orderUUID,
			UserUUID:  gofakeit.UUID(),
			Status:    model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("GetOrder", mock.Anything, orderUUID).Return(testOrder, nil)

	transactionUUID, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, model.PaymentMethodPAYMENTMETHODCARD)

	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(transactionUUID)
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	s.orderRepository.AssertNotCalled(s.T(), "PayOrder", mock.Anything, mock.Anything, mock.Anything)
	s.orderProducer.AssertNotCalled(s.T(), "ProduceOrderPaid", mock.Anything, mock.Anything)
}
