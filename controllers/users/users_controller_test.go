package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tokenalert_user-api/domain/users"
	"tokenalert_user-api/services"

	"github.com/gin-gonic/gin"
	"github.com/rafawilliner/tokenalert_utils-go/rest_errors"
	"github.com/stretchr/testify/assert"
)

var (
	createUserFunc func(user users.User) (*users.User, rest_errors.RestErr)
	getUserFunc func(id int64) (*users.User, rest_errors.RestErr)
	loginUserFunc  func(request users.LoginRequest) (*users.User, rest_errors.RestErr)
)

type usersServiceMock struct{}

func (*usersServiceMock) CreateUser(user users.User) (*users.User, rest_errors.RestErr) {
	return createUserFunc(user)
}

func (*usersServiceMock) GetUser(id int64) (*users.User, rest_errors.RestErr) {
	return getUserFunc(id)
}

func (*usersServiceMock) LoginUser(request users.LoginRequest) (*users.User, rest_errors.RestErr) {
	return loginUserFunc(request)
}

func TestUserCreateOK(t *testing.T) {

	createUserFunc = func(user users.User) (*users.User, rest_errors.RestErr) {
		return &users.User{Id: 123, Name: "Serge", Email: "serge@gmail.com", TelegramUser: "@serge"}, nil
	}

	services.UsersService = &usersServiceMock{}

	bodyUser := users.User{
		Id: 123,
		Name:  "Serge",
		Email: "serge@gmail.com",
		TelegramUser: "@serge",
	}

	body, _ := json.Marshal(bodyUser)

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))

	Create(c)

	var userResponse users.User
	error := json.Unmarshal(response.Body.Bytes(), &userResponse)

	assert.Nil(t, error)
	assert.EqualValues(t, http.StatusCreated, response.Code)
	assert.EqualValues(t, 123, userResponse.Id)	
}

func TestUserCreateBadRequestError(t *testing.T) {

	bodyUser := map[string]interface{}{
		"Id": "123ABC",
		"Name": "John",
		"Email": "email@mail.com",
	}

	body, _ := json.Marshal(bodyUser)

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))

	Create(c)

	var userResponse users.User
	error := json.Unmarshal(response.Body.Bytes(), &userResponse)

	assert.NotNil(t, error)
	assert.EqualValues(t, http.StatusBadRequest, response.Code)
}

func TestUserCreateInternalError(t *testing.T) {

	createUserFunc = func(user users.User) (*users.User, rest_errors.RestErr) {
		return nil, rest_errors.NewInternalServerError("internal error creating user", nil)
	}

	services.UsersService = &usersServiceMock{}

	bodyUser := users.User{
		Id: 123,
		Name:  "Serge",
		Email: "serge@gmail.com",
		TelegramUser: "@serge",
	}

	body, _ := json.Marshal(bodyUser)

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))

	Create(c)

	var userResponse users.User
	error := json.Unmarshal(response.Body.Bytes(), &userResponse)

	assert.NotNil(t, error)
	assert.EqualValues(t, http.StatusInternalServerError, response.Code)
}

func TestUserGetOK(t *testing.T) {

	getUserFunc = func(int64) (*users.User, rest_errors.RestErr) {
		return &users.User{Id: 123, Name: "Serge", Email: "serge@gmail.com", TelegramUser: "@serge"}, nil
	}

	services.UsersService = &usersServiceMock{}	
	
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodGet, "/users/123", nil)
	c.Params = gin.Params{
		{Key: "user_id", Value: "123"},
	}

	Get(c)

	var userResponse users.User
	error := json.Unmarshal(response.Body.Bytes(), &userResponse)

	assert.Nil(t, error)
	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, 123, userResponse.Id)	
	assert.EqualValues(t, "Serge", userResponse.Name)	
	assert.EqualValues(t, "serge@gmail.com", userResponse.Email)	
}


func TestUserGetBadRequestError(t *testing.T) {

	getUserFunc = func(id int64) (*users.User, rest_errors.RestErr) {
		return nil, rest_errors.NewBadRequestError("wrong parameter format")
	}

	services.UsersService = &usersServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodGet, "/users", nil)
	c.Params = gin.Params{
		{Key: "user_id", Value: "ABC"},
	}

	Get(c)

	var userResponse users.User
	error := json.Unmarshal(response.Body.Bytes(), &userResponse)

	assert.NotNil(t, error)
	assert.EqualValues(t, http.StatusBadRequest, response.Code)
}

func TestUserGetInternalServerError(t *testing.T) {

	getUserFunc = func(id int64) (*users.User, rest_errors.RestErr) {
		return nil, rest_errors.NewInternalServerError("internal error", nil)
	}

	services.UsersService = &usersServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodGet, "/users", nil)
	c.Params = gin.Params{
		{Key: "user_id", Value: "123"},
	}

	Get(c)

	var userResponse users.User
	error := json.Unmarshal(response.Body.Bytes(), &userResponse)

	assert.NotNil(t, error)
	assert.EqualValues(t, http.StatusInternalServerError, response.Code)
}

func TestUserLoginOK(t *testing.T) {

	loginUserFunc = func(request users.LoginRequest) (*users.User, rest_errors.RestErr) {
		return &users.User{Id: 123, Name: "Serge", Email: "email@email.com", TelegramUser: "@serge"}, nil
	}

	services.UsersService = &usersServiceMock{}

	bodyLogin := users.LoginRequest{
		Email: "email@email.com",
		Password: "admin",
	}

	body, _ := json.Marshal(bodyLogin)

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))

	Login(c)

	var userResponse users.User
	error := json.Unmarshal(response.Body.Bytes(), &userResponse)

	assert.Nil(t, error)
	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, 123, userResponse.Id)	
}

func TestUserLoginBadRequestError(t *testing.T) {

	bodyLogin := map[string]interface{}{		
	}

	body, _ := json.Marshal(bodyLogin)

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))

	Login(c)
	
	assert.EqualValues(t, http.StatusBadRequest, response.Code)
}

func TestUserLoginInternalError(t *testing.T) {

	loginUserFunc = func(request users.LoginRequest) (*users.User, rest_errors.RestErr) {
		return nil, rest_errors.NewInternalServerError("internal error login user", nil)
	}
	
	services.UsersService = &usersServiceMock{}

	bodyLogin := users.LoginRequest{
		Email: "email@email.com",
		Password: "admin",
	}

	body, _ := json.Marshal(bodyLogin)

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))

	Login(c)

	var loginResponse users.User
	error := json.Unmarshal(response.Body.Bytes(), &loginResponse)

	assert.NotNil(t, error)
	assert.EqualValues(t, http.StatusInternalServerError, response.Code)
}