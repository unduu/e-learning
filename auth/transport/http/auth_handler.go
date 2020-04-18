package http

import (
	"github.com/dongri/phonenumber"
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
	router.POST("register", handler.Register)
	router.POST("register/verification", mw.TokenCheckMiddleware, handler.Verify)
	router.POST("register/verification/resend", mw.TokenCheckMiddleware, handler.ResendVerifCode)
	router.POST("password/forgot", handler.ForgotPassword)
	router.POST("password/reset", handler.ResetPassword)
}

// Login return auth token
func (a *AuthHandler) Login(c *gin.Context) {
	// Request
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

	// Processing
	user, token := a.AuthUsecase.Login(req.Username, req.Password)

	// Response
	if token == "" {
		errResponse := make([]string, 0)
		response.RespondErrorJSON(c.Writer, errResponse, "Incorrect username or password")
		return
	}
	msg := "Welcome " + req.Username
	res := LoginResponse{
		User:  User{req.Username, "menthor", user.Status, user.StatusCode},
		Token: token,
	}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// Register new user
func (a *AuthHandler) Register(c *gin.Context) {
	// Request validation
	var req RequestRegister
	err := c.ShouldBind(&req)
	if err != nil {
		var errValidation []response.Error
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(a.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}

	// Processing
	phoneNormalize := phonenumber.Parse(req.Phone, "ID")
	verificationCode, _ := a.AuthUsecase.Register(req.Fullname, phoneNormalize, req.Email, req.Username, req.Password)
	user, token := a.AuthUsecase.Login(req.Username, req.Password)

	// Respponse
	msg := "We have sent a verification code to your email address"
	res := LoginResponseTemp{
		User:       User{req.Username, "menthor", user.Status, user.StatusCode},
		Token:      token,
		Activation: verificationCode,
	}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// Verify to activate user account
func (a *AuthHandler) Verify(c *gin.Context) {
	// Request validation
	var req RequestVerify
	err := c.ShouldBind(&req)
	if err != nil {
		var errValidation []response.Error
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(a.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}

	// Processing
	loggedIn := a.Middleware.GetLoggedInUser(c)
	success := a.AuthUsecase.Verify(loggedIn.Username, req.Code)

	// Response
	if success <= 0 {
		msg := "Incorrect code"
		err := response.Error{"code", "Enter a valid activation code"}
		response.RespondErrorJSON(c.Writer, err, msg)
		return
	}
	msg := "Your account is activated"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// Forgot to reset user password
func (a *AuthHandler) ForgotPassword(c *gin.Context) {
	// Request validation
	var req RequestForgotPassword
	err := c.ShouldBind(&req)
	if err != nil {
		var errValidation []response.Error
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(a.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}

	// Processing
	phoneNormalize := phonenumber.Parse(req.Phone, "ID")
	affected, code := a.AuthUsecase.ForgotPassword(phoneNormalize)
	if affected > 0 {
		body := "You requested to reset your password, please enter this code " + code
		a.AuthUsecase.SendVerificationCode(code, phoneNormalize, body)
	}

	// Response
	msg := "We have sent a confirmation code to reset your password"
	//res := make([]string, 0)
	res := ForgotPasswordResponseTemp{ConfirmationCode: code}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// Forgot to reset user password
func (a *AuthHandler) ResetPassword(c *gin.Context) {
	// Request validation
	var req RequestResetPassword
	err := c.ShouldBind(&req)
	if err != nil {
		var errValidation []response.Error
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(a.Validator.Translation)
			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}

	// Processing
	a.AuthUsecase.ResetPassword(req.PasswordNew, req.PasswordKey)

	// Response
	msg := "Your password has been reset successfully"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// ResendVerifCode resend verification code to user phone
func (a *AuthHandler) ResendVerifCode(c *gin.Context) {
	// Session
	loggedIn := a.Middleware.GetLoggedInUser(c)

	// Processing
	ok := a.AuthUsecase.ResendVerificationCode(loggedIn.Username)

	// Response Error
	if !ok {
		errResponse := make([]string, 0)
		response.RespondErrorJSON(c.Writer, errResponse, "You have entered wrong phone number")
		return
	}

	// Response Success
	msg := "We have sent a verification code to your phone number"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// Logout remove user access token
func (a *AuthHandler) Logout(c *gin.Context) {
	msg := "You have successfully logged out"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}
