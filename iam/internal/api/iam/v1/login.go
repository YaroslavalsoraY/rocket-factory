package iam_v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iam/internal/model"
	iam_v1 "shared/pkg/proto/iam/v1"
)

func (iam *IamApi) Login(ctx context.Context, req *iam_v1.LoginRequest) (*iam_v1.LoginResponse, error) {
	sessionUUID, err := iam.sessionService.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, model.ErrUserNotFound.Error())
		}
		if errors.Is(err, model.ErrUnauthorized) {
			return nil, status.Error(codes.Unauthenticated, model.ErrUnauthorized.Error())
		}

		return nil, status.Error(codes.Internal, model.ErrInternal.Error())
	}

	return &iam_v1.LoginResponse{SessionUuid: sessionUUID}, nil
}
