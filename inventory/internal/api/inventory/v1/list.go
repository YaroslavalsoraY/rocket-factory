package v1

import (
	"context"
	"inventory/internal/converter"
	inventory_v1 "shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	parts, err := a.invService.List(ctx, *converter.PartsFiltersToModelFilters(req.Filter))
	if err != nil {
		return nil, err
	}

	return &inventory_v1.ListPartsResponse{
		Parts: converter.ArrayOfPartsToProto(parts),
	}, nil
}