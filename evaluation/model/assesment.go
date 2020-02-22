package model

type Assesment struct {
	Start        string
	End          string
	Status       string
	QuestionList []*Question
}

func (a *Assesment) AddQuestion(Q *Question) {
	a.QuestionList = append(a.QuestionList, Q)
}
