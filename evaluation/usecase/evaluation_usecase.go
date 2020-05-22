package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/unduu/e-learning/evaluation"
	"github.com/unduu/e-learning/evaluation/model"
	"strconv"
	"time"
)

type EvaluationUsecase struct {
	repository evaluation.Repository
}

func NewEvaluationUsecase(repository evaluation.Repository) *EvaluationUsecase {
	return &EvaluationUsecase{
		repository: repository,
	}
}

// StartEvaluation return list of pre test question
func (a *EvaluationUsecase) StartEvaluation(module string, page int, limit int) (model.Assesment, int) {
	assesment := model.Assesment{Status: "active"}
	assesment.SetDuration()

	questionsList, totalData, err := a.repository.GetQuestions(module, page, limit)

	if err != nil {
		fmt.Println("ERROR StartEvaluation : ", err)
	}

	for _, questionRow := range questionsList {
		// Decode choices from db then set as choice struct
		cc := model.Choice{}
		err := json.Unmarshal([]byte(questionRow.ChoicesDB), &cc)
		if err != nil {
			fmt.Println("ERROR StartEvaluation->Unmarshal : ", err)
		}
		questionRow.Choices = cc
		assesment.AddQuestion(questionRow)
	}

	return assesment, totalData
}

// StartPostEvaluation return list of answered pre test question
func (a *EvaluationUsecase) StartPostEvaluation(username string, page int, limit int) (*model.Assesment, int) {
	questionIdArr := []string{}
	answerArr := []model.Answer{}

	assesment := &model.Assesment{Status: "active"}
	assesment.SetDuration()

	answer := a.repository.GetUserAnswers(username, "pretest")
	if answer == nil {
		return nil, 0
	}

	err := json.Unmarshal([]byte(answer.Selected), &answerArr)
	if err != nil {
		fmt.Println("ERROR StartPostEvaluation->Unmarshal ", err)
	}

	for _, row := range answerArr {
		// List of user pre test question
		questionIdArr = append(questionIdArr, strconv.Itoa(row.Id))
	}
	questionsList, totalData, err := a.repository.GetQuestionByIds(questionIdArr, 1, 20)

	if err != nil {
		fmt.Println("ERROR StartEvaluation : ", err)
	}

	for _, questionRow := range questionsList {
		// Decode choices from db then set as choice struct
		cc := model.Choice{}
		err := json.Unmarshal([]byte(questionRow.ChoicesDB), &cc)
		if err != nil {
			fmt.Println("ERROR StartEvaluation->Unmarshal : ", err)
		}
		questionRow.Choices = cc
		assesment.AddQuestion(questionRow)
	}

	return assesment, totalData
}

// IsAnswerExists check if user answer already exists in pre / post test
func (a *EvaluationUsecase) IsAnswerExists(username string, module string) (isExist bool, answer *model.UserAnswerDB) {
	answer = a.repository.GetUserAnswers(username, module)
	if answer == nil {
		return false, nil
	}
	return true, answer
}

// CheckAnswerResult compare user test answer with correct answer
func (a *EvaluationUsecase) CheckAnswerResult(answer string) *model.Answer {

	var totalRightAnswer int
	var totalWrongAnswer int
	var totalQuestions int
	var grade float64

	answerArr := []model.Answer{}

	err := json.Unmarshal([]byte(answer), &answerArr)
	if err != nil {
		fmt.Println("ERROR CompareAnswer->Unmarshal ", err)
	}

	for _, row := range answerArr {
		// Answer from user
		selectedAnswerJson, err := json.Marshal(row.Selected)
		selectedAnswer := string(selectedAnswerJson)
		if err != nil {
			fmt.Println("ERROR CompareAnswer->Marshal selected answer", err)
		}

		// Answer from db
		rightAnswerObj := a.repository.GetCorrectAnswerByQuestionID(row.Id)
		rightAnswer := rightAnswerObj.Selected

		// Answer right ?
		if selectedAnswer == rightAnswer {
			totalRightAnswer++
		}
	}
	totalQuestions = len(answerArr)
	totalWrongAnswer = totalQuestions - totalRightAnswer
	if totalRightAnswer >= 0 && totalQuestions > 0 {
		grade = float64(totalRightAnswer) / float64(totalQuestions) * 100
	}

	return &model.Answer{
		TotalAnswer:  totalQuestions,
		TotalWrong:   totalWrongAnswer,
		TotalCorrect: totalRightAnswer,
		Grade:        grade,
	}
}

