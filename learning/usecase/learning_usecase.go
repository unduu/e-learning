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

func (a *LearningUsecase) GetCourseList() (courseArr []*model.Course) {
	courseArr = a.repository.GetCourses()
	for _, course := range courseArr {
		participants := a.repository.GetParticipantByCourse(course.Id)
		course.AddParticipant(participants)
	}

	return courseArr
}
