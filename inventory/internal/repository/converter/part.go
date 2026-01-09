package converter

import (
	repoModel "inventory/internal/repository/model"
	"inventory/internal/model"
)

func ModelToRepoModelPart(part model.PartInfo) repoModel.PartInfo {
	var dimensions *repoModel.DimensionsInfo
	if part.Dimensions != nil {
		dimensions = &repoModel.DimensionsInfo{
			Length: part.Dimensions.Length,
			Width: part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		}
	}

	var manufacturer *repoModel.ManufacturerInfo
	if part.Manufacturer != nil {
		manufacturer = &repoModel.ManufacturerInfo{
			Name: part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}
	
	return repoModel.PartInfo{
		UUID: part.UUID,
		Name: part.Name,
		Description: part.Description,
		Price: part.Price,
		StockQuantity: part.StockQuantity,
		Category: repoModel.CategoryEnum(part.Category),
		Dimensions: dimensions,
		Manufacturer: manufacturer,
		Tags: part.Tags,
		Metadata: part.Metadata,
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
}

func RepoModelToModelPart(partInfo repoModel.PartInfo) model.PartInfo {
	var dimensions *model.DimensionsInfo
	if partInfo.Dimensions != nil {
		dimensions = &model.DimensionsInfo{
			Height: partInfo.Dimensions.Height,
			Length: partInfo.Dimensions.Length,
			Weight: partInfo.Dimensions.Weight,
			Width: partInfo.Dimensions.Width,
		}
	}

	var manufacturer *model.ManufacturerInfo
	if partInfo.Manufacturer != nil {
		manufacturer = &model.ManufacturerInfo{
			Name: partInfo.Manufacturer.Name,
			Country: partInfo.Manufacturer.Country,
			Website: partInfo.Manufacturer.Website,
		}
	}
	
	return model.PartInfo{
		UUID: partInfo.UUID,
		Name: partInfo.Name,
		Description: partInfo.Description,
		Price: partInfo.Price,
		StockQuantity: partInfo.StockQuantity,
		Category: model.CategoryEnum(partInfo.Category),
		Dimensions: dimensions,
		Manufacturer: manufacturer,
		Tags: partInfo.Tags,
		Metadata: partInfo.Metadata,
		CreatedAt: partInfo.CreatedAt,
		UpdatedAt: partInfo.UpdatedAt,
	}
}

func ModelToRepoModelFilters(filters model.Filters) repoModel.Filters {
	categories := make([]repoModel.CategoryEnum, len(filters.Categories))
	if len(filters.Categories) != 0 {
		for i, category := range filters.Categories {
			categories[i] = repoModel.CategoryEnum(category)
		}
	}

	return repoModel.Filters{
		UUIDs: filters.UUIDs,
		Names: filters.Names,
		Categories: categories,
		ManufacturerCountries: filters.ManufacturerCountries,
		Tags: filters.Tags,
	}
}

func RepoModelFiltersToModelFilters(filters repoModel.Filters) model.Filters {
	categories := make([]model.CategoryEnum, len(filters.Categories))
	if len(filters.Categories) != 0 {
		for i, category := range filters.Categories {
			categories[i] = model.CategoryEnum(category)
		}
	}

	return model.Filters{
		UUIDs: filters.UUIDs,
		Names: filters.Names,
		Categories: categories,
		ManufacturerCountries: filters.ManufacturerCountries,
		Tags: filters.Tags,
	}
}