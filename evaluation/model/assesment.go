package model

import (
	"time"
)

type Assesment struct {
	Start        string
	End          string
	Status       string
	QuestionList []*Question
}

func (a *Assesment) SetDuration(dur int32) {
	durTime := time.Duration(dur)
	layout := "2006-01-02 15:04:05"
	startTime := time.Now()
	endTime := startTime.Add(time.Minute * durTime)
	a.Start = startTime.Format(layout)
	a.End = endTime.Format(layout)
}

func (a *Assesment) AddQuestion(Q *Question) {
	a.QuestionList = append(a.QuestionList, Q)
}
