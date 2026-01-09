package v1

import (
	"inventory/internal/service"
	inventory_v1 "shared/pkg/proto/inventory/v1"
)

type api struct {
	inventory_v1.UnimplementedInventoryServiceServer

	invService service.InventoryService
}

func NewApi(invService service.InventoryService) *api {
	return &api{
		invService: invService,
	}
}