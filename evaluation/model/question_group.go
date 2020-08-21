package model

type QuestionGroup struct {
	Name          string `db:"module"`
	Type          string
	TotalQuestion string `db:"total_questions"`
}
