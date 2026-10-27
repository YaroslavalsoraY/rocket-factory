package client

import (
	"context"

	"google.golang.org/grpc/metadata"
	"order/internal/client/converter"
	"order/internal/model"
	grpcAuth "platform/pkg/middleware/grpc"
	inventory_v1 "shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.Filters) ([]*model.PartInfo, error) {
	sessionUUID, ok := grpcAuth.GetSessionUUIDFromContext(ctx)
	if !ok {
		return nil, model.ErrUnauthorized
	}
	ctx = metadata.AppendToOutgoingContext(ctx, grpcAuth.SessionUUIDMetadataKey, sessionUUID)
	parts, err := c.generatedClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: converter.ModelFiltersToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return converter.ArrayOfProtoToParts(parts.Parts), nil
}
