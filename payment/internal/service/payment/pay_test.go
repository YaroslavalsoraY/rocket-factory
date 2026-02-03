package payment

func (s *ServiceSuite) TestPayOrder() {
	transactionID, err := s.service.PayOrder(s.ctx)

	s.NoError(err)
	s.NotEmpty(transactionID)
}
