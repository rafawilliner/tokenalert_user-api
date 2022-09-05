package repositories

import (
	"database/sql"
	"errors"
	"testing"
	"tokenalert_user-api/src/datasources/mysql/users_db"
	"tokenalert_user-api/src/domain/users"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rafawilliner/tokenalert_utils-go/src/logger"
	"github.com/rafawilliner/tokenalert_utils-go/src/rest_errors"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.Error("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestSaveOK(t *testing.T) {

	db, mock := NewMock()
	users_db.Client = db	
	defer func() {
		users_db.Client.Close()
	}()

	user := users.User{Name: "John", Email: "john@mail.com", TelegramUser: "@john", Password: "admin", DateCreated: "2022-01-01"}

	query := "INSERT INTO users(name, email, telegram_user, status, password, date_created) VALUES(?, ?, ?, ?, ?, ?);"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(user.Name, user.Email, user.TelegramUser, user.Status, user.Password, user.DateCreated).WillReturnResult(sqlmock.NewResult(667, 1))

	err := UsersRepository.Save(&user)
	
	assert.NoError(t, err)
	assert.Equal(t, int64(667), user.Id)
}

func TestSavePrepareQueryFailed(t *testing.T) {

	db, mock := NewMock()
	users_db.Client = db	
	defer func() {
		users_db.Client.Close()
	}()

	user := users.User{Name: "John", Email: "john@mail.com", TelegramUser: "@john", Password: "admin", DateCreated: "2022-01-01"}

	query := "INSERT INTO users(name, email, telegram_user, status, password, date_created) VALUES(?, ?, ?, ?, ?);"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(user.Name, user.Email, user.TelegramUser, user.Status, user.Password, user.DateCreated).WillReturnResult(sqlmock.NewResult(667, 1))

	err := UsersRepository.Save(&user)
	
	assert.Error(t, err)
	assert.Equal(t, 500, err.Status())	
}

func TestSaveExecutionFailed(t *testing.T) {

	db, mock := NewMock()
	users_db.Client = db	
	defer func() {
		users_db.Client.Close()
	}()

	user := users.User{Name: "John", Email: "john@mail.com", TelegramUser: "@john", Password: "admin", DateCreated: "2022-01-01"}

	query := "INSERT INTO users(name, email, telegram_user, status, password, date_created) VALUES(?, ?, ?, ?, ?, ?);"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(user.Name, user.Email, user.TelegramUser, user.Status, user.Password, user.DateCreated).WillReturnError(rest_errors.NewInternalServerError("internal_server_error", errors.New("database error")))

	err := UsersRepository.Save(&user)
	
	assert.Error(t, err)
	assert.Equal(t, 500, err.Status())	
	assert.Equal(t, "error saving user", err.Message())	
}

func TestGetOK(t *testing.T) {

	db, mock := NewMock()
	users_db.Client = db	
	defer func() {
		users_db.Client.Close()
	}()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "telegram_user", "date_created", "status"}).
		AddRow(667, "john", "john@mail.com", "@john", "2022-01-01", "active")		

	query := "SELECT id, name, email, telegram_user, status, date_created FROM users WHERE id=?;"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(667).WillReturnRows(rows)

	user, err := UsersRepository.Get(667)
	
	assert.NoError(t, err)
	assert.Equal(t, int64(667), user.Id)
	assert.Equal(t, "@john", user.TelegramUser)
}

func TestGetPrepareQueryFailed(t *testing.T) {

	db, mock := NewMock()
	users_db.Client = db	
	defer func() {
		users_db.Client.Close()
	}()

	query := "SELECT id, name, email, telegram_user, status, date_created FROM users WHERE id=?;"	
	expected := mock.ExpectPrepare(query).WillReturnError(rest_errors.NewInternalServerError("internal_server_error_prepare", errors.New("database error")))
	
	_, err := UsersRepository.Get(667)
	
	assert.NotNil(t, err)
	assert.NotNil(t, expected)
	assert.Equal(t, 500, err.Status())	
	assert.Equal(t, "error fetching user", err.Message())	
}

func TestGetExecutionFailed(t *testing.T) {

	db, mock := NewMock()
	users_db.Client = db	
	defer func() {
		users_db.Client.Close()
	}()

	query := "SELECT id, name, email, telegram_user, status, date_created FROM users WHERE id=?;"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(667).WillReturnError(rest_errors.NewInternalServerError("internal_server_error", errors.New("database error")))

	_, err := UsersRepository.Get(667)
	
	assert.Error(t, err)
	assert.Equal(t, 500, err.Status())	
	assert.Equal(t, "error fetching user", err.Message())	
}

func TestFindByEmailAndPasswordOK(t *testing.T) {

	db, mock := NewMock()
	users_db.Client = db	
	defer func() {
		users_db.Client.Close()
	}()

	loginRequest := users.LoginRequest{Email: "john@mail.com", Password: "ABC123"}
	rows := sqlmock.NewRows([]string{"id", "name", "email", "telegram_user", "status"}).
		AddRow(667, "john", "john@mail.com", "@john", "active")		

	query := "SELECT id, name, email, telegram_user, status FROM users WHERE email=? AND password=? AND status=?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(loginRequest.Email, loginRequest.Password, users.StatusActive).WillReturnRows(rows)

	user, err := UsersRepository.FindByEmailAndPassword(loginRequest)
	
	assert.NoError(t, err)
	assert.Equal(t, int64(667), user.Id)
	assert.Equal(t, "john@mail.com", user.Email)
	assert.Equal(t, "@john", user.TelegramUser)
}

func TestFindByEmailAndPasswordPrepareQueryFailed(t *testing.T) {

	db, mock := NewMock()
	users_db.Client = db	
	defer func() {
		users_db.Client.Close()
	}()

	loginRequest := users.LoginRequest{Email: "john@mail.com", Password: "ABC123"}
	query := "SELECT id, name, email, telegram_user, status FROM users WHERE email=? AND password=? AND status=?"
	expected := mock.ExpectPrepare(query).WillReturnError(rest_errors.NewInternalServerError("internal_server_error_prepare", errors.New("database error")))
	
	_, err := UsersRepository.FindByEmailAndPassword(loginRequest)
	
	assert.NotNil(t, err)
	assert.NotNil(t, expected)
	assert.Equal(t, 500, err.Status())	
	assert.Equal(t, "error when tying to find user", err.Message())	
}

func TestFindByEmailAndPasswordExecutionFailed(t *testing.T) {

	db, mock := NewMock()
	users_db.Client = db	
	defer func() {
		users_db.Client.Close()
	}()

	loginRequest := users.LoginRequest{Email: "john@mail.com", Password: "ABC123"}
	query := "SELECT id, name, email, telegram_user, status FROM users WHERE email=? AND password=? AND status=?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(667).WillReturnError(rest_errors.NewInternalServerError("internal_server_error", errors.New("database error")))

	_, err := UsersRepository.FindByEmailAndPassword(loginRequest)
	
	assert.Error(t, err)
	assert.Equal(t, 500, err.Status())	
	assert.Equal(t, "error when trying to find user", err.Message())	
}