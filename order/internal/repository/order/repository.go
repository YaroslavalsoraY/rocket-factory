package order

import (
	repoModel "order/internal/repository/model"
	"sync"
)

type Repository struct {
	storage map[string]repoModel.OrderInfo
	mu sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		storage: make(map[string]repoModel.OrderInfo),
	}
}