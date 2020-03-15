package auth

import "github.com/unduu/e-learning/auth/model"

type Repository interface {
	GetByUsernamePassword(username string, password string) (*model.User, error)
	InsertNewUser(user model.User, verifCode string) (affected int64)
	UpdateUserStatus(username string, code string) (affected int64)
}
