package http

// Request data
type RequestEvaluation struct {
	Page  int `form:"page,default=1" json:"page" xml:"page"  binding:"number"`
	Limit int `form:"limit,default=5" json:"limit" xml:"limit" binding:"number"`
}

// Request data
type RequestProcessEvaluationAnswer struct {
	Answer string `form:"answer" json:"answer" xml:"answer"  binding:"required"`
}
