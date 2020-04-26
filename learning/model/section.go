package model

type SectionLessons struct {
	Name      string `db:"section_name"`
	Desc      string `db:"section_desc"`
	Type      string `db:"type"`
	Title     string `db:"title"`
	Permalink string `db:"permalink"`
	Duration  int    `db:"duration"`
	Video     string `db:"content"`
	Timebar   int    `db:"progress_time"`
	LessonID  int    `db:"course_content_id"`
}

type Section struct {
	Name    string
	Desc    string
	Lessons []*Lesson
}

func (s *Section) AddLesson(lesson *Lesson) {
	s.Lessons = append(s.Lessons, lesson)
}
