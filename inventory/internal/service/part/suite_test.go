package part

import (
	"context"
	"inventory/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context
	InvRepository *mocks.InventoryRepository
	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.InvRepository = mocks.NewInventoryRepository(s.T())

	s.service = NewService(
		s.InvRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {

}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}