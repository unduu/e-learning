package auth

import "github.com/e-learning/auth/model"

type Repository interface {
	GetByUsernamePassword(username string, password string) (*model.User, error)
}
