package payment

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"platform/pkg/logger"
)

func (s *service) PayOrder(ctx context.Context) (string, error) {
	transactionID := uuid.New()

	logger.Info(ctx, "Оплата прошла успешно", zap.String("transaction ID", transactionID.String()))

	return transactionID.String(), nil
}
