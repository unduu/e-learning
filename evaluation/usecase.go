package evaluation

import "github.com/unduu/e-learning/evaluation/model"

type Usecase interface {
	StartEvaluation(module string, page int, limit int) (model.Assesment, int)
	StartPostEvaluation(username string, page int, limit int) (*model.Assesment, int)
	IsAnswerExists(username string, module string) bool
	CheckAnswerResult(answer string)
	SaveAnswer(username string, testType string, answer string)
	ArchivedPrePostAnswer(username string)
}
