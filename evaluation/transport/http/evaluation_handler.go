package http

import (
	"github.com/unduu/e-learning/evaluation"
	customValidator "github.com/unduu/e-learning/helper/validator"
	"github.com/unduu/e-learning/middleware"
	"github.com/unduu/e-learning/response"
	"gopkg.in/go-playground/validator.v9"
	"math"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	//"gopkg.in/go-playground/validator.v9"
)

type EvaluationHandler struct {
	EvaluationUsecase evaluation.Usecase
	Middleware        *middleware.Middleware
	Validator         *customValidator.CustomValidator
}

func NewHttpAuthHandler(router *gin.RouterGroup, mw *middleware.Middleware, v *customValidator.CustomValidator, evaluationUC evaluation.Usecase) {
	handler := &EvaluationHandler{
		EvaluationUsecase: evaluationUC,
		Middleware:        mw,
		Validator:         v,
	}
	router.GET("test/pre", mw.AuthMiddleware, handler.PreEvaluation)
	router.POST("test/pre", mw.AuthMiddleware, handler.ProcessEvaluationAnswer)
	router.GET("test/post", mw.AuthMiddleware, handler.PostEvaluation)
	router.POST("test/post", mw.AuthMiddleware, handler.ProcessPostAnswer)
	router.GET("test/quiz", mw.AuthMiddleware, handler.QuizEvaluation)
	router.POST("test/quiz", mw.AuthMiddleware, handler.ProcessQuizAnswer)
	router.POST("test/pre/reset", mw.AuthMiddleware, handler.ResetPrePostStatus)
	router.POST("test/post/reset", mw.AuthMiddleware, handler.ResetPrePostStatus)
}

// PreEvaluation return pre test question
func (e *EvaluationHandler) PreEvaluation(c *gin.Context) {
	// Form Data
	var req RequestEvaluation
	// Validation
	err := c.ShouldBind(&req)
	if err != nil {
		//a.Middleware.CheckValidate(err, c)
		var errValidation []response.Error
		if reflect.TypeOf(err).String() != "validator.ValidationErrors" {
			error := response.Error{"", err.Error()}
			errValidation = append(errValidation, error)
			response.RespondErrorJSON(c.Writer, errValidation)
			return
		}
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(e.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}

	assesment, totalData := e.EvaluationUsecase.StartEvaluation("prepost", req.Page, req.Limit)

	// Pagination
	totalPage := int(math.Round(float64(totalData) / float64(req.Limit)))
	prevPage := req.Page - 1
	if prevPage <= 0 {
		prevPage = 1
	}
	nextPage := req.Page + 1
	if nextPage >= totalPage {
		nextPage = totalPage
	}

	msg := "List of questions"
	res := PreEvaluationResponse{
		StartTime: assesment.Start,
		EndTime:   assesment.End,
		Pagination: PaginationResponse{
			TotalData:   totalData,
			TotalPage:   totalPage,
			Limit:       req.Limit,
			Current:     req.Page,
			PreviousUrl: "/test/pre?page=" + strconv.Itoa(prevPage) + "&limit=" + strconv.Itoa(req.Limit),
			NextUrl:     "/test/pre?page=" + strconv.Itoa(nextPage) + "&limit=" + strconv.Itoa(req.Limit),
		},
	}

	for _, question := range assesment.QuestionList {
		q := Question{
			Id:         question.Id,
			Type:       question.Type,
			AttachType: question.AttachType,
			Attachment: question.Attachment,
			Question:   question.Text,
			Choices: Choice{
				Type:    question.Choices.Type,
				Options: question.Choices.Options,
			},
		}
		res.Test = append(res.Test, q)
	}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

// PostEvaluation return post test question
func (e *EvaluationHandler) PostEvaluation(c *gin.Context) {
	// Form Data
	var req RequestEvaluation
	// Validation
	err := c.ShouldBind(&req)
	if err != nil {
		//a.Middleware.CheckValidate(err, c)
		var errValidation []response.Error
		if reflect.TypeOf(err).String() != "validator.ValidationErrors" {
			error := response.Error{"", err.Error()}
			errValidation = append(errValidation, error)
			response.RespondErrorJSON(c.Writer, errValidation)
			return
		}
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(e.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}
	// User session
	loggedIn := e.Middleware.GetLoggedInUser(c)

	assesment, totalData := e.EvaluationUsecase.StartPostEvaluation(loggedIn.Username, req.Page, req.Limit)
	if assesment == nil {
		msg := "You have to join pre test first"
		err := make([]string, 0)
		response.RespondErrorJSON(c.Writer, err, msg)
		return
	}

	// Pagination
	totalPage := int(math.Round(float64(totalData) / float64(req.Limit)))
	prevPage := req.Page - 1
	if prevPage <= 0 {
		prevPage = 1
	}
	nextPage := req.Page + 1
	if nextPage >= totalPage {
		nextPage = totalPage
	}

	msg := "List of questions"
	res := PreEvaluationResponse{
		StartTime: assesment.Start,
		EndTime:   assesment.End,
		Pagination: PaginationResponse{
			TotalData:   totalData,
			TotalPage:   totalPage,
			Limit:       req.Limit,
			Current:     req.Page,
			PreviousUrl: "/test/pre?page=" + strconv.Itoa(prevPage) + "&limit=" + strconv.Itoa(req.Limit),
			NextUrl:     "/test/pre?page=" + strconv.Itoa(nextPage) + "&limit=" + strconv.Itoa(req.Limit),
		},
	}

	for _, question := range assesment.QuestionList {
		q := Question{
			Id:         question.Id,
			Type:       question.Type,
			AttachType: question.AttachType,
			Attachment: question.Attachment,
			Question:   question.Text,
			Choices: Choice{
				Type:    question.Choices.Type,
				Options: question.Choices.Options,
			},
		}
		res.Test = append(res.Test, q)
	}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

// QuizEvaluation return quiz test question
func (e *EvaluationHandler) QuizEvaluation(c *gin.Context) {
	// Form Data
	var req RequestEvaluation

	// Validation
	err := c.ShouldBind(&req)
	if err != nil {
		//a.Middleware.CheckValidate(err, c)
		var errValidation []response.Error
		if reflect.TypeOf(err).String() != "validator.ValidationErrors" {
			error := response.Error{"", err.Error()}
			errValidation = append(errValidation, error)
			response.RespondErrorJSON(c.Writer, errValidation)
			return
		}
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(e.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}

	assesment, totalData := e.EvaluationUsecase.StartEvaluation(req.Title, req.Page, req.Limit)

	// Pagination
	totalPage := int(math.Round(float64(totalData) / float64(req.Limit)))
	prevPage := req.Page - 1
	if prevPage <= 0 {
		prevPage = 1
	}
	nextPage := req.Page + 1
	if nextPage >= totalPage {
		nextPage = totalPage
	}

	msg := "List of questions"
	res := PreEvaluationResponse{
		StartTime: assesment.Start,
		EndTime:   assesment.End,
		Pagination: PaginationResponse{
			TotalData:   totalData,
			TotalPage:   totalPage,
			Limit:       req.Limit,
			Current:     req.Page,
			PreviousUrl: "/test/pre?page=" + strconv.Itoa(prevPage) + "&limit=" + strconv.Itoa(req.Limit),
			NextUrl:     "/test/pre?page=" + strconv.Itoa(nextPage) + "&limit=" + strconv.Itoa(req.Limit),
		},
	}

	for _, question := range assesment.QuestionList {
		q := Question{
			Id:         question.Id,
			Type:       question.Type,
			AttachType: question.AttachType,
			Attachment: question.Attachment,
			Question:   question.Text,
			Choices: Choice{
				Type:    question.Choices.Type,
				Options: question.Choices.Options,
			},
		}
		res.Test = append(res.Test, q)
	}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

// ProcessEvaluationAnswer receive answer from user
func (e *EvaluationHandler) ProcessEvaluationAnswer(c *gin.Context) {
	// Form Data
	var req RequestProcessEvaluationAnswer
	// Validation
	err := c.ShouldBind(&req)
	if err != nil {
		//a.Middleware.CheckValidate(err, c)
		var errValidation []response.Error
		if reflect.TypeOf(err).String() != "validator.ValidationErrors" {
			error := response.Error{"", err.Error()}
			errValidation = append(errValidation, error)
			response.RespondErrorJSON(c.Writer, errValidation)
			return
		}
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(e.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}
	// User Logged in session
	loggedIn := e.Middleware.GetLoggedInUser(c)

	// Check if user not join pre test yet
	exists := e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, "pretest")
	if exists {
		msg := "You already join pre test"
		err := make([]string, 0)
		response.RespondErrorJSON(c.Writer, err, msg)
		return
	}

	e.EvaluationUsecase.CheckAnswerResult(req.Answer)
	e.EvaluationUsecase.SaveAnswer(loggedIn.Username, "pretest", req.Answer)

	msg := "Thank you, We have recieve your answer"
	res := make([]string, 0)
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// ProcessPostAnswer receive post test answer from user
func (e *EvaluationHandler) ProcessPostAnswer(c *gin.Context) {
	// Form Data
	var req RequestProcessEvaluationAnswer
	// Validation
	err := c.ShouldBind(&req)
	if err != nil {
		//a.Middleware.CheckValidate(err, c)
		var errValidation []response.Error
		if reflect.TypeOf(err).String() != "validator.ValidationErrors" {
			error := response.Error{"", err.Error()}
			errValidation = append(errValidation, error)
			response.RespondErrorJSON(c.Writer, errValidation)
			return
		}
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(e.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}
	// User Logged in session
	loggedIn := e.Middleware.GetLoggedInUser(c)

	// Check if user not join pre test yet
	exists := e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, "pretest")
	if !exists {
		msg := "You have to join pre test first"
		err := make([]string, 0)
		response.RespondErrorJSON(c.Writer, err, msg)
		return
	}
	// Check if user already join post test
	exists = e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, "posttest")
	if exists {
		msg := "You already join post test"
		err := make([]string, 0)
		response.RespondErrorJSON(c.Writer, err, msg)
		return
	}

	e.EvaluationUsecase.CheckAnswerResult(req.Answer)
	e.EvaluationUsecase.SaveAnswer(loggedIn.Username, "posttest", req.Answer)

	msg := "Thank you, We have recieve your answer"
	res := make([]string, 0)
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// ProcessQuizAnswer receive quiz answer from user
func (e *EvaluationHandler) ProcessQuizAnswer(c *gin.Context) {
	// Form Data
	var req RequestProcessQuizAnswer
	// Validation
	err := c.ShouldBind(&req)
	if err != nil {
		//a.Middleware.CheckValidate(err, c)
		var errValidation []response.Error
		if reflect.TypeOf(err).String() != "validator.ValidationErrors" {
			error := response.Error{"", err.Error()}
			errValidation = append(errValidation, error)
			response.RespondErrorJSON(c.Writer, errValidation)
			return
		}
		for _, fieldErr := range err.(validator.ValidationErrors) {
			e := fieldErr.Translate(e.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}
	// User Logged in session
	loggedIn := e.Middleware.GetLoggedInUser(c)

	e.EvaluationUsecase.CheckAnswerResult(req.Answer)
	e.EvaluationUsecase.SaveAnswer(loggedIn.Username, req.Title, req.Answer)

	msg := "Thank you, We have recieve your answer"
	res := make([]string, 0)
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// ResetPrePostStatus set current pre post answer status to archived
func (e *EvaluationHandler) ResetPrePostStatus(c *gin.Context) {
	// User session
	loggedIn := e.Middleware.GetLoggedInUser(c)
	e.EvaluationUsecase.ArchivedPrePostAnswer(loggedIn.Username)

	// Response
	msg := "Your pre post status has been reset"
	res := make([]string, 0)
	response.RespondSuccessJSON(c.Writer, res, msg)
}
