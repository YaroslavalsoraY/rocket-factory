package client

import (
	iam_v1 "shared/pkg/proto/iam/v1"
)

type client struct {
	iam_v1.IamServiceClient
}

func NewClient(generatedClient iam_v1.IamServiceClient) *client {
	return &client{
		generatedClient,
	}
}
