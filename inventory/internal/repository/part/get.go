package part

import (
	"context"
	"inventory/internal/model"
	repoConverter "inventory/internal/repository/converter"
)

func (inv *inventory) GetPart(ctx context.Context, uuid string) (model.PartInfo, error) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	part, ok := inv.storage[uuid]
	if !ok {
		return model.PartInfo{}, model.ErrPartsNotFound
	}

	return repoConverter.RepoModelToModelPart(part), nil
}
