package http

import (
	"github.com/gin-gonic/gin"
	customValidator "github.com/unduu/e-learning/helper/validator"
	"github.com/unduu/e-learning/learning"
	"github.com/unduu/e-learning/learning/model"
	"github.com/unduu/e-learning/middleware"
	"github.com/unduu/e-learning/response"

	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"strconv"
)

type LearningHandler struct {
	LearningUsecase learning.Usecase
	Middleware      *middleware.Middleware
	Validator       *customValidator.CustomValidator
}

func NewHttpLearningHandler(router *gin.RouterGroup, mw *middleware.Middleware, v *customValidator.CustomValidator, learningUC learning.Usecase) {
	handler := &LearningHandler{
		LearningUsecase: learningUC,
		Middleware:      mw,
		Validator:       v,
	}
	router.GET("module", mw.AuthMiddleware, handler.ModuleList)
	router.GET("module/:alias/content", mw.AuthMiddleware, handler.LearningContent)
	router.POST("module/:alias/:lesson/video", mw.AuthMiddleware, handler.SaveVideoProgress)

	router.POST("module", mw.AuthMiddleware, handler.AddCourse)
	router.PUT("module/:alias", mw.AuthMiddleware, handler.EditCourse)
	router.DELETE("module/:alias", mw.AuthMiddleware, handler.DeleteCourse)

	router.POST("course/:alias/:section/lessons", mw.AuthMiddleware, handler.AddCourseContent)
}

// ModuleList return list of courses / modules
func (l *LearningHandler) ModuleList(c *gin.Context) {
	// Logged in User Session
	loggedIn := l.Middleware.GetLoggedInUser(c)

	// Processing
	courses := l.LearningUsecase.GetCourseList()
	moduleArr := []Module{}
	for _, course := range courses {
		totalLessons := strconv.Itoa(course.GetTotalLesson()) + " Lessons"
		totalSections := strconv.Itoa(len(course.Sections)) + " Sections"

		// Formatting time duration
		courseDuration := model.CourseDuration{Duration: course.CountDuration()}

		status, statusCode := course.GetParticipantStatus(loggedIn.Username)

		c := Module{
			course.Alias,
			course.Title,
			course.Subtitle,
			totalLessons,
			courseDuration.Format(),
			status,
			statusCode,
			course.Thumbnail,
			totalSections,
		}
		moduleArr = append(moduleArr, c)
	}

	// Response
	msg := "Learning module list"
	res := ResponseModuleList{Modules: moduleArr}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

// LearningContent return course content
func (l *LearningHandler) LearningContent(c *gin.Context) {
	// Path param
	alias := c.Params.ByName("alias")
	// Session
	loggedIn := l.Middleware.GetLoggedInUser(c)
	// Processing
	contentArr := []Section{}
	course := l.LearningUsecase.GetCourseLessons(alias, loggedIn.Username)
	_, statusCode := course.GetParticipantStatus(loggedIn.Username)
	// Response cannot access course
	if statusCode == 0 && loggedIn.Role != "admin" {
		msg := "You not allowed to access this course"
		err := make([]string, 0)
		response.RespondErrorJSON(c.Writer, err, msg)
		return
	}

	for i, sectionObj := range course.Sections {
		lessonResArr := []Lesson{}
		for _, lessonObj := range sectionObj.Lessons {
			// Formatting time duration
			courseDuration := model.CourseDuration{Duration: lessonObj.Duration}
			l.LearningUsecase.SetLessonProgress(loggedIn.Username, lessonObj)
			// Format lessonsplit question choices
			lessonSplits := make([]LessonSplit, 0)

			for _, split := range lessonObj.Split {
				split.FormatChoices()
				lessonSplits = append(lessonSplits, LessonSplit{
					Type:    split.Type,
					Video:   split.Video,
					Answer:  split.Answer,
					Choices: split.Choices,
				})
			}
			lessonRes := Lesson{
				Type:         lessonObj.Type,
				Title:        lessonObj.Title,
				Permalink:    lessonObj.Permalink,
				Duration:     courseDuration.Minute(),
				Video:        lessonObj.Video,
				Timebar:      lessonObj.Timebar,
				Progress:     lessonObj.GetProgressName(),
				ProgressCode: lessonObj.Progress,
				LessonSplit:  lessonSplits,
			}
			lessonResArr = append(lessonResArr, lessonRes)
		}

		// Set status to open for first section
		sectionStatus, sectionCode := "open", 1
		if i != 0 {
			sectionStatus, sectionCode = course.Sections[i-1].GetParticipantStatus(loggedIn.Username)
		}

		sectionRes := Section{
			Section:    sectionObj.Name,
			Name:       sectionObj.Desc,
			Lessons:    lessonResArr,
			Status:     sectionStatus,
			StatusCode: sectionCode,
		}
		contentArr = append(contentArr, sectionRes)
	}

	if len(course.Sections) <= 0 {
		msg := "Course not found"
		err := response.Error{"alias", "Enter a valid course"}
		response.RespondErrorJSON(c.Writer, err, msg)
		return
	}

	// Response
	msg := "Learning module list"
	res := ResponseLearningContent{Content: contentArr}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

// SaveVideoProgress Save user video time progress to resume later
func (l *LearningHandler) SaveVideoProgress(c *gin.Context) {
	// Path param
	alias := c.Params.ByName("alias")
	lesson := c.Params.ByName("lesson")
	// Form Data
	var req RequestSaveVideoProgress
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
			e := fieldErr.Translate(l.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}
	// Session
	loggedIn := l.Middleware.GetLoggedInUser(c)
	// Processing
	l.LearningUsecase.UpdateVideoProgress(loggedIn.Username, alias, lesson, req.Timebar)
	// Response
	msg := "User video progres has been saved"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// AddCourse Add new course
func (l *LearningHandler) AddCourse(c *gin.Context) {
	// Form Data
	var req RequestAddCourse
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
			e := fieldErr.Translate(l.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}

	l.LearningUsecase.AddCourse(req.Title, req.Subtitle, req.Thumbnail)

	// Response
	msg := "New Course has been added"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// EditCourse update course content
func (l *LearningHandler) EditCourse(c *gin.Context) {
	// Form Data
	var req RequestAddCourse
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
			e := fieldErr.Translate(l.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}
	alias := c.Params.ByName("alias")

	l.LearningUsecase.EditCourse(alias, req.Title, req.Subtitle, req.Thumbnail)

	// Response
	msg := "Course data has been updated"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// DeleteCourse delete course
func (l *LearningHandler) DeleteCourse(c *gin.Context) {
	// Form Data
	alias := c.Params.ByName("alias")

	l.LearningUsecase.DeleteCourse(alias)

	// Response
	msg := "Course has been deleted"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}

// AddCourseContent
func (l *LearningHandler) AddCourseContent(c *gin.Context) {
	course := c.Params.ByName("alias")
	section := c.Params.ByName("section")

	// Form Data
	var req RequestAddCourseContent
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
			e := fieldErr.Translate(l.Validator.Translation)

			error := response.Error{fieldErr.Field(), e}
			errValidation = append(errValidation, error)
		}
		response.RespondErrorJSON(c.Writer, errValidation)
		return
	}

	l.LearningUsecase.AddCourseContent(course, section, req.SectionDesc, req.Type, req.Title, req.Video)

	// Response
	msg := "New lesson has been added"
	res := struct{}{}
	response.RespondSuccessJSON(c.Writer, res, msg)
}
