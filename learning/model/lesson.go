package model

import "strings"

type Lesson struct {
	Type      string `db:"type"`
	Title     string `db:"title"`
	Permalink string `db:"permalink"`
	Duration  int    `db:"duration"`
	Video     string `db:"content"`
	Timebar   int
	Progress  int
	Split     []*LessonSplit
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

func (l *Lesson) GeneratePermalink() {
	lower := strings.ToLower(l.Title)
	l.Permalink = strings.Replace(lower, " ", "-", -1)
}
