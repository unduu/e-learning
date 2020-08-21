package learning

import "github.com/unduu/e-learning/learning/model"

type Repository interface {
	GetCourses() []*model.Course
	GetParticipantByCourse(id int) []*model.Participant
	GetCourseByAlias(alias string) *model.Course
	GetCourseByQuiz(quiz string) *model.Course
	GetLessonsByCourseId(id int, username string) []*model.SectionLessons
	GetLessonByPermalink(course string, permalink string) *model.SectionLessons
	AddCourseParticipant(username string, courseId int, status int) (affected int64)
	AddLessonProgress(username string, learningId int) (affected int64)
	DeleteUserFromAllCourse(username string) (affected int64)
	DeleteUserAllLessonProgress(username string) (affected int64)
	UpdateParticipantStatus(username string, id int, newStatus int) (affected int64)
	UpdateLearningVideoTimebar(username string, learningId int, time int) (affected int64)
	InsertCourse(course *model.Course) (affected int64)
	UpdateCourse(course *model.Course) (affected int64)
	DeleteCourse(course *model.Course) (affected int64)
	SaveCourseContent(courseId int, sectionName string, sectionDesc string, content *model.Lesson) (affected int64)
	FetchSectionContentByCourseAndSection(courseID int, sectionName string) *model.SectionLessons
}
