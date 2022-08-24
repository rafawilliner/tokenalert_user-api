package users

import (
	"encoding/json"
)

type PublicUser struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	TelegramUser string `json:"telegram_user"`
	Status       string `json:"status"`
	DateCreated  string `json:"date_created"`
}

func (users Users) Marshall(isPublic bool) []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Marshall(isPublic)
	}
	return result
}

func (user *User) Marshall(isPublic bool) interface{} {
	userJson, _ := json.Marshal(user)
	var publicUser PublicUser
	json.Unmarshal(userJson, &publicUser)
	return publicUser
}
