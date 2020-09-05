package http

import (
	evaluationModel "github.com/unduu/e-learning/evaluation/model"
)

type Module struct {
	Alias        string `json:"alias"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	TotalLesson  string `json:"total_lessons"`
	TotalHours   string `json:"total_hours"`
	Status       string `json:"status"`
	StatusCode   int    `json:"status_code"`
	Thumbnail    string `json:"thumbnail"`
	TotalSection string `json:"total_sections"`
}

type ResponseModuleList struct {
	Modules []Module `json:"modules"`
}

type LessonSplit struct {
	Type    string                  `json:"type"`
	Video   string                  `json:"video"`
	Answer  string                  `json:"answer"`
	Choices *evaluationModel.Choice `json:"choices"`
}

type Lesson struct {
	Type         string        `json:"type"`
	Title        string        `json:"title"`
	Permalink    string        `json:"permalink"`
	Duration     string        `json:"duration"`
	Video        string        `json:"video"`
	Timebar      int           `json:"timebar"`
	Progress     string        `json:"progress"`
	ProgressCode int           `json:"progress_code"`
	LessonSplit  []LessonSplit `json:"lesson_split"`
}

type Section struct {
	Section    string   `json:"section"`
	Name       string   `json:"name"`
	Lessons    []Lesson `json:"lessons"`
	Status     string   `json:"status"`
	StatusCode int      `json:"status_code"`
}

type ResponseLearningContent struct {
	Content []Section `json:"content"`
}
