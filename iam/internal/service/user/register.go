package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"iam/internal/model"
	"platform/pkg/logger"
)

func (us *UserService) Register(ctx context.Context, user model.User) (string, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return "", model.ErrPasswordTooLong
		}
		return "", err
	}

	user.Password = hashedPassword

	userUUID, err := us.userRepository.Create(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", model.ErrUserAlreadyExists
		}
		return "", err
	}

	logger.Info(ctx, "Created user with uuid", zap.String("user_uuid", userUUID))

	return userUUID, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
