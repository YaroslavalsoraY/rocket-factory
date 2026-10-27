package user_v1

import (
	"iam/internal/service"
	user_v1 "shared/pkg/proto/user/v1"
)

type UserApi struct {
	user_v1.UnimplementedUserServiceServer

	UserService service.UserService
}

func NewUserApi(userService service.UserService) *UserApi {
	return &UserApi{
		UserService: userService,
	}
}
