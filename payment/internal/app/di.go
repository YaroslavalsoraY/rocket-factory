package app

import (
	"context"

	api "payment/internal/api/payment/v1"
	"payment/internal/service"
	"payment/internal/service/payment"
	payment_v1 "shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentService service.PaymentService

	paymentAPI payment_v1.PaymentServiceServer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = payment.NewService()
	}

	return d.paymentService
}

func (d *diContainer) PaymentAPI(ctx context.Context) payment_v1.PaymentServiceServer {
	if d.paymentAPI == nil {
		d.paymentAPI = api.NewApi(d.PaymentService(ctx))
	}

	return d.paymentAPI
}
