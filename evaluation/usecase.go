package evaluation

import "github.com/unduu/e-learning/evaluation/model"

type Usecase interface {
	StartEvaluation(module string, page int, limit int) (model.Assesment, int)
	StartPostEvaluation(username string) (*model.Assesment, int)
	IsAnswerExists(username string, module string) (isExist bool, answer *model.UserAnswerDB)
	CheckAnswerResult(answer string) *model.Answer
	PostTestResult(username string) *model.Result
	SaveAnswer(username string, testType string, answer string, grade float64)
	ArchivedPrePostAnswer(username string)
	ArchivedPostAnswer(username string)
	ArchivedQuizAnswer(username string, quizName string)
	AddQuestion(question string, module string, option string, answer string)
	EditQuestion(id int, question string, option string, answer string)
	DeleteQuestion(id int)
	ListQuestion(page int, limit int) (model.Assesment, int)
}
