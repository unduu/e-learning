package http

type Choice struct {
	Type    string   `json:"type"`
	Options []string `json:"options"`
}

type Question struct {
	Id         int    `json:"id"`
	Type       string `json:"type"`
	AttachType string `json:"attachment_type"`
	Attachment string `json:"attachment"`
	Question   string `json:"question"`
	Choices    Choice
}

type PreEvaluationResponse struct {
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Test       []Question
	Pagination PaginationResponse
}

type PaginationResponse struct {
	TotalData   int    `json:"total_data"`
	TotalPage   int    `json:"total_page"`
	Limit       int    `json:"limit"`
	Current     int    `json:"curr_page"`
	PreviousUrl string `json:"prev_page_url"`
	NextUrl     string `json:"next_page_url"`
}
