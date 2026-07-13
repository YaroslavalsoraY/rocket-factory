package user

import (
	"context"
	"errors"
	"iam/internal/model"
	"platform/pkg/logger"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (us *UserService) GetUser(ctx context.Context, userUUID string) (model.User, error) {
	user, err := us.userRepository.Get(ctx, userUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		} else {
			return model.User{}, err
		}
	}

	logger.Info(ctx, "got user with uuid", zap.String("user_uuid", userUUID))

	return user, nil
}