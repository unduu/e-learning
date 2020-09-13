package model

type CorrectAnswerDB struct {
	Id       int    `db:"id"`
	Selected string `db:"answer"`
}

type UserAnswerDB struct {
	Username string  `db:"username"`
	Selected string  `db:"answer"`
	Grade    float64 `db:"grade"`
}

type Answer struct {
	Id           int   `json:"id" db:"id"`
	Selected     []int `json:"answer" db:"answer"`
	TotalAnswer  int
	TotalWrong   int
	TotalCorrect int
	Grade        float64
}

func (a *Answer) IsPass() bool {
	if a.Grade >= 80 {
		return true
	}
	return false
}
