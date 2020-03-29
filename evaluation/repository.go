package evaluation

import "github.com/unduu/e-learning/evaluation/model"

type Repository interface {
	GetQuestions(module string, page int, limit int) ([]*model.Question, int, error)
	GetQuestionByIds(ids []string, page int, limit int) ([]*model.Question, int, error)
	GetCorrectAnswerByQuestionID(id int) *model.CorrectAnswerDB
	GetUserAnswers(username string, module string) *model.UserAnswerDB
	InsertAnswer(username string, testType string, answer string) (affected int64)
	UpdateUserAnswerStatus(username string, module string, newStatus string) (affected int64)
}
