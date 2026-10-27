package user_v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"iam/internal/converter"
	"iam/internal/model"
	user_v1 "shared/pkg/proto/user/v1"
)

func (ua *UserApi) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {
	user, err := ua.UserService.GetUser(ctx, req.UserUuid)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with id %s was not found", req.UserUuid)
		}

		return nil, status.Error(codes.Internal, model.ErrInternal.Error())
	}

	return &user_v1.GetUserResponse{User: converter.UserToProto(user)}, nil
}
