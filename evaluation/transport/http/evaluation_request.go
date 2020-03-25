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
