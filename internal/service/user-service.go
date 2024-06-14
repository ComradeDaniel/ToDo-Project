package service

import (
	"database/sql"
	"errors"
	"log"
	"todolist/internal/auth"
	"todolist/internal/database"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(database.User) (string, error)
	LoginUser(database.User) (string, error)
}

type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

var (
	ErrNoSuchUser        error = errors.New("there is no user with this username")
	ErrWrongPassword     error = errors.New("wrong password")
	ErrUserAlreadyExists error = errors.New("a user with this username already exists")
)

// Registers the new user. Returns nil on success or ErrUserAlreadyExists if the user already exists
func (service *userService) RegisterUser(user database.User) (string, error) {
	err := database.AddUser(user)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			return "", ErrUserAlreadyExists
		}
		return "", err
	}
	token := auth.GenerateToken(user.Username)
	return token, nil
}

// Returns a jwt, nil if the login was successful
func (service *userService) LoginUser(user database.User) (string, error) {

	dbUser, err := database.GetUserByUsername(user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoSuchUser
		}
		log.Fatal(err)
	}

	result := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if result != nil {
		if errors.Is(result, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrWrongPassword
		}
		log.Fatal(result)
	}

	token := auth.GenerateToken(user.Username)
	return token, nil
}
