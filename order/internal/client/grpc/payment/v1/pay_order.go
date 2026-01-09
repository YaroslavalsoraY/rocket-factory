package client

import (
	"context"
	payment_v1 "shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	res, err := c.generatedClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid: orderUUID,
		UserUuid: userUUID,
		PaymentMethod: payment_v1.PaymentMethod(payment_v1.PaymentMethod_value[paymentMethod]),
	})
	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
}