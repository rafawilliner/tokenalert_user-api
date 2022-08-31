package services

import (
	"errors"
	"testing"
	"tokenalert_user-api/src/domain/users"
	"tokenalert_user-api/src/repositories"

	"github.com/rafawilliner/tokenalert_utils-go/rest_errors"
	"github.com/stretchr/testify/assert"
)

var (
	createUserRepoFunc func(user *users.User) rest_errors.RestErr
	getUserRepoFunc func(int64) (*users.User, rest_errors.RestErr)
)

type usersRepoMock struct{}

func (*usersRepoMock) Save(user *users.User) rest_errors.RestErr {
	return createUserRepoFunc(user)
}

func (*usersRepoMock) Get(int64) (*users.User, rest_errors.RestErr) {
	return nil, nil
}

func TestCreateOK(t *testing.T) {

	user := users.User{Id: 666, Name: "John", Email: "john@mail.com", Password: "admin"}
	createUserRepoFunc = func(user *users.User) rest_errors.RestErr {
		return nil
	}

	repositories.UsersRepository = &usersRepoMock{}

	_, err := UsersService.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, int64(666), user.Id)
}

func TestCreateMissingPasswordReturnBadRequest(t *testing.T) {
	user := users.User{Id: 666, Name: "John", Email: "john@mail.com"}
	
	repositories.UsersRepository = &usersRepoMock{}

	_, err := UsersService.CreateUser(user)

	assert.Equal(t, 400, err.Status())
}

func TestCreateRepositoryFailReturnInternalServerError(t *testing.T) {
	user := users.User{Id: 666, Name: "John", Email: "john@mail.com", Password: "admin"}
	createUserRepoFunc = func(user *users.User) rest_errors.RestErr {
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}

	repositories.UsersRepository = &usersRepoMock{}

	_, err := UsersService.CreateUser(user)

	assert.Equal(t, 500, err.Status())	
}

func TestGetOK(t *testing.T) {

	user := users.User{Id: 666, Name: "John", Email: "john@mail.com", Password: "admin"}
	getUserRepoFunc = func(Id int64) (*users.User, rest_errors.RestErr) {
		return &user, nil
	}

	repositories.UsersRepository = &usersRepoMock{}

	_, err := UsersService.GetUser(666)

	assert.NoError(t, err)
	assert.Equal(t, int64(666), user.Id)
	assert.Equal(t, "John", user.Name)
}