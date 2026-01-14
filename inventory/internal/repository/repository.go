package repository

import (
	"context"

	"inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, uuid string) (model.PartInfo, error)
	ListParts(ctx context.Context, filetrs model.Filters) ([]model.PartInfo, error)
}
