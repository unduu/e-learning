package evaluation

import "github.com/unduu/e-learning/evaluation/model"

type Usecase interface {
	StartEvaluation(page int, limit int) (model.Assesment, int)
	CheckAnswerResult(answer string)
	SaveAnswer(username, testType, answer string)
}
