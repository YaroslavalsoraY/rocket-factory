package user

import (
	"context"
	"iam/internal/model"

	"iam/internal/repository/converter"
	repoModel "iam/internal/repository/model"

	sq "github.com/Masterminds/squirrel"
)

func (ur *UserRepository) Get(ctx context.Context, userUUID string) (model.User, error) {
	builderSelect := sq.Select("id", "login", "email", "notification_methods").
	From("users").
	PlaceholderFormat(sq.Dollar).
	Where(sq.Eq{"id": userUUID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return model.User{}, err
	}

	row := ur.pool.QueryRow(ctx, query, args...)

	var user repoModel.User

	err = row.Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.NotificationMethods,
	)
	if err != nil {
		return model.User{}, err
	}

	return *converter.ToServiceUser(user), nil
}

func (ur *UserRepository) GetAuthInfo(ctx context.Context, login string) (model.AuthInfo, error) {
	builderSelect := sq.Select("email", "id", "password").
	From("users").
	PlaceholderFormat(sq.Dollar).
	Where(sq.Eq{"login": login})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return model.AuthInfo{}, err
	}
	row := ur.pool.QueryRow(ctx, query, args...)

	var authInfo repoModel.AuthInfo
	err = row.Scan(&authInfo.Email, &authInfo.UserUUID, &authInfo.HashedPassword)
	if err != nil {
		return model.AuthInfo{}, err
	}

	authInfo.Login = login

	return converter.ToServiceAuthInfo(authInfo), nil
}