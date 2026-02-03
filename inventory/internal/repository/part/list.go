package part

import (
	"context"

	"inventory/internal/model"
	"inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

func (inv *inventory) ListParts(ctx context.Context, filters model.Filters) ([]model.PartInfo, error) {
	repoFilters := converter.ModelToRepoModelFilters(filters)

	cursor, err := inv.collection.Find(ctx, converter.RepoModelFiltersToMongoFilters(repoFilters))
	if err != nil {
		return nil, err
	}

	var parts []repoModel.PartInfo
	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, err
	}

	if len(parts) == 0 {
		return nil, model.ErrPartsNotFound
	}

	result := make([]model.PartInfo, len(parts))

	for i, part := range parts {
		result[i] = converter.RepoModelToModelPart(part)
	}

	return result, nil
}
