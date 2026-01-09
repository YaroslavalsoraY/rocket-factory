package part

import "inventory/internal/repository"

type service struct {
	invRepo repository.InventoryRepository
}

func NewService(invRepo repository.InventoryRepository) *service {
	return &service{
		invRepo: invRepo,
	}
}