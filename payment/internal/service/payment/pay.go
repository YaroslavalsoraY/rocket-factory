package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
)

func (s *service) PayOrder(ctx context.Context) (string, error) {
	transactionID := uuid.New()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionID.String())

	return transactionID.String(), nil
}