package converter

import (
	"iam/internal/model"
	repoModel "iam/internal/repository/model"
)

func ToRepoSession(session model.Session) *repoModel.Session {
	return &repoModel.Session{
		UserUUID: session.UserUUID,
		Email:    session.Email,
		Login:    session.Login,
	}
}

func ToServiceSession(session repoModel.Session) model.Session {
	return model.Session{
		UserUUID: session.UserUUID,
		Email:    session.Email,
		Login:    session.Login,
	}
}
