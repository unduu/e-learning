package http

// Request data
type RequestEvaluation struct {
	Title string `form:"title" json:"title" xml:"title"  binding:""`
	Page  int    `form:"page,default=1" json:"page" xml:"page"  binding:"number"`
	Limit int    `form:"limit,default=5" json:"limit" xml:"limit" binding:"number"`
}

// Request data
type RequestProcessEvaluationAnswer struct {
	Answer string `form:"answer" json:"answer" xml:"answer"  binding:"required"`
}

// Request form quiz answer
type RequestProcessQuizAnswer struct {
	Title  string `form:"title" json:"title" xml:"title"  binding:"required"`
	Answer string `form:"answer" json:"answer" xml:"answer"  binding:"required"`
}

// Request form add question
type RequestAddQuestion struct {
	Question     string `form:"question" json:"question" xml:"question"  binding:"required"`
	Choices      string `form:"choices" json:"choices" xml:"choices"  binding:"required"`
	Answer       string `form:"answer" json:"answer" xml:"answer"  binding:"required"`
	QuestionType string `form:"type" json:"type" xml:"type"  binding:"oneof=prepost quiz"`
}

// Request form edit question
type RequestEditQuestion struct {
	Question string `form:"question" json:"question" xml:"question"  binding:"required"`
	Choices  string `form:"choices" json:"choices" xml:"choices"  binding:"required"`
	Answer   string `form:"answer" json:"answer" xml:"answer"  binding:"required"`
}

// Request data
type RequestListQuestion struct {
	Page  int `form:"page,default=1" json:"page" xml:"page"  binding:"number"`
	Limit int `form:"limit,default=5" json:"limit" xml:"limit" binding:"number"`
}

type RequestListOfGroupsQuestion struct {
	Status string `form:"status,default=available" json:"status" xml:"status"  binding:"oneof=available not_available"`
	Page   int    `form:"page,default=1" json:"page" xml:"page"  binding:"number"`
	Limit  int    `form:"limit,default=5" json:"limit" xml:"limit" binding:"number"`
}
