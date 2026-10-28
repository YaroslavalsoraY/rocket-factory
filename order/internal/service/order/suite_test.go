package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/metric/noop"
	clientMocks "order/internal/client/grpc/mocks"
	repositoryMock "order/internal/repository/mocks"
	kafkaMocks "order/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *repositoryMock.OrderRepository

	inventoryClient *clientMocks.InventoryClient

	paymentClient *clientMocks.PaymentClient

	orderProducer *kafkaMocks.ProducerService

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = repositoryMock.NewOrderRepository(s.T())

	s.inventoryClient = clientMocks.NewInventoryClient(s.T())

	s.paymentClient = clientMocks.NewPaymentClient(s.T())

	s.orderProducer = kafkaMocks.NewProducerService(s.T())

	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
		s.orderProducer,
		noop.Int64Counter{},
		noop.Float64Counter{},
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
