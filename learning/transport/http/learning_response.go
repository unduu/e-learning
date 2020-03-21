package http

type Module struct {
	Alias       string `json:"alias"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	TotalLesson string `json:"total_lessons"`
	TotalHours  string `json:"total_hours"`
	Status      string `json:"status"`
	StatusCode  int    `json:"status_code"`
}

type ModulesResponse struct {
	Modules []Module `json:"modules"`
}
