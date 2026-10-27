package model

type Session struct {
	Login    string `json:"login"`
	UserUUID string `json:"userUUID"`
	Email    string `json:"email"`
}
