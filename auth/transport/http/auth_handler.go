package http

import (
	"github.com/gin-gonic/gin"
	"github.com/unduu/e-learning/auth"
	customValidator "github.com/unduu/e-learning/helper/validator"
	"github.com/unduu/e-learning/middleware"
	"github.com/unduu/e-learning/response"
	"gopkg.in/go-playground/validator.v9"
)

type AuthHandler struct {
	AuthUsecase auth.Usecase
	Middleware  *middleware.Middleware
	Validator   *customValidator.CustomValidator
}

func NewHttpAuthHandler(router *gin.RouterGroup, mw *middleware.Middleware, v *customValidator.CustomValidator, authUC auth.Usecase) {
	handler := &AuthHandler{
		AuthUsecase: authUC,
		Middleware:  mw,
		Validator:   v,
	}
	router.POST("login", handler.Login)
	router.GET("logout", mw.AuthMiddleware, handler.Logout)
}

// Login return auth token
func (a *AuthHandler) Login(c *gin.Context) {
	var req RequestLogin

	err := c.ShouldBind(&req)
	if err != nil {
		//a.Middleware.CheckValidate(err, c)
		var errValidation []response.Error
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(a.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}

	token := a.AuthUsecase.Login(req.Username, req.Password)

	if token == "" {
		errResponse := make([]string, 0)
		response.RespondErrorJSON(c.Writer, errResponse, "Incorrect username or password")
		return
	}

	msg := "Welcome " + req.Username
	res := LoginResponse{
		User:  User{req.Username, "menthor"},
		Token: token,
	}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

// Logout remove user access token
func (a *AuthHandler) Logout(c *gin.Context) {
	msg := "You have successfully logged out"
	res := make([]string, 0)
	response.RespondSuccessJSON(c.Writer, res, msg)
}
