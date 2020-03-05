package evaluation

import "github.com/e-learning/evaluation/model"

type Repository interface {
	GetQuestions(page int, limit int) ([]*model.Question, int, error)
	GetAnswerByQuestionID(id int) *model.AnswerDB
	InsertAnswer(username, testType, answer string) (affected int64)
}
