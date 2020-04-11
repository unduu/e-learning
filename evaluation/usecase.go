package evaluation

import "github.com/unduu/e-learning/evaluation/model"

type Usecase interface {
	StartEvaluation(module string, page int, limit int) (model.Assesment, int)
	StartPostEvaluation(username string, page int, limit int) (*model.Assesment, int)
	IsAnswerExists(username string, module string) (isExist bool, answer *model.UserAnswerDB)
	CheckAnswerResult(answer string) *model.Answer
	SaveAnswer(username string, testType string, answer string, grade float64)
	ArchivedPrePostAnswer(username string)
	ArchivedPostAnswer(username string)
	ArchivedQuizAnswer(username string, quizName string)
}
