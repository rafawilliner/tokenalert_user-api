package users

import (
	"strings"

	"github.com/rafawilliner/tokenalert_utils-go/rest_errors"
)

type User struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	TelegramUser string `json:"telegram_user"`
	Status string `json:"status"`
	Password string `json:"password"`
}

func (user *User) Validate() rest_errors.RestErr {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.TelegramUser = strings.TrimSpace(strings.ToLower(user.TelegramUser))
	if user.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return rest_errors.NewBadRequestError("invalid password")
	}
	return nil
}