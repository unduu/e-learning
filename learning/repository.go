package learning

import "github.com/unduu/e-learning/learning/model"

type Repository interface {
	GetCourses() []*model.Course
	GetParticipantByCourse(id int) []*model.Participant
	GetCourseByAlias(alias string) *model.Course
	GetLessonsByCourseId(id int) []*model.SectionLessons
	AddCourseParticipant(username string, courseId int, status int) (affected int64)
	DeleteUserFromAllCourse(username string) (affected int64)
}
