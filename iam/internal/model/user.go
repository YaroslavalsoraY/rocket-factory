package model

type User struct {
	ID                  string
	Login               string
	Password            string
	Email               string
	NotificationMethods map[string]string
}
