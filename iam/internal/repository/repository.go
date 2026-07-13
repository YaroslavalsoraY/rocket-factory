package repository

import (
	"context"
	"iam/internal/model"
)

type SessionRepository interface {
	AddSessionToUserSet(ctx context.Context, userUUID string, sessionToken string) error
	GetUserSessions(ctx context.Context, userUUID string) ([]model.Session, error)
	Create(ctx context.Context, session model.Session) (string, error)
	Get(ctx context.Context, sessionUUID string) (model.Session, error)
}

type UserRepository interface {
	Create(ctx context.Context, user model.User) (string, error)
	Get(ctx context.Context, userUUID string) (model.User, error)
	GetAuthInfo(ctx context.Context, login string) (model.AuthInfo, error)
}