package evaluation

import "github.com/unduu/e-learning/evaluation/model"

type Usecase interface {
	StartEvaluation(module string, page int, limit int) (model.Assesment, int)
	StartPostEvaluation(username string, page int, limit int) (*model.Assesment, int)
	IsPrePostExists(username string) bool
	CheckAnswerResult(answer string)
	SaveAnswer(username, testType, answer string)
}
