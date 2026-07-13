package iam_v1

import (
	"iam/internal/service"
	iam_v1 "shared/pkg/proto/iam/v1"
)

type IamApi struct {
	iam_v1.UnimplementedIamServiceServer

	sessionService service.IamService
	userService service.UserService
}

func NewIamApi(sessionService service.IamService, userService service.UserService) *IamApi{
	return &IamApi{
		sessionService: sessionService,
		userService: userService,
	}
}