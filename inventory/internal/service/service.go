package service

import (
	"context"
	"inventory/internal/model"
)

type InventoryService interface {
	Get(ctx context.Context, uuid string) (model.PartInfo, error)
	List(ctx context.Context, filters model.Filters) ([]model.PartInfo, error)
}