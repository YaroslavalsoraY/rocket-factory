package converter

import (
	"inventory/internal/model"
	inventory_v1 "shared/pkg/proto/inventory/v1"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartToPartInfo(part *inventory_v1.Part) *model.PartInfo {
	var dimensions *model.DimensionsInfo
	if part.Dimensions != nil {
		dimensions = &model.DimensionsInfo{
			Height: part.Dimensions.Height,
			Length: part.Dimensions.Length,
			Weight: part.Dimensions.Weight,
			Width: part.Dimensions.Width,
		}
	}

	var manufacturer *model.ManufacturerInfo
	if part.Manufacturer != nil {
		manufacturer = &model.ManufacturerInfo{
			Name: part.Manufacturer.Name,
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
		UUID: part.Uuid,
		Name: part.Name,
		Description: part.Description,
		Price: part.Price,
		StockQuantity: part.StockQuantity,
		Category: model.CategoryEnum(part.Category),
		Dimensions: dimensions,
		Manufacturer: manufacturer,
		Tags: part.Tags,
		Metadata: metadata,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}
}

func PartInfoToProto(partInfo *model.PartInfo) *inventory_v1.Part {
	var dimensions *inventory_v1.Dimensions
	if partInfo.Dimensions != nil {
		dimensions = &inventory_v1.Dimensions{
			Height: partInfo.Dimensions.Height,
			Length: partInfo.Dimensions.Length,
			Weight: partInfo.Dimensions.Weight,
			Width: partInfo.Dimensions.Width,
		}
	}

	var manufacturer *inventory_v1.Manufacturer
	if partInfo.Manufacturer != nil {
		manufacturer = &inventory_v1.Manufacturer{
			Name: partInfo.Manufacturer.Name,
			Country: partInfo.Manufacturer.Country,
			Website: partInfo.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*inventory_v1.Value)
	if len(partInfo.Metadata) != 0 {
		for k, v := range partInfo.Metadata {
			var value inventory_v1.Value
			switch (*v).(type) {
			case string:
				value = inventory_v1.Value{
					Value: &inventory_v1.Value_StringValue{
						StringValue: (*v).(string),
					},
				}
			case int64:
				value = inventory_v1.Value{
					Value: &inventory_v1.Value_Int64Value{
						Int64Value: (*v).(int64),
					},
				}
			case float64:
				value = inventory_v1.Value{
					Value: &inventory_v1.Value_DoubleValue{
						DoubleValue: (*v).(float64),
					},
				}
			case bool:
				value = inventory_v1.Value{
					Value: &inventory_v1.Value_BoolValue{
						BoolValue: (*v).(bool),
					},
				}
			}

			metadata[k] = &value
		}
	}

	var createdAt timestamppb.Timestamp
	if partInfo.CreatedAt != nil {
		createdAt.Nanos = int32(partInfo.CreatedAt.UnixNano())
	}

	var updatedAt timestamppb.Timestamp
	if partInfo.UpdatedAt != nil {
		updatedAt.Nanos = int32(partInfo.UpdatedAt.UnixNano())
	}
	
	return &inventory_v1.Part{
		Uuid: partInfo.UUID,
		Name: partInfo.Name,
		Description: partInfo.Description,
		Price: partInfo.Price,
		StockQuantity: partInfo.StockQuantity,
		Category: inventory_v1.Category(partInfo.Category),
		Dimensions: dimensions,
		Manufacturer: manufacturer,
		Tags: partInfo.Tags,
		Metadata: metadata,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}
}

func PartsFiltersToModelFilters(filters *inventory_v1.PartsFilter) *model.Filters {
	categories := make([]model.CategoryEnum, len(filters.Categories))
	if len(filters.Categories) != 0 {
		for i, category := range filters.Categories {
			categories[i] = model.CategoryEnum(category)
		}
	}

	return &model.Filters{
		UUIDs: filters.Uuids,
		Names: filters.Names,
		Categories: categories,
		ManufacturerCountries: filters.ManufacturerCountries,
		Tags: filters.Tags,
	}
}

func ModelFiltersToPartsFilters(filters *model.Filters) *inventory_v1.PartsFilter {
	categories := make([]inventory_v1.Category, len(filters.Categories))
	if len(filters.Categories) != 0 {
		for i, category := range filters.Categories {
			categories[i] = inventory_v1.Category(category)
		}
	}

	return &inventory_v1.PartsFilter{
		Uuids: filters.UUIDs,
		Names: filters.Names,
		Categories: categories,
		ManufacturerCountries: filters.ManufacturerCountries,
		Tags: filters.Tags,
	}
}

func ArrayOfPartsToProto(parts []model.PartInfo) []*inventory_v1.Part {
	result := make([]*inventory_v1.Part, len(parts))
	for i, part := range parts {
		result[i] = PartInfoToProto(&part)
	}

	return result
}