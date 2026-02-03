package converter

import (
	"go.mongodb.org/mongo-driver/bson"
	"inventory/internal/model"
	repoModel "inventory/internal/repository/model"
)

func ModelToRepoModelPart(part model.PartInfo) repoModel.PartInfo {
	dimensions := (*repoModel.DimensionsInfo)(part.Dimensions)

	manufacturer := (*repoModel.ManufacturerInfo)(part.Manufacturer)

	return repoModel.PartInfo{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repoModel.CategoryEnum(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func RepoModelToModelPart(partInfo repoModel.PartInfo) model.PartInfo {
	dimensions := (*model.DimensionsInfo)(partInfo.Dimensions)

	manufacturer := (*model.ManufacturerInfo)(partInfo.Manufacturer)

	return model.PartInfo{
		UUID:          partInfo.UUID,
		Name:          partInfo.Name,
		Description:   partInfo.Description,
		Price:         partInfo.Price,
		StockQuantity: partInfo.StockQuantity,
		Category:      model.CategoryEnum(partInfo.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          partInfo.Tags,
		Metadata:      partInfo.Metadata,
		CreatedAt:     partInfo.CreatedAt,
		UpdatedAt:     partInfo.UpdatedAt,
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
		UUIDs:                 filters.UUIDs,
		Names:                 filters.Names,
		Categories:            categories,
		ManufacturerCountries: filters.ManufacturerCountries,
		Tags:                  filters.Tags,
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
		UUIDs:                 filters.UUIDs,
		Names:                 filters.Names,
		Categories:            categories,
		ManufacturerCountries: filters.ManufacturerCountries,
		Tags:                  filters.Tags,
	}
}

func RepoModelFiltersToMongoFilters(filters repoModel.Filters) bson.M {
	conditions := make([]bson.M, 0)

	if len(filters.UUIDs) > 0 {
		conditions = append(conditions, bson.M{"_id": bson.M{"$in": filters.UUIDs}})
	}

	if len(filters.Names) > 0 {
		conditions = append(conditions, bson.M{"name": bson.M{"$in": filters.Names}})
	}

	if len(filters.Categories) > 0 {
		conditions = append(conditions, bson.M{"category": bson.M{"$in": filters.Categories}})
	}

	if len(filters.ManufacturerCountries) > 0 {
		conditions = append(conditions, bson.M{"manufacturer.country": bson.M{"$in": filters.ManufacturerCountries}})
	}

	if len(filters.Tags) > 0 {
		conditions = append(conditions, bson.M{"tags": bson.M{"$in": filters.Tags}})
	}

	if len(conditions) == 0 {
		return bson.M{}
	}

	return bson.M{
		"$and": conditions,
	}
}
