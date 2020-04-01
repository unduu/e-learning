package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/unduu/e-learning/auth/model"
	_customValidator "github.com/unduu/e-learning/helper/validator"
	"github.com/unduu/e-learning/response"
	"gopkg.in/go-playground/validator.v9"
	"os"
	"strings"
)

type UserSession struct {
	Username   string
	StatusCode int
}

type Middleware struct {
	conn *sqlx.DB
}

func NewMiddleware(db *sqlx.DB) *Middleware {
	return &Middleware{conn: db}
}

func (m *Middleware) CheckValidate(err error, c *gin.Context) bool {
	var cv = _customValidator.NewCustomValidator(m.conn)
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
	session := m.ValidateToken(c)
	m.SetLoggedInUserInfo(session, c)

	if !m.IsActivated(session) {
		response.RespondUnverifyJSON(c.Writer, "Please verify your account to access this page")
		c.Abort()
		return
	}
}

func (m *Middleware) TokenCheckMiddleware(c *gin.Context) {
	session := m.ValidateToken(c)
	m.SetLoggedInUserInfo(session, c)
	// Call the next handler, which can be another middleware in the chain, or the final handler.
	c.Next()
}

func (m *Middleware) ValidateToken(c *gin.Context) UserSession {
	session := UserSession{}
	// Grab the token from the header
	tokenHeader := c.GetHeader("Authorization")

	// Token is missing, returns with error code 403 Unauthorized
	if tokenHeader == "" {
		response.RespondUnauthorizedJSON(c.Writer, "Missing auth token")
		c.Abort()
		return session
	}

	// The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		response.RespondUnauthorizedJSON(c.Writer, "Invalid/Malformed auth token")
		c.Abort()
		return session
	}

	// Grab the token part, what we are truly interested in
	tokenPart := splitted[1]
	tk := &model.Claims{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		// User logged in session
		session = UserSession{
			Username:   tk.Username,
			StatusCode: tk.StatusCode,
		}
		return []byte(os.Getenv("token_password")), nil
	})

	// Malformed token, returns with http code 403 as usual
	if err != nil {
		response.RespondUnauthorizedJSON(c.Writer, "Malformed authentication token")
		c.Abort()
		return session
	}

	// Token is invalid, maybe not signed on this server
	if !token.Valid {
		response.RespondUnauthorizedJSON(c.Writer, "Token is not valid.")
		c.Abort()
		return session
	}

	return session
}

func (m *Middleware) SetLoggedInUserInfo(userinfo UserSession, c *gin.Context) {
	c.Set("userinfo", userinfo)
}

func (m *Middleware) GetLoggedInUser(c *gin.Context) UserSession {
	userInfo := c.MustGet("userinfo").(UserSession)
	return userInfo
}

func (m *Middleware) IsActivated(userinfo UserSession) bool {
	if userinfo.StatusCode == 0 {
		return false
	}
	return true
}
