package http

// Request data
type RequestPreEvaluation struct {
	Page  int `form:"page,default=1" json:"page" xml:"page"  binding:"number"`
	Limit int `form:"limit,default=5" json:"limit" xml:"limit" binding:"number"`
}
