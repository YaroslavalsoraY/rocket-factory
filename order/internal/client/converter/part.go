package converter

import (
	"order/internal/model"
	inventory_v1 "shared/pkg/proto/inventory/v1"
	"time"
)

func ModelFiltersToProto(filters model.Filters) *inventory_v1.PartsFilter {
	categories := make([]inventory_v1.Category, len(filters.Categories))
	if len(filters.Categories) != 0 {
		for i, category := range filters.Categories {
			categories[i] = inventory_v1.Category(category)
		}
	}

	return &inventory_v1.PartsFilter{
		Uuids:                 filters.UUIDs,
		Names:                 filters.Names,
		Categories:            categories,
		ManufacturerCountries: filters.ManufacturerCountries,
		Tags:                  filters.Tags,
	}
}

func PartToPartInfo(part *inventory_v1.Part) *model.PartInfo {
	var dimensions *model.DimensionsInfo
	if part.Dimensions != nil {
		dimensions = &model.DimensionsInfo{
			Height: part.Dimensions.Height,
			Length: part.Dimensions.Length,
			Weight: part.Dimensions.Weight,
			Width:  part.Dimensions.Width,
		}
	}

	var manufacturer *model.ManufacturerInfo
	if part.Manufacturer != nil {
		manufacturer = &model.ManufacturerInfo{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*any)
	if len(part.Metadata) != 0 {
		for k, v := range part.Metadata {
			var value any
			switch v.GetValue().(type) {
			case *inventory_v1.Value_StringValue:
				value = v.GetStringValue()
			case *inventory_v1.Value_Int64Value:
				value = v.GetInt64Value()
			case *inventory_v1.Value_DoubleValue:
				value = v.GetDoubleValue()
			case *inventory_v1.Value_BoolValue:
				value = v.GetBoolValue()
			}

			metadata[k] = &value
		}
	}

	var createdAt time.Time
	if part.CreatedAt != nil {
		createdAt = part.CreatedAt.AsTime()
	}

	var updatedAt time.Time
	if part.UpdatedAt != nil {
		updatedAt = part.UpdatedAt.AsTime()
	}

	return &model.PartInfo{
		UUID:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.CategoryEnum(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      metadata,
		CreatedAt:     &createdAt,
		UpdatedAt:     &updatedAt,
	}
}

func ArrayOfProtoToParts(parts []*inventory_v1.Part) []*model.PartInfo {
	result := make([]*model.PartInfo, len(parts))
	for i, part := range parts {
		result[i] = PartToPartInfo(part)
	}

	return result
}
