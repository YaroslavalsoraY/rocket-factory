package order

import (
	"context"
	clientMocks "order/internal/client/grpc/mocks"
	repositoryMock "order/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *repositoryMock.OrderRepository

	inventoryClient *clientMocks.InventoryClient

	paymentClient *clientMocks.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = repositoryMock.NewOrderRepository(s.T())

	s.inventoryClient = clientMocks.NewInventoryClient(s.T())

	s.paymentClient = clientMocks.NewPaymentClient(s.T())

	s.service = NewService(
		s.orderRepository, 
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {

}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}