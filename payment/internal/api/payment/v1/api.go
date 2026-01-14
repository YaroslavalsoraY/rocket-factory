package v1

import (
	"payment/internal/service"
	payment_v1 "shared/pkg/proto/payment/v1"
)

type api struct {
	payment_v1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewApi(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
