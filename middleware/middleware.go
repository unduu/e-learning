package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/unduu/e-learning/auth/model"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"strings"

	_customValidator "github.com/unduu/e-learning/helper/validator"
	"github.com/unduu/e-learning/response"
)

type UserSession struct {
	Username string
}

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) CheckValidate(err error, c *gin.Context) bool {
	var cv = _customValidator.NewCustomValidator()
	var errValidation []response.Error
	for _, fieldErr := range err.(validator.ValidationErrors) {
		e := fieldErr.Translate(cv.Translation)

		error := response.Error{fieldErr.Field(), e}
		errValidation = append(errValidation, error)
	}
	response.RespondErrorJSON(c.Writer, errValidation)
	return false
}

func (m *Middleware) AuthMiddleware(c *gin.Context) {
	// Grab the token from the header
	tokenHeader := c.GetHeader("Authorization")

	// Token is missing, returns with error code 403 Unauthorized
	if tokenHeader == "" {
		response.RespondUnauthorizedJSON(c.Writer, "Missing auth token")
		c.Abort()
		return
	}

	// The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		response.RespondUnauthorizedJSON(c.Writer, "Invalid/Malformed auth token")
		c.Abort()
		return
	}

	// Grab the token part, what we are truly interested in
	session := UserSession{}
	tokenPart := splitted[1]
	tk := &model.Claims{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		// User logged in session
		session = UserSession{
			Username: tk.Username,
		}
		m.SetLoggedInUserInfo(session, c)

		return []byte(os.Getenv("token_password")), nil
	})

	// Malformed token, returns with http code 403 as usual
	if err != nil {
		response.RespondUnauthorizedJSON(c.Writer, "Malformed authentication token")
		c.Abort()
		return
	}

	// Token is invalid, maybe not signed on this server
	if !token.Valid {
		response.RespondUnauthorizedJSON(c.Writer, "Token is not valid.")
		c.Abort()
		return
	}
	// Call the next handler, which can be another middleware in the chain, or the final handler.
	c.Next()
}

func (m *Middleware) SetLoggedInUserInfo(userinfo UserSession, c *gin.Context) {
	c.Set("userinfo", userinfo)
}

func (m *Middleware) GetLoggedInUser(c *gin.Context) UserSession {
	userInfo := c.MustGet("userinfo").(UserSession)
	return userInfo
}
