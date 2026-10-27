package model

import "errors"

var (
	ErrSessionNotFound   = errors.New("session was not found")
	ErrUserNotFound      = errors.New("user was not found")
	ErrInternal          = errors.New("Internal error")
	ErrUserAlreadyExists = errors.New("User with this login or email is already exists")
	ErrPasswordTooLong   = errors.New("password can't be longer than 72 symbols")
	ErrUnauthorized      = errors.New("unauthorized")
)
