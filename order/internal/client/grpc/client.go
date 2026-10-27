package grpc

import (
	"context"

	"google.golang.org/grpc"
	"order/internal/model"
	iam_v1 "shared/pkg/proto/iam/v1"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.Filters) ([]*model.PartInfo, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}

type IamClient interface {
	Whoami(ctx context.Context, in *iam_v1.WhoamiRequest, opts ...grpc.CallOption) (*iam_v1.WhoamiResponse, error)
	Login(ctx context.Context, in *iam_v1.LoginRequest, opts ...grpc.CallOption) (*iam_v1.LoginResponse, error)
}
