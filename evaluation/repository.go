package evaluation

import "github.com/e-learning/evaluation/model"

type Repository interface {
	GetQuestions(page int, limit int) ([]*model.Question, int, error)
}