func (a *EvaluationUsecase) PostTestResult(username string) *model.Result {
	result := &model.Result{}
	// Check if user already join post test
	exists, answer := a.IsAnswerExists(username, "posttest")
	if exists {
		answerObj := a.CheckAnswerResult(answer.Selected)
		result.Pass = answerObj.IsPass()
		result.Grade = answerObj.Grade
	}
	return result
}

// SaveAnswer insert user answer to databases
func (a *EvaluationUsecase) SaveAnswer(username string, testType string, answer string, grade float64) {
	a.repository.InsertAnswer(username, testType, answer, grade)
}

// ArchivedPrePostAnswer reset user pre post test, so user can retry pre post test
func (a *EvaluationUsecase) ArchivedPrePostAnswer(username string) {
	t := time.Now()
	timeStr := t.Format("20060102150405")
	archivedName := "archived_" + timeStr
	a.repository.UpdateUserAnswerStatus(username, "pretest", archivedName)
	a.repository.UpdateUserAnswerStatus(username, "posttest", archivedName)
}

// ArchivedPrePostAnswer reset user pre post test, so user can retry pre post test
func (a *EvaluationUsecase) ArchivedPostAnswer(username string) {
	t := time.Now()
	timeStr := t.Format("20060102150405")
	archivedName := "archived_" + timeStr
	a.repository.UpdateUserAnswerStatus(username, "posttest", archivedName)
}

func (a *EvaluationUsecase) ArchivedQuizAnswer(username string, quizName string) {
	t := time.Now()
	timeStr := t.Format("20060102150405")
	archivedName := "archived_" + timeStr
	a.repository.UpdateUserAnswerStatus(username, quizName, archivedName)
}

// AddQuestion add a new question
func (a *EvaluationUsecase) AddQuestion(question string, module string, option string, answer string) {
	q := model.Question{
		Module:     module,
		Type:       "multiple_choice",
		AttachType: "",
		Attachment: "",
		Text:       question,
		ChoicesDB:  option,
		Answer:     answer,
	}
	a.repository.InsertQuestion(q)
}

// EditQuestion change question content
func (a *EvaluationUsecase) EditQuestion(id int, question string, option string, answer string) {
	q := model.Question{
		Id:        id,
		Text:      question,
		ChoicesDB: option,
		Answer:    answer,
	}
	a.repository.UpdateQuestion(q)
}

// DeleteQuestion remove question
func (a *EvaluationUsecase) DeleteQuestion(id int) {
	q := model.Question{
		Id: id,
	}
	a.repository.DeleteQuestion(q)
}

func (a *EvaluationUsecase) ListQuestion(page int, limit int) (model.Assesment, int) {
	assesment := model.Assesment{Status: "active"}
	assesment.SetDuration()

	questionsList, totalData, err := a.repository.GetAllQuestions(page, limit)
	if err != nil {
		fmt.Println("ERROR StartEvaluation : ", err)
	}
	for _, questionRow := range questionsList {
		// Decode choices from db then set as choice struct
		cc := model.Choice{}
		err := json.Unmarshal([]byte(questionRow.ChoicesDB), &cc)
		if err != nil {
			fmt.Println("ERROR StartEvaluation->Unmarshal : ", err)
		}
		questionRow.Choices = cc
		assesment.AddQuestion(questionRow)
	}
	return assesment, totalData
}
