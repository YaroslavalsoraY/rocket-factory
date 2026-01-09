package part

import (
	"context"
	"inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.PartInfo, error) {
	part, err := s.invRepo.GetPart(ctx, uuid)
	if err != nil {
		return model.PartInfo{}, err
	}

	return part, nil
}