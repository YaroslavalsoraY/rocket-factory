package v1

import (
	"context"

	"inventory/internal/converter"
	inventory_v1 "shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	part, err := a.invService.Get(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventory_v1.GetPartResponse{
		Part: converter.PartInfoToProto(&part),
	}, nil
}
