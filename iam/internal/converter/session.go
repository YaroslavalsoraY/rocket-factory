package converter

import (
	"iam/internal/model"
	common_v1 "shared/pkg/proto/common/v1"
)

func SessionToProto(session model.Session) *common_v1.Session {
	return &common_v1.Session{
		UserId: session.UserUUID,
		Login: session.Login,
		Email: session.Email,
	}
}