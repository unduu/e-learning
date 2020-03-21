package http

import (
	"github.com/gin-gonic/gin"
	customValidator "github.com/unduu/e-learning/helper/validator"
	"github.com/unduu/e-learning/learning"
	"github.com/unduu/e-learning/learning/model"
	"github.com/unduu/e-learning/middleware"
	"github.com/unduu/e-learning/response"
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

}

func (l *LearningHandler) ModuleList(c *gin.Context) {
	// Logged in User Session
	loggedIn := l.Middleware.GetLoggedInUser(c)

	// Processing
	courses := l.LearningUsecase.GetCourseList()
	moduleArr := []Module{}
	for _, course := range courses {
		totalLessons := strconv.Itoa(course.TotalLesson) + " Lessons"

		// Formatting time duration
		courseDuration := model.NewCourseDuration(course.Duration)

		status, statusCode := course.GetParticipantStatus(loggedIn.Username)

		c := Module{
			course.Alias,
			course.Title,
			course.Subtitle,
			totalLessons,
			courseDuration.Format(),
			status,
			statusCode,
		}
		moduleArr = append(moduleArr, c)
	}

	// Response
	msg := "Learning module list"
	res := ModulesResponse{Modules: moduleArr}

	response.RespondSuccessJSON(c.Writer, res, msg)
}
