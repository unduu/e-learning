package auth

import "github.com/unduu/e-learning/auth/model"

type Usecase interface {
	Login(username string, password string) (user *model.User, tokenString string)
	Register(fullname string, phone string, email string, username string, password string) (affected int64)
}
