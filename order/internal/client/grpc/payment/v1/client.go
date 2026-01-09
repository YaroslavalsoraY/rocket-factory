package client

import (
	payment_v1 "shared/pkg/proto/payment/v1"
)

type client struct {
	generatedClient payment_v1.PaymentServiceClient
}

func NewClient(generatedClient payment_v1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}