package usecase

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/e-learning/auth"
	"github.com/e-learning/auth/model"
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

func (a *AuthUsecase) Login(username string, password string) string {
	tokenString := ""
	user, err := a.repository.GetByUsernamePassword(username, password)

	if err == nil {
		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(5 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &model.Claims{
			Username: user.Username,
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

	return tokenString
}
