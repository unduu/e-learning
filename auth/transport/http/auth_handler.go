package http

import (
	"github.com/e-learning/auth"
	customValidator "github.com/e-learning/helper/validator"
	"github.com/e-learning/middleware"
	"github.com/e-learning/response"
	"github.com/gin-gonic/gin"
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
}

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

	msg := "Welcome Johndoe"
	res := LoginResponse{
		User:  User{req.Username, "menthor"},
		Token: token,
	}

	response.RespondSuccessJSON(c.Writer, res, msg)
}
