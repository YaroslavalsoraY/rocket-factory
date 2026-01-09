package part

import (
	repoModel "inventory/internal/repository/model"
	"sync"
)

type inventory struct {
	storage map[string]repoModel.PartInfo
	mu sync.RWMutex
}

func NewInventory() *inventory {
	storage := make(map[string]repoModel.PartInfo)
	storage["123e4567-e89b-12d3-a456-426614174000"] = repoModel.PartInfo{
		UUID:          "123e4567-e89b-12d3-a456-426614174000",
		Name:          "Тормозной диск передний",
		Description:   "Высококачественный вентилируемый тормозной диск для передних колес",
		Price:         12499.99,
		StockQuantity: 25,
		Category:      repoModel.CategoryEnum_CATEGORY_PORTHOLE,
		Tags:      []string{"тормоза", "диск", "передний", "вентилируемый"},
	}
	return &inventory{
		storage: storage,
	}
}