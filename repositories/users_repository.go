package repositories

import (	
	"errors"
	"tokenalert_user-api/datasources/mysql/users_db"
	"tokenalert_user-api/domain/users"

	"github.com/rafawilliner/tokenalert_utils-go/logger"
	"github.com/rafawilliner/tokenalert_utils-go/rest_errors"
)

const (
	queryInsertUser = "INSERT INTO users(name, email, telegram_user, status, password, date_created) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser    = "SELECT id, name, email, telegram_user, status, date_created FROM users WHERE id=?;"
)

var (
	UsersRepository userRepositoryInterface = &usersRepository{}
)

type usersRepository struct{}

type userRepositoryInterface interface {
	Save(*users.User) rest_errors.RestErr
	Get(int64) (*users.User, rest_errors.RestErr)
	FindByEmailAndPassword() rest_errors.RestErr
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

func (u *usersRepository) FindByEmailAndPassword() rest_errors.RestErr {

	return nil
}
