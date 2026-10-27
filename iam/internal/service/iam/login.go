package iam

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"iam/internal/model"
)

func (iam *IamService) Login(ctx context.Context, login, password string) (string, error) {
	authInfo, err := iam.userRepository.GetAuthInfo(ctx, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", model.ErrUserNotFound
		}
		return "", err
	}

	if compareHashWithPassword(authInfo.HashedPassword, password) {
		return iam.sessionRepository.Create(ctx, model.Session{
			UserUUID: authInfo.UserUUID,
			Login:    authInfo.Login,
			Email:    authInfo.Email,
		})
	}

	return "", model.ErrUnauthorized
}

func compareHashWithPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}
