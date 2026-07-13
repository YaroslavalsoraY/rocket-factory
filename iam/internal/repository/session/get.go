package session

import (
	"context"
	"encoding/json"
	"iam/internal/model"
	"iam/internal/repository/converter"
	repoModel "iam/internal/repository/model"
)

func (sr *SessionRepository) Get(ctx context.Context, sessionUUID string) (model.Session, error) {
	raw, err := sr.redisClient.Get(ctx, "sessions:"+sessionUUID)
	if err != nil {
		return model.Session{}, err
	}

	if len(raw) == 0 {
		return model.Session{}, model.ErrSessionNotFound
	}

	var session repoModel.Session
	err = json.Unmarshal(raw, &session)
	if err != nil {
		return model.Session{}, err
	}

	return converter.ToServiceSession(session), nil
}