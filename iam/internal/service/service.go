package service

import (
	"context"
	"iam/internal/model"
)

type IamService interface {
	Login(ctx context.Context, login, password string) (string, error)
	WhoAmI(ctx context.Context, sessionUUID string) (model.Session, error)
}

type UserService interface {	
	Register(ctx context.Context, user model.User) (string, error)
	GetUser(ctx context.Context, userUUID string) (model.User, error)
}
