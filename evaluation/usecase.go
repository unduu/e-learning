package evaluation

import "github.com/e-learning/evaluation/model"

type Usecase interface {
	StartEvaluation(page int, limit int) (model.Assesment, int)
	CheckAnswerResult(answer string)
	SaveAnswer(username, testType, answer string)
}
