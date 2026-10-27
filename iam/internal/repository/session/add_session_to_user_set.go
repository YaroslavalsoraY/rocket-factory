package session

import (
	"context"
)

func (sr *SessionRepository) AddSessionToUserSet(ctx context.Context, userUUID, sessionToken string) error {
	err := sr.setOperator.SAdd(ctx, "user_sessions:"+userUUID, sessionToken)
	if err != nil {
		return err
	}

	return nil
}
