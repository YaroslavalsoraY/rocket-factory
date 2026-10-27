package session

import (
	"context"

	"iam/internal/model"
)

func (sr *SessionRepository) GetUserSessions(ctx context.Context, userUUID string) ([]model.Session, error) {
	indexKey := "user_sessions:" + userUUID

	tokens, err := sr.redisClient.SMembers(ctx, indexKey)
	if err != nil {
		return nil, err
	}

	var sessions []model.Session
	var toRemove []string

	for _, token := range tokens {
		exists, err := sr.redisClient.Exists(ctx, "sessions:"+token)
		if err != nil {
			continue
		}
		if !exists {
			toRemove = append(toRemove, token)
			continue
		}

		session, err := sr.Get(ctx, token)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}

	if len(toRemove) > 0 {
		sr.setOperator.SRem(ctx, indexKey, toRemove...)
	}

	return sessions, nil
}
