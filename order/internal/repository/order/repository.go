package order

import (
	"sync"

	repoModel "order/internal/repository/model"
)

type Repository struct {
	storage map[string]repoModel.OrderInfo
	mu      sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		storage: make(map[string]repoModel.OrderInfo),
	}
}
