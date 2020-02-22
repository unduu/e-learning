package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/e-learning/evaluation"
	"github.com/e-learning/evaluation/model"
)

type EvaluationUsecase struct {
	repository evaluation.Repository
}

func NewEvaluationUsecase(repository evaluation.Repository) *EvaluationUsecase {
	return &EvaluationUsecase{
		repository: repository,
	}
}

func (a *EvaluationUsecase) StartEvaluation(page int, limit int) (model.Assesment, int) {
	assesment := model.Assesment{Start: "2020-02-19 19:30:00", End: "2020-02-19 21:30:00", Status: "active"}

	questionsList, totalData, err := a.repository.GetQuestions(page, limit)

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
