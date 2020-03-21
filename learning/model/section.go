package model

type SectionLessons struct {
	Name     string `db:"section_name"`
	Desc     string `db:"section_desc"`
	Type     string `db:"type"`
	Title    string `db:"title"`
	Duration int    `db:"duration"`
	Video    string `db:"content"`
}

type Section struct {
	Name    string
	Desc    string
	Lessons []*Lesson
}

func (s *Section) AddLesson(lesson *Lesson) {
	s.Lessons = append(s.Lessons, lesson)
}
