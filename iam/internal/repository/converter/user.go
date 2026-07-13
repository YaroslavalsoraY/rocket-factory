package converter

import (
	"iam/internal/model"
	repoModel "iam/internal/repository/model"
)

func ToRepoUser(user model.User) *repoModel.User {
	return &repoModel.User{
		ID: user.ID,
		Login: user.Login,
		Email: user.Email,
		Password: user.Password,
		NotificationMethods: user.NotificationMethods,
	}
}

func ToServiceUser(user repoModel.User) *model.User {
	return &model.User{
		ID: user.ID,
		Login: user.Login,
		Email: user.Email,
		Password: user.Password,
		NotificationMethods: user.NotificationMethods,
	}
}