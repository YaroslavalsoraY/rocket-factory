package model

import "time"

type PartInfo struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      CategoryEnum
	Dimensions    *DimensionsInfo
	Manufacturer  *ManufacturerInfo
	Tags          []string
	Metadata      map[string]any
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type Filters struct {
	UUIDs                 []string
	Names                 []string
	Categories            []CategoryEnum
	ManufacturerCountries []string
	Tags                  []string
}

type ManufacturerInfo struct {
	Name    string
	Country string
	Website string
}

type DimensionsInfo struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type CategoryEnum int32

const (
	CategoryEnum_CATEGORY_UNKNOWN_UNSPECIFIED CategoryEnum = 0
	CategoryEnum_CATEGORY_ENGINE              CategoryEnum = 1
	CategoryEnum_CATEGORY_FUEL                CategoryEnum = 2
	CategoryEnum_CATEGORY_PORTHOLE            CategoryEnum = 3
	CategoryEnum_CATEGORY_WING                CategoryEnum = 4
)
