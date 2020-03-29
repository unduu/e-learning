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
	router.GET("module/:alias/content", mw.AuthMiddleware, handler.LearningContent)

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
		}
		moduleArr = append(moduleArr, c)
	}

	// Response
	msg := "Learning module list"
	res := ResponseModuleList{Modules: moduleArr}

	response.RespondSuccessJSON(c.Writer, res, msg)
}

// Learning return course content
func (l *LearningHandler) LearningContent(c *gin.Context) {
	// Path param
	alias := c.Params.ByName("alias")
	// Session
	loggedIn := l.Middleware.GetLoggedInUser(c)
	// Processing
	contentArr := []Section{}
	course := l.LearningUsecase.GetCourseLessons(alias)
	_, statusCode := course.GetParticipantStatus(loggedIn.Username)
	// Response cannot access course
	if statusCode == 0 {
		msg := "You not allowed to access this course"
		err := make([]string, 0)
		response.RespondErrorJSON(c.Writer, err, msg)
		return
	}

	for _, sectionObj := range course.Sections {
		lessonResArr := []Lesson{}
		for _, lessonObj := range sectionObj.Lessons {
			// Formatting time duration
			courseDuration := model.CourseDuration{Duration: lessonObj.Duration}
			lessonRes := Lesson{
				Type:     lessonObj.Type,
				Title:    lessonObj.Title,
				Duration: courseDuration.Minute(),
				Video:    lessonObj.Video,
			}
			lessonResArr = append(lessonResArr, lessonRes)
		}
		sectionRes := Section{
			Section: sectionObj.Name,
			Name:    sectionObj.Desc,
			Lessons: lessonResArr,
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
