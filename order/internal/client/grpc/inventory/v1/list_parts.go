package client

import (
	"context"

	"order/internal/client/converter"
	"order/internal/model"
	inventory_v1 "shared/pkg/proto/inventory/v1"

	grpcAuth "platform/pkg/middleware/grpc"
	"google.golang.org/grpc/metadata"
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
