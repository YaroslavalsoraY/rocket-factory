package grpc

import (
	"context"
	iam_v1 "shared/pkg/proto/iam/v1"

	"google.golang.org/grpc"
)

type IamClient interface {
	Whoami(ctx context.Context, in *iam_v1.WhoamiRequest, opts ...grpc.CallOption) (*iam_v1.WhoamiResponse, error)
	Login(ctx context.Context, in *iam_v1.LoginRequest, opts ...grpc.CallOption) (*iam_v1.LoginResponse, error)
}