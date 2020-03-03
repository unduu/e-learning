package model

type AnswerDB struct {
	Id       int    `db:"id"`
	Selected string `db:"answer"`
}

type Answer struct {
	Id       int   `json:"id" db:"id"`
	Selected []int `json:"answer" db:"answer"`
}
