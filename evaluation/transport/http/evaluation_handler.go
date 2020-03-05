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
}

func (e *EvaluationHandler) PreEvaluation(c *gin.Context) {
	// Form Data
	var req RequestPreEvaluation
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

	assesment, totalData := e.EvaluationUsecase.StartEvaluation(req.Page, req.Limit)
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

	e.EvaluationUsecase.CheckAnswerResult(req.Answer)
	e.EvaluationUsecase.SaveAnswer(loggedIn.Username, "pretest", req.Answer)

	msg := "Thank you, We have recieve your answer"
	res := make([]string, 0)
	response.RespondSuccessJSON(c.Writer, res, msg)
}
