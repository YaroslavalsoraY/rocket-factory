package part

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"inventory/internal/model"
	repoConverter "inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

func (inv *inventory) GetPart(ctx context.Context, uuid string) (model.PartInfo, error) {
	var part repoModel.PartInfo

	err := inv.collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&part)
	if err != nil {
		return model.PartInfo{}, model.ErrPartsNotFound
	}

	return repoConverter.RepoModelToModelPart(part), nil
}
