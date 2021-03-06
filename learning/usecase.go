package learning

import "github.com/unduu/e-learning/learning/model"

type Usecase interface {
	GetCourseList() (courseArr []*model.Course)
	GetCourseLessons(alias string, username string) (course *model.Course)
	SetDefaultCourse(username string)
	UpdateUserCourseProgress(username string, quiz string)
	SetLessonProgress(username string, lesson *model.Lesson) *model.Lesson
	UpdateVideoProgress(username string, course string, lesson string, time int)
	AddCourse(title string, subtitle string, thumbnail string)
	EditCourse(alias string, title string, subtitle string, thumbnail string)
	DeleteCourse(alias string)
	AddCourseContent(courseAlias string, sectionName string, sectionDesc string, module string, title string, video string)
}
