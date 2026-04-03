package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"platform/pkg/logger"
)

const (
	testLogLvl = "info"
	testAsJson = false
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	err := logger.Init(testLogLvl, testAsJson)
	if err != nil {
		return
	}

	s.service = NewService()
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
