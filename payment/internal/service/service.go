package service

import "context"

type PaymentService interface {
	PayOrder(ctx context.Context) (string, error)
}
