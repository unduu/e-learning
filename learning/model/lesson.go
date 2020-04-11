package model

import "strings"

type Lesson struct {
	Type     string
	Title    string
	Duration int
	Video    string
	Progress int
}

func (l *Lesson) IsQuiz() bool {
	t := strings.ToLower(l.Type)
	if t == "quiz" {
		return true
	}
	return false
}
