package learning

import "github.com/unduu/e-learning/learning/model"

type Usecase interface {
	GetCourseList() (courseArr []*model.Course)
	GetCourseLessons(alias string) (course *model.Course)
	SetDefaultCourse(username string)
	UpdateUserCourseProgress(username string, quiz string)
	SetLessonProgress(username string, lesson *model.Lesson) *model.Lesson
}
