package usecase

import (
	"github.com/unduu/e-learning/evaluation"
	"github.com/unduu/e-learning/learning"
	"github.com/unduu/e-learning/learning/model"
)

type LearningUsecase struct {
	repository   learning.Repository
	evaluationUC evaluation.Usecase
}

func NewLearningUsecase(repository learning.Repository, evaluationUC evaluation.Usecase) *LearningUsecase {
	return &LearningUsecase{
		repository:   repository,
		evaluationUC: evaluationUC,
	}
}

// GetCourseList return course / modules list
func (a *LearningUsecase) GetCourseList() (results []*model.Course) {
	courseArr := a.repository.GetCourses()
	for _, course := range courseArr {
		participants := a.repository.GetParticipantByCourse(course.Id)
		// Get detail total lessons & hours
		course := a.GetCourseLessons(course.Alias)
		course.AddParticipant(participants)
		results = append(results, course)
	}

	return results
}

// GetCourseLessons return lessons from a course
func (a *LearningUsecase) GetCourseLessons(alias string) (course *model.Course) {

	course = a.repository.GetCourseByAlias(alias)
	participants := a.repository.GetParticipantByCourse(course.Id)
	course.AddParticipant(participants)
	data := a.repository.GetLessonsByCourseId(course.Id)
	for _, row := range data {
		lesson := &model.Lesson{
			Type:     row.Type,
			Title:    row.Title,
			Duration: row.Duration,
			Video:    row.Video,
			Progress: 0,
		}

		currSection := course.GetSection(row.Name)
		section := currSection
		if currSection == nil {
			section = &model.Section{
				Name: row.Name,
				Desc: row.Desc,
			}
		}
		section.AddLesson(lesson)

		if currSection == nil {
			course.AddSection(section)
		}
	}

	return course
}

// SetDefaultCourse set default course for new registered user
func (a *LearningUsecase) SetDefaultCourse(username string) {
	a.repository.DeleteUserFromAllCourse(username)
	a.repository.AddCourseParticipant(username, 1, 1)
	a.repository.AddCourseParticipant(username, 2, 0)
	a.repository.AddCourseParticipant(username, 3, 0)
	a.repository.AddCourseParticipant(username, 4, 0)
	a.repository.AddCourseParticipant(username, 5, 0)
}

func (a *LearningUsecase) UpdateUserCourseProgress(username string, quiz string) {
	pass := false
	// Get course by quiz name
	course := a.repository.GetCourseByQuiz(quiz)
	// Get lesson from quiz
	lessons := a.repository.GetLessonsByCourseId(course.Id)
	for _, row := range lessons {
		lesson := &model.Lesson{
			Type:     row.Type,
			Title:    row.Title,
			Duration: row.Duration,
			Video:    row.Video,
		}
		lesson = a.SetLessonProgress(username, lesson)
		if lesson.Progress == 1 {
			pass = true
		}
	}
	if pass {
		// Mark current course as finish
		a.repository.UpdateParticipantStatus(username, course.Id, course.GetStatusCode("completed"))
		// Open next course
		a.repository.UpdateParticipantStatus(username, course.GetNextCourseId(), course.GetStatusCode("open"))
	}
}

func (a *LearningUsecase) SetLessonProgress(username string, lesson *model.Lesson) *model.Lesson {
	if lesson.IsQuiz() {
		lesson.Progress = 0
		exist, answer := a.evaluationUC.IsAnswerExists(username, lesson.Title)
		if !exist {
			return lesson
		}

		lesson.Progress = 1
		if answer.Grade == 100 {
			lesson.Progress = 2
		}

	}
	return lesson
}
