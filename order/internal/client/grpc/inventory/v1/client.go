package client

import inventory_v1 "shared/pkg/proto/inventory/v1"

type client struct {
	generatedClient inventory_v1.InventoryServiceClient
}

func NewClient(generatedClient inventory_v1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
