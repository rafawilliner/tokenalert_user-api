package repositories

import (
	"errors"
	"strings"
	"tokenalert_user-api/src/datasources/mysql/users_db"
	"tokenalert_user-api/src/domain/users"
	"tokenalert_user-api/src/utils/mysql_utils"
	"github.com/rafawilliner/tokenalert_utils-go/src/logger"
	"github.com/rafawilliner/tokenalert_utils-go/src/rest_errors"
)

const (
	queryInsertUser             = "INSERT INTO users(name, email, telegram_user, status, password, date_created) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, name, email, telegram_user, status, date_created FROM users WHERE id=?;"
	queryFindByEmailAndPassword = "SELECT id, name, email, telegram_user, status FROM users WHERE email=? AND password=? AND status=?"
)

var (
	UsersRepository userRepositoryInterface = &usersRepository{}
)

type usersRepository struct{}

type userRepositoryInterface interface {
	Save(*users.User) rest_errors.RestErr
	Get(int64) (*users.User, rest_errors.RestErr)
	FindByEmailAndPassword(users.LoginRequest) (*users.User, rest_errors.RestErr)
}

func (u *usersRepository) Save(user *users.User) rest_errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error saving user", errors.New("database error"))
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(&user.Name, user.Email, user.TelegramUser, user.Status, user.Password, user.DateCreated)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return rest_errors.NewInternalServerError("error saving user", errors.New("database error"))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("error saving user", errors.New("database error"))
	}
	user.Id = userId
	return nil
}

func (u *usersRepository) Get(id int64) (*users.User, rest_errors.RestErr) {

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return nil, rest_errors.NewInternalServerError("error fetching user", errors.New("database error"))
	}
	defer stmt.Close()

	result := stmt.QueryRow(id)

	var user users.User
	if getErr := result.Scan(&user.Id, &user.Name, &user.Email, &user.TelegramUser, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error when trying to get user by id", getErr)
		return nil, rest_errors.NewInternalServerError("error fetching user", errors.New("database error"))
	}
	return &user, nil
}

func (u *usersRepository) FindByEmailAndPassword(login users.LoginRequest) (*users.User, rest_errors.RestErr) {

	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return nil,  rest_errors.NewInternalServerError("error when tying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	var user users.User
	result := stmt.QueryRow(login.Email, login.Password, users.StatusActive)
	if getErr := result.Scan(&user.Id, &user.Name, &user.Email, &user.TelegramUser, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return nil, rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return nil, rest_errors.NewInternalServerError("error when trying to find user", errors.New("database error"))
	}

	return &user, nil
}
