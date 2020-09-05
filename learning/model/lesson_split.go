package model

import (
	"encoding/json"
	"github.com/unduu/e-learning/evaluation/model"
)

type LessonSplit struct {
	Type       string
	Video      string `db:"video_file"`
	ChoicesRaw string `db:"quiz_choices"`
	Answer     string `db:"quiz_answer"`
	Choices    *model.Choice
}

func (s *LessonSplit) FormatChoices() error {
	if s.ChoicesRaw == "" {
		s.Type = "video"
		s.Choices = nil
		return nil
	}
	s.Type = "quiz"
	err := json.Unmarshal([]byte(s.ChoicesRaw), &s.Choices)
	return err
}
