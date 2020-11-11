package http

import (
	"fmt"
	"github.com/unduu/e-learning/evaluation"
	"github.com/unduu/e-learning/evaluation/model"
	customValidator "github.com/unduu/e-learning/helper/validator"
	"github.com/unduu/e-learning/learning"
	"github.com/unduu/e-learning/middleware"
	"github.com/unduu/e-learning/response"
	"gopkg.in/go-playground/validator.v9"
	"math"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EvaluationHandler struct {
	EvaluationUsecase evaluation.Usecase
	LearningUsecase   learning.Usecase
	Middleware        *middleware.Middleware
	Validator         *customValidator.CustomValidator
}

func NewHttpAuthHandler(router *gin.RouterGroup, mw *middleware.Middleware, v *customValidator.CustomValidator, evaluationUC evaluation.Usecase, learningUC learning.Usecase) {
	handler := &EvaluationHandler{
		EvaluationUsecase: evaluationUC,
		LearningUsecase:   learningUC,
		Middleware:        mw,
		Validator:         v,
	}
	router.GET("test/pre", mw.AuthMiddleware, handler.PreEvaluation)
	router.POST("test/pre", mw.AuthMiddleware, handler.ProcessEvaluationAnswer)
	router.GET("test/post", mw.AuthMiddleware, handler.PostEvaluation)
	router.GET("test/post/result", mw.AuthMiddleware, handler.PostTestResult)
	router.POST("test/post", mw.AuthMiddleware, handler.ProcessPostAnswer)
	router.GET("test/quiz", mw.AuthMiddleware, handler.QuizEvaluation)
	router.POST("test/quiz", mw.AuthMiddleware, handler.ProcessQuizAnswer)
	router.POST("test/pre/reset", mw.AuthMiddleware, handler.ResetPrePostStatus)
	router.POST("test/post/reset", mw.AuthMiddleware, handler.ResetPostStatus)
	router.PUT("test/question/:id", mw.AuthMiddleware, handler.EditQuestion)
	router.DELETE("test/question/:id", mw.AuthMiddleware, handler.DeleteQuestion)
	router.GET("test/question", mw.AuthMiddleware, handler.ListOfQuestion)
	router.GET("test/questions/groups", mw.AuthMiddleware, handler.ListOfGroupsQuestion)
	router.POST("test/groups/:name/questions", mw.AuthMiddleware, handler.AddQuestion)
	router.GET("test/groups/:name/questions", mw.AuthMiddleware, handler.ListOfQuestionsInGroup)
	router.DELETE("test/questions/groups/:name", mw.AuthMiddleware, handler.DeleteGroupQuestion)
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
	// User Logged in session
	loggedIn := e.Middleware.GetLoggedInUser(c)

	// Check if user not join pre test yet
	exists, _ := e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, "pretest")
	if exists {
		msg := "Anda telah mengikuti pre test"
		err := struct{}{}
		response.RespondSuccessJSON(c.Writer, err, msg)
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

	// Check if user already join post test
	exists, answer := e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, "posttest")
	if exists {
		answerObj := e.EvaluationUsecase.CheckAnswerResult(answer.Selected)
		if answerObj.IsPass() {
			msg := "Anda telah lulus post test"
			err := struct{}{}
			response.RespondSuccessJSON(c.Writer, err, msg)
			return
		}
		e.EvaluationUsecase.ArchivedPostAnswer(loggedIn.Username)
	}

	assesment, totalData := e.EvaluationUsecase.StartPostEvaluation(loggedIn.Username)
	if assesment == nil {
		msg := "You have to join pre test first"
		err := struct{}{}
		response.RespondSuccessJSON(c.Writer, err, msg)
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

	msg := "Daftar pertanyaan"
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

	// User session
	loggedIn := e.Middleware.GetLoggedInUser(c)

	// Check if user already pass this quiz
	exists, answer := e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, req.Title)
	if exists {
		answerObj := e.EvaluationUsecase.CheckAnswerResult(answer.Selected)
		if answerObj.IsPass() {
			msg := "Anda telah lulus kuis"
			err := struct{}{}
			response.RespondSuccessJSON(c.Writer, err, msg)
			return
		}
	}

	assesment, totalData := e.EvaluationUsecase.StartEvaluation(req.Title, req.Page, 5)

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

	msg := "Daftar pertanyaan"
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
	exists, _ := e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, "pretest")
	if exists {
		msg := "Anda telah mengikuti pre test"
		err := struct{}{}
		response.RespondSuccessJSON(c.Writer, err, msg)
		return
	}

	// Remove double quote start & end from answer
	req.Answer = req.Answer[1 : len(req.Answer)-1]

	answerObj := e.EvaluationUsecase.CheckAnswerResult(req.Answer)
	e.EvaluationUsecase.SaveAnswer(loggedIn.Username, "pretest", req.Answer, answerObj.Grade)
	e.LearningUsecase.SetDefaultCourse(loggedIn.Username)

	msg := "Terima kasih, kami telah menerima jawaban anda"
	res := struct{}{}
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
	exists, _ := e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, "pretest")
	if !exists {
		msg := "Anda harus mengikuti pre test terlebih dahulu"
		err := struct{}{}
		response.RespondSuccessJSON(c.Writer, err, msg)
		return
	}
	// Check if user already join post test
	exists, _ = e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, "posttest")
	if exists {
		msg := "Anda telah mengikuti post test"
		err := struct{}{}
		response.RespondSuccessJSON(c.Writer, err, msg)
		return
	}

	// Remove double quote start & end from answer
	req.Answer = req.Answer[1 : len(req.Answer)-1]

	answerObj := e.EvaluationUsecase.CheckAnswerResult(req.Answer)
	e.EvaluationUsecase.SaveAnswer(loggedIn.Username, "posttest", req.Answer, answerObj.Grade)

	grade := fmt.Sprintf("%.1f", answerObj.Grade)
	msg := "Your submission grade " + grade + "%"
	res := ProcessPostAnswerResponse{Grade: grade}
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
	// Check if user already pass this quiz
	exists, answer := e.EvaluationUsecase.IsAnswerExists(loggedIn.Username, req.Title)
	if exists {
		answerObj := e.EvaluationUsecase.CheckAnswerResult(answer.Selected)
		if answerObj.IsPass() {
			msg := "Anda telah lulus kuis"
			err := struct{}{}
			response.RespondSuccessJSON(c.Writer, err, msg)
			return
		}
	}

	// Remove double quote start & end from answer
	req.Answer = req.Answer[1 : len(req.Answer)-1]

	// Check quiz result
	answerObj := e.EvaluationUsecase.CheckAnswerResult(req.Answer)

	e.EvaluationUsecase.ArchivedQuizAnswer(loggedIn.Username, req.Title)
	e.EvaluationUsecase.SaveAnswer(loggedIn.Username, req.Title, req.Answer, answerObj.Grade)

	// User pass the quiz
	if answerObj.IsPass() {
		e.LearningUsecase.UpdateUserCourseProgress(loggedIn.Username, req.Title)
	}

	msg := "Terima kasih, kami telah menerima jawaban anda"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// ResetPrePostStatus set current pre post answer status to archived
