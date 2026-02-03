package part

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/mongo"
	repoModel "inventory/internal/repository/model"
)

type inventory struct {
	collection *mongo.Collection
}

func NewInventory(db *mongo.Database) *inventory {
	collection := db.Collection("parts")

	_, err := collection.InsertOne(context.Background(), repoModel.PartInfo{
		UUID:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   "Test part",
		Price:         999.9,
		StockQuantity: 15,
		Category:      repoModel.CategoryEnum_CATEGORY_ENGINE,
		Dimensions: &repoModel.DimensionsInfo{
			Length: 1.0,
			Width:  2.0,
			Weight: 3.0,
			Height: 4.0,
		},
		Manufacturer: &repoModel.ManufacturerInfo{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      []string{"Test", "Fake"},
		Metadata:  map[string]any{"power": int64(500), "model": "sosalik", "is_kaif": true, "kaif_percents": float64(89.234)},
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil
	}

	return &inventory{
		collection: collection,
	}
}
