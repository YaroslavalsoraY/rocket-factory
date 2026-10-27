package iam

import (
	"context"
	"errors"

	"github.com/gomodule/redigo/redis"
	"iam/internal/model"
)

func (iam *IamService) WhoAmI(ctx context.Context, sessionUUID string) (model.Session, error) {
	session, err := iam.sessionRepository.Get(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return model.Session{}, model.ErrSessionNotFound
		}
		return model.Session{}, err
	}

	return session, nil
}