func (e *EvaluationHandler) ResetPrePostStatus(c *gin.Context) {
	// User session
	loggedIn := e.Middleware.GetLoggedInUser(c)
	e.EvaluationUsecase.ArchivedPrePostAnswer(loggedIn.Username)

	e.LearningUsecase.SetDefaultCourse(loggedIn.Username)

	// Response
	msg := "Pre test status anda telah diatur ulang"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// ResetPostStatus set current post answer status to archived
func (e *EvaluationHandler) ResetPostStatus(c *gin.Context) {
	// User session
	loggedIn := e.Middleware.GetLoggedInUser(c)
	e.EvaluationUsecase.ArchivedPostAnswer(loggedIn.Username)

	// Response
	msg := "Post test status anda telah diatur ulang"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// PostTestResult return user post test result
func (e *EvaluationHandler) PostTestResult(c *gin.Context) {
	// User session
	loggedIn := e.Middleware.GetLoggedInUser(c)

	// Result
	result := e.EvaluationUsecase.PostTestResult(loggedIn.Username)

	// Post test status
	postTestStatus := 1
	courses := e.LearningUsecase.GetCourseList()
	for _, course := range courses {
		_, statusCode := course.GetParticipantStatus(loggedIn.Username)
		if statusCode < 2 {
			postTestStatus = 0
		}
	}

	// User complete past test
	certificate := ""
	if result.Pass && postTestStatus == 1 {
		postTestStatus = 2
		certificate = "Coursera%20XZS4R52UDZFN.pdf"
	}

	// Response
	grade := fmt.Sprintf("%.f", result.Grade)
	msg := "To Pass get 80%"
	res := PostTestResultResponse{
		Grade:          grade + "%",
		Pass:           result.Pass,
		Status:         postTestStatus,
		CertificateUrl: certificate,
	}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// AddQuestion add new question
func (e *EvaluationHandler) AddQuestion(c *gin.Context) {
	groupName := c.Params.ByName("name")
	// Form Data
	var req RequestAddQuestion
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

	if req.QuestionType == "prepost" {
		groupName = "prepost"
	}

	e.EvaluationUsecase.AddQuestion(req.Question, groupName, req.Choices, req.Answer)

	// Response
	msg := "Pertanyaan baru telah ditambahkan"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// AddQuestion add new question
func (e *EvaluationHandler) EditQuestion(c *gin.Context) {
	// Form Data
	var req RequestEditQuestion
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
	id := c.Params.ByName("id")

	i, err := strconv.Atoi(id)
	e.EvaluationUsecase.EditQuestion(i, req.Question, req.Choices, req.Answer)

	// Response
	msg := "Pertanyaan ini telah diupdate"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// DeleteQuestion delete a question
func (e *EvaluationHandler) DeleteQuestion(c *gin.Context) {
	id := c.Params.ByName("id")
	i, _ := strconv.Atoi(id)

	e.EvaluationUsecase.DeleteQuestion(i)

	// Response
	msg := "Pertanyaan ini telah dihapus"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// DeleteQuestion all question list
func (e *EvaluationHandler) ListOfQuestion(c *gin.Context) {
	// Form Data
	var req RequestListQuestion
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

	assesment, totalData := e.EvaluationUsecase.ListQuestion(req.Page, req.Limit)
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
	res := ListQuestionResponse{
		Pagination: PaginationResponse{
			TotalData:   totalData,
			TotalPage:   totalPage,
			Limit:       10,
			Current:     1,
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
		res.Questions = append(res.Questions, q)
	}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

func (e *EvaluationHandler) ListOfQuestionsInGroup(c *gin.Context) {
	groupName := c.Params.ByName("name")

	// Form Data
	var req RequestListQuestion
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

	assesment, totalData := e.EvaluationUsecase.GetQuestionByGroup(groupName, req.Page, req.Limit)
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

	msg := "Daftar pertanyaan"
	res := ListQuestionResponse{
		Pagination: PaginationResponse{
			TotalData:   totalData,
			TotalPage:   totalPage,
			Limit:       10,
			Current:     1,
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
		res.Questions = append(res.Questions, q)
	}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

// ListOfGroupsQuestion list of groups
func (e *EvaluationHandler) ListOfGroupsQuestion(c *gin.Context) {
	// Form Data
	var req RequestListOfGroupsQuestion
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

	questionGroupResp := []QuestionGroup{}

	questionGroups := []*model.QuestionGroup{}
	if req.Status != "available" {
		questionGroups = e.EvaluationUsecase.GetAvailableQuestionGroups()
	} else {
		questionGroups = e.EvaluationUsecase.GetQuestionGroups()
	}

	for _, group := range questionGroups {
		questionType := "prepost"
		if group.Name != "prepost" {
			questionType = "quiz"
		}
		questionGroupResp = append(questionGroupResp, QuestionGroup{
			Name:           group.Name,
			Type:           questionType,
			TotalQuestions: group.TotalQuestion,
		})
	}

	msg := "List of groups question"
	res := &ListOfGroupsQuestionResponse{
		Questions: questionGroupResp,
	}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

// DeleteQuestion
func (e *EvaluationHandler) DeleteGroupQuestion(c *gin.Context) {
	name := c.Params.ByName("name")

	e.EvaluationUsecase.RemoveGroupByname(name)

	// Response
	msg := "This question group has been deleted"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}
