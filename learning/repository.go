package learning

import "github.com/unduu/e-learning/learning/model"

type Repository interface {
	GetCourses() []*model.Course
	GetParticipantByCourse(id int) []*model.Participant
}
