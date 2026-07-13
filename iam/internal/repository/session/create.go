package session

import (
	"context"
	"encoding/json"
	"iam/internal/config"
	"iam/internal/model"
	"iam/internal/repository/converter"

	"github.com/google/uuid"
)

func (sr *SessionRepository) Create(ctx context.Context, session model.Session) (string, error) {
	sessionData := converter.ToRepoSession(session)
	sessionToken := uuid.New().String()
	
	data, err := json.Marshal(sessionData)
	if err != nil {
		return "", err
	}
	err = sr.redisClient.SetWithTTL(ctx, "sessions:"+sessionToken, data, config.AppConfig().Session.TTL())
	if err != nil {
		return "", err
	}

	err = sr.AddSessionToUserSet(ctx, session.UserUUID, sessionToken)
	if err != nil {
		sr.redisClient.Del(ctx, "sessions:"+sessionToken)
		return "", err
	}

	return sessionToken, err
}