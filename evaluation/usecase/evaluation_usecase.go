package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/unduu/e-learning/evaluation"
	"github.com/unduu/e-learning/evaluation/model"
	"strconv"
)

type EvaluationUsecase struct {
	repository evaluation.Repository
}

func NewEvaluationUsecase(repository evaluation.Repository) *EvaluationUsecase {
	return &EvaluationUsecase{
		repository: repository,
	}
}

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

func (a *EvaluationUsecase) StartPostEvaluation(page int, limit int) (model.Assesment, int) {
	questionIdArr := []string{}
	answerArr := []model.Answer{}

	assesment := model.Assesment{Status: "active"}
	assesment.SetDuration()

	answer := a.repository.GetUserAnswers("johndoe")

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

func (a *EvaluationUsecase) CheckAnswerResult(answer string) {

	var grade float64
	var totalRightAnswer int
	var totalWrongAnswer int
	var totalQuestions int

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

	fmt.Println("Total Question ", totalQuestions)
	fmt.Println("Total Right Answer ", totalRightAnswer)
	fmt.Println("Total Wrong Answer ", totalWrongAnswer)
	fmt.Println("Grade ", grade)
}

func (a *EvaluationUsecase) SaveAnswer(username string, testType string, answer string) {

	a.repository.InsertAnswer(username, testType, answer)
}
