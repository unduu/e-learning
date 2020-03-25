package evaluation

import "github.com/unduu/e-learning/evaluation/model"

type Repository interface {
	GetQuestions(module string, page int, limit int) ([]*model.Question, int, error)
	GetQuestionByIds(ids []string, page int, limit int) ([]*model.Question, int, error)
	GetCorrectAnswerByQuestionID(id int) *model.CorrectAnswerDB
	GetUserAnswers(username string) *model.UserAnswerDB
	InsertAnswer(username, testType, answer string) (affected int64)
}
