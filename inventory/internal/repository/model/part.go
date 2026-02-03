package model

import "time"

type PartInfo struct {
	UUID          string            `bson:"_id"`
	Name          string            `bson:"name"`
	Description   string            `bson:"desription"`
	Price         float64           `bson:"price"`
	StockQuantity int64             `bson:"stock_quantity"`
	Category      CategoryEnum      `bson:"category"`
	Dimensions    *DimensionsInfo   `bson:"dimensions,omitempty"`
	Manufacturer  *ManufacturerInfo `bson:"manufacturer,omitempty"`
	Tags          []string          `bson:"tags"`
	Metadata      map[string]any    `bson:"metadata"`
	CreatedAt     time.Time         `bson:"created_at"`
	UpdatedAt     *time.Time        `bson:"updated_at,omitempty"`
}

type Filters struct {
	UUIDs                 []string
	Names                 []string
	Categories            []CategoryEnum
	ManufacturerCountries []string
	Tags                  []string
}

type ManufacturerInfo struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type DimensionsInfo struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type CategoryEnum int32

const (
	CategoryEnum_CATEGORY_UNKNOWN_UNSPECIFIED CategoryEnum = 0
	CategoryEnum_CATEGORY_ENGINE              CategoryEnum = 1
	CategoryEnum_CATEGORY_FUEL                CategoryEnum = 2
	CategoryEnum_CATEGORY_PORTHOLE            CategoryEnum = 3
	CategoryEnum_CATEGORY_WING                CategoryEnum = 4
)
