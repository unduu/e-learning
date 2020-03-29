package usecase

import (
	"github.com/unduu/e-learning/learning"
	"github.com/unduu/e-learning/learning/model"
)

type LearningUsecase struct {
	repository learning.Repository
}

func NewLearningUsecase(repository learning.Repository) *LearningUsecase {
	return &LearningUsecase{
		repository: repository,
	}
}

// GetCourseList return course / modules list
func (a *LearningUsecase) GetCourseList() (results []*model.Course) {
	courseArr := a.repository.GetCourses()
	for _, course := range courseArr {
		participants := a.repository.GetParticipantByCourse(course.Id)
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
