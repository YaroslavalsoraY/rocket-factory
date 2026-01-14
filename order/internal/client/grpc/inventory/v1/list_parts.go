package client

import (
	"context"

	"order/internal/client/converter"
	"order/internal/model"
	inventory_v1 "shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.Filters) ([]*model.PartInfo, error) {
	parts, err := c.generatedClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: converter.ModelFiltersToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return converter.ArrayOfProtoToParts(parts.Parts), nil
}
