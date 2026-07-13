package converter

import (
	"iam/internal/model"
	repoModel "iam/internal/repository/model"
)

func ToServiceAuthInfo(authInfo repoModel.AuthInfo) model.AuthInfo {
	return model.AuthInfo{
		HashedPassword: authInfo.HashedPassword,
		UserUUID: authInfo.UserUUID,
		Email: authInfo.Email,
		Login: authInfo.Login,
	}
}