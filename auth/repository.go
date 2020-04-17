package auth

import "github.com/unduu/e-learning/auth/model"

type Repository interface {
	GetByUsername(username string) (*model.User, error)
	GetByUsernamePassword(username string, password string) (*model.User, error)
	GetByPhone(phone string) (*model.User, error)
	InsertNewUser(user model.User, verifCode string) (affected int64)
	UpdateUserStatus(username string, code string) (affected int64)
	InsertPasswordKey(phone string, passkey string) (affected int64)
	UpdateNewPassword(password string, passkey string) (affected int64)
}
