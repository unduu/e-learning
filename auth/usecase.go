package auth

import "github.com/unduu/e-learning/auth/model"

type Usecase interface {
	Login(username string, password string) (user *model.User, tokenString string)
	Register(fullname string, phone string, email string, username string, password string) (verifivationCode string, affected int64)
	Verify(username, code string) (affected int64)
	ForgotPassword(phone string) (affected int64, passKey string)
	ResetPassword(password string, passkey string) (affected int64)
	SendVerificationCode(code string, phone string, body string)
	ResendVerificationCode(username string) bool
}
