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

func (l *Lesson) GetProgressName() (status string) {
	switch l.Progress {
	case 0:
		status = "open"
	case 1:
		status = "failed"
	case 2:
		status = "complete"
	}
	return status
}
