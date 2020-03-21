package learning

import "github.com/unduu/e-learning/learning/model"

type Usecase interface {
	GetCourseList() (courseArr []*model.Course)
}
