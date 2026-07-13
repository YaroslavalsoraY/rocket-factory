package user_v1

import (
	"context"
	"errors"
	"iam/internal/converter"
	"iam/internal/model"
	user_v1 "shared/pkg/proto/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ua *UserApi) Register(ctx context.Context, req *user_v1.RegisterRequest) (*user_v1.RegisterResponse, error) {
	userUUID, err := ua.UserService.Register(ctx, converter.UserFromProto(req.User))
	if err != nil {
		if errors.Is(err, model.ErrPasswordTooLong) {
			return nil, status.Error(codes.InvalidArgument, model.ErrPasswordTooLong.Error())
		}
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, model.ErrUserAlreadyExists.Error())
		}
		return nil, status.Error(codes.Internal, model.ErrInternal.Error())
	}

	return &user_v1.RegisterResponse{UserUuid: userUUID}, nil
}