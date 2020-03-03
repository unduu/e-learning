package repository

import (
	"fmt"
	"github.com/e-learning/evaluation/model"
	"github.com/jmoiron/sqlx"
)

type EvaluationRepository struct {
	conn *sqlx.DB
}

func NewEvaluationRepository(db *sqlx.DB) *EvaluationRepository {
	return &EvaluationRepository{
		conn: db,
	}
}

func (a *EvaluationRepository) GetQuestions(page int, limit int) ([]*model.Question, int, error) {
	offset := (page - 1) * limit
	questions := make([]*model.Question, 0)

	var count struct {
		Total int `db:"total"`
	}
	queryParams := map[string]interface{}{
		"offset": offset,
		"limit":  limit,
	}
	query, err := a.conn.PrepareNamed(`SELECT id, type, attachment_type, attachment, question, choices FROM questions LIMIT :offset, :limit `)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	queryTotal, err := a.conn.PrepareNamed(`SELECT COUNT(*) AS total FROM questions`)
	if err != nil {
		fmt.Println("Error db GetQuestions->PrepareNamed : ", err)
	}

	err = query.Select(&questions, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	err = queryTotal.Get(&count, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	return questions, count.Total, err
}

func (a *EvaluationRepository) GetAnswerByQuestionID(id int) *model.AnswerDB {
	// DB Response struct
	answers := make([]*model.AnswerDB, 0)

	// Data for query
	queryParams := map[string]interface{}{
		"id": id,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`SELECT id,answer FROM questions WHERE id = :id`)
	if err != nil {
		fmt.Println("Error db GetAnswerByQuestionID->PrepareNamed : ", err)
	}

	// Execute query
	err = query.Select(&answers, queryParams)
	if err != nil {
		fmt.Println("Error db GetQuestions->query.Get : ", err)
	}

	return answers[0]
}
