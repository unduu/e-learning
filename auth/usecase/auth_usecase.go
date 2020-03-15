package usecase

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/unduu/e-learning/auth"
	"github.com/unduu/e-learning/auth/model"
	"io"
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

// Login user to system
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

// Register new user
func (a *AuthUsecase) Register(fullname string, phone string, email string, username string, password string) (verifivationCode string, affected int64) {
	vcode := EncodeToString(6)
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

	return vcode, a.repository.InsertNewUser(user, vcode)
}

// Register new user
func (a *AuthUsecase) Verify(username, code string) (affected int64) {
	isValid := a.repository.UpdateUserStatus(username, code)
	return isValid
}

// EncodeToString return auto generated 6 digit number
func EncodeToString(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
