package session

import "platform/pkg/cache"

type SessionRepository struct {
	redisClient cache.RedisClient
	setOperator cache.SetOperator
}

func NewSessionRepository(redisClient cache.RedisClient, setOperator cache.SetOperator) *SessionRepository {
	return &SessionRepository{
		redisClient: redisClient,
		setOperator: setOperator,
	}
}
