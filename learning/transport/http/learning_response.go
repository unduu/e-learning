package http

type Module struct {
	Alias       string `json:"alias"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	TotalLesson string `json:"total_lessons"`
	TotalHours  string `json:"total_hours"`
	Status      string `json:"status"`
	StatusCode  int    `json:"status_code"`
	Thumbnail   string `json:"thumbnail"`
}

type ResponseModuleList struct {
	Modules []Module `json:"modules"`
}

type Lesson struct {
	Type         string `json:"type"`
	Title        string `json:"title"`
	Duration     string `json:"duration"`
	Video        string `json:"video"`
	Progress     string `json:"progress"`
	ProgressCode int    `json:"progress_code"`
}

type Section struct {
	Section string   `json:"section"`
	Name    string   `json:"name"`
	Lessons []Lesson `json:"lessons"`
}

type ResponseLearningContent struct {
	Content []Section `json:"content"`
}
