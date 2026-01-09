package part

import (
	"context"
	"inventory/internal/model"
)

func (s *service) List(ctx context.Context, filters model.Filters) ([]model.PartInfo, error) {
	parts, err := s.invRepo.ListParts(ctx, filters)
	if err != nil {
		return parts, err
	}
	
	return parts, nil
}