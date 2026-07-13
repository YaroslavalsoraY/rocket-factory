package user

import (
	"context"
	"iam/internal/model"
	"iam/internal/repository/converter"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (ur *UserRepository) Create(ctx context.Context, user model.User) (string, error) {
	repoUser := converter.ToRepoUser(user)

	userUUID := uuid.New().String()

	builderInsert := sq.Insert("users").
	Columns("id", "login", "password", "email", "notification_methods").
	PlaceholderFormat(sq.Dollar).
	Values(userUUID, repoUser.Login, repoUser.Password, repoUser.Email, repoUser.NotificationMethods)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return "", err
	}

	_, err = ur.pool.Exec(ctx, query, args...)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}