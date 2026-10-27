package converter

import (
	"iam/internal/model"
	common_v1 "shared/pkg/proto/common/v1"
	user_v1 "shared/pkg/proto/user/v1"
)

func UserFromProto(protoUser *common_v1.User) model.User {
	notificationMethods := make(map[string]string)
	for _, el := range protoUser.NotificationMethods {
		notificationMethods[el.ProviderName] = el.Target
	}

	return model.User{
		ID:                  protoUser.Id,
		Login:               protoUser.Login,
		Password:            protoUser.Password,
		Email:               protoUser.Email,
		NotificationMethods: notificationMethods,
	}
}

func UserToProto(user model.User) *user_v1.User {
	notificationMethods := make([]*common_v1.NotificationMethod, 0)
	for key, val := range user.NotificationMethods {
		notificationMethods = append(notificationMethods, &common_v1.NotificationMethod{ProviderName: key, Target: val})
	}

	return &user_v1.User{
		Id:                  user.ID,
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: notificationMethods,
	}
}
