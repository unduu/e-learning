package http

// Request data
type RequestSaveVideoProgress struct {
	Timebar int `form:"timebar,default=0" json:"timebar" xml:"timebar"  binding:"number"`
}

// Request Add new course
type RequestAddCourse struct {
	Title     string `form:"title" json:"title" xml:"title"  binding:"required"`
	Subtitle  string `form:"subtitle" json:"subtitle" xml:"subtitle"  binding:"required"`
	Thumbnail string `form:"thumbnail" json:"thumbnail" xml:"thumbnail"  binding:"required"`
}

// Request Add new course content
type RequestAddCourseContent struct {
	Section string `form:"section" json:"section" xml:"section"  binding:"required"`
	Name    string `form:"name" json:"name" xml:"name"  binding:"required"`
	Type    string `form:"type" json:"type" xml:"type"  binding:"required"`
	Title   string `form:"title" json:"title" xml:"title"  binding:"required"`
	Video   string `form:"video" json:"video" xml:"video"  binding:"required"`
}
