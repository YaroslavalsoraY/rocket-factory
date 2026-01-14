package part

import (
	"context"
	"slices"

	"inventory/internal/model"
	"inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

func (inv *inventory) ListParts(ctx context.Context, filters model.Filters) ([]model.PartInfo, error) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	repoFilters := converter.ModelToRepoModelFilters(filters)
	result := make([]model.PartInfo, 0)

	if IsNoFilters(repoFilters) {
		for _, v := range inv.storage {
			result = append(result, converter.RepoModelToModelPart(v))
		}
		return result, nil
	}

	for _, v := range inv.storage {
		if isInFilters(repoFilters, v) {
			result = append(result, converter.RepoModelToModelPart(v))
		}
	}

	if len(result) == 0 {
		return result, model.ErrPartsNotFound
	}

	return result, nil
}

func IsNoFilters(filters repoModel.Filters) bool {
	if len(filters.Names) == 0 &&
		len(filters.UUIDs) == 0 &&
		len(filters.Categories) == 0 &&
		len(filters.ManufacturerCountries) == 0 &&
		len(filters.Tags) == 0 {
		return true
	}
	return false
}

func isInFilters(filters repoModel.Filters, part repoModel.PartInfo) bool {
	if (filters.UUIDs == nil || slices.Contains(filters.UUIDs, part.UUID) || len(filters.UUIDs) == 0) &&
		(filters.Names == nil || slices.Contains(filters.Names, part.Name) || len(filters.Names) == 0) &&
		(filters.Categories == nil || slices.Contains(filters.Categories, part.Category) || len(filters.Categories) == 0) &&
		(filters.ManufacturerCountries == nil || slices.Contains(filters.ManufacturerCountries, part.Manufacturer.Country) || len(filters.ManufacturerCountries) == 0) &&
		isInTags(filters.Tags, part.Tags) {
		return true
	}
	return false
}

func isInTags(filterTags, partTags []string) bool {
	if len(filterTags) == 0 || filterTags == nil {
		return true
	}
	for _, filterTag := range partTags {
		if !slices.Contains(partTags, filterTag) {
			return false
		}
	}
	return true
}
