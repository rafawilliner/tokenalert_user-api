package services

import (
	"tokenalert_user-api/domain/users"
	"tokenalert_user-api/repositories"
	"tokenalert_user-api/utils/crypto_utils"
	"tokenalert_user-api/utils/date_utils"

	"github.com/rafawilliner/tokenalert_utils-go/rest_errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, rest_errors.RestErr)
	GetUser(int64) (*users.User, rest_errors.RestErr)
	LoginUser(users.LoginRequest) (*users.User, rest_errors.RestErr)
}

func (s *usersService) CreateUser(user users.User) (*users.User, rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := repositories.UsersRepository.Save(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) GetUser(userId int64) (*users.User, rest_errors.RestErr) {

	var user *users.User
	var err rest_errors.RestErr
	if user, err = repositories.UsersRepository.Get(userId); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *usersService) LoginUser(request users.LoginRequest) (*users.User, rest_errors.RestErr) {
	/*dao := &users.User{
		Email:    request.Email,
		Password: crypto_utils.GetMd5(request.Password),
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}*/
	return nil, nil
}
