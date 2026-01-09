package order

import (
	"context"
	"order/internal/model"

	"github.com/google/uuid"
)

func (s *Service) CreateOrder(ctx context.Context, userUUID string, partsUUIDs []string) (string, float32, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.Filters{UUIDs: partsUUIDs})
	if err != nil {
		return "", 0, err
	}

	if len(parts) != len(partsUUIDs) {
		return "", 0, model.ErrPartsNotFound
	}

	var totalPrice float32

	for _, part := range parts {
		totalPrice += float32(part.Price)
	}
	orderUUID := uuid.New().String()

	order := model.OrderInfo{
		OrderUUID: orderUUID,
		UserUUID: userUUID,
		PartUuids: partsUUIDs,
		Status: model.OrderStatusPENDINGPAYMENT,
		TotalPrice: float64(totalPrice),
	}

	err = s.OrderRepository.CreateOrder(ctx, order)
	if err != nil {
		return "", 0, err
	}

	return orderUUID, totalPrice, nil
}