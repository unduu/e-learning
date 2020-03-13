package model

type CorrectAnswerDB struct {
	Id       int    `db:"id"`
	Selected string `db:"answer"`
}

type UserAnswerDB struct {
	Username string `db:"username"`
	Selected string `db:"answer"`
}

type Answer struct {
	Id       int   `json:"id" db:"id"`
	Selected []int `json:"answer" db:"answer"`
}
