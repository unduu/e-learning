package usecase

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/unduu/e-learning/auth"
	"github.com/unduu/e-learning/auth/model"
	"os"
	"time"
)

type AuthUsecase struct {
	repository auth.Repository
}

func NewAuthUsecase(repository auth.Repository) *AuthUsecase {
	return &AuthUsecase{
		repository: repository,
	}
}

func (a *AuthUsecase) Login(username string, password string) (user *model.User, tokenString string) {
	passMD5 := md5.Sum([]byte(password))
	passMD5String := fmt.Sprintf("%x", passMD5)
	user, err := a.repository.GetByUsernamePassword(username, passMD5String)

	if err == nil {
		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(120 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &model.Claims{
			Username:   user.Username,
			StatusCode: user.StatusCode,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}
		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, _ = token.SignedString([]byte(os.Getenv("token_password")))
	}

	return user, tokenString
}

func (a *AuthUsecase) Register(fullname string, phone string, email string, username string, password string) (affected int64) {
	passMD5 := md5.Sum([]byte(password))
	user := model.User{
		Username:   username,
		Password:   fmt.Sprintf("%x", passMD5),
		Fullname:   fullname,
		Phone:      phone,
		Email:      email,
		Status:     "inactive",
		StatusCode: 0,
		Role:       "menthor",
	}

	return a.repository.InsertNewUser(user)
}
