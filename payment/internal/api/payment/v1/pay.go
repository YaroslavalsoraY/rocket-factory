package v1

import (
	"context"
	payment_v1 "shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	transactionID, err := a.paymentService.PayOrder(ctx)

	if err != nil {
		return nil, err
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionID,
	}, nil
}