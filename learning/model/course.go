package model

import (
	"strings"
)

type Course struct {
	Id           int
	Title        string `db:"title"`
	Subtitle     string `db:"subtitle"`
	Alias        string `db:"alias"`
	TotalLesson  int
	Thumbnail    string `db:"thumbnail"`
	Participants []*Participant
	Sections     []*Section
}

func (c *Course) AddParticipant(users []*Participant) {
	c.Participants = users
}

func (c *Course) GetParticipantStatus(username string) (status string, code int) {
	for _, p := range c.Participants {
		user := p.User
		if user.Username == username {
			code = p.Status
		}
	}

	switch code {
	case 0:
		status = "Locked"
	case 1:
		status = "Open"
	case 2:
		status = "Completed"

	}

	return status, code
}

func (c *Course) AddSection(section *Section) {
	c.Sections = append(c.Sections, section)
}

func (c *Course) GetSection(name string) *Section {
	for _, section := range c.Sections {
		if section.Name == name {
			return section
		}
	}
	return nil
}

func (c *Course) CountDuration() (total int) {
	for _, section := range c.Sections {
		for _, lesson := range section.Lessons {
			total = total + lesson.Duration
		}
	}
	return total
}

func (c *Course) GetTotalLesson() (total int) {
	for _, section := range c.Sections {
		for _, lesson := range section.Lessons {
			if lesson.Type != "quiz" {
				total = total + 1
			}
		}
	}
	return total
}
func (c *Course) GetStatusCode(status string) (code int) {
	status = strings.ToLower(status)
	switch status {
	case "locked":
		code = 0
	case "open":
		code = 1
	case "completed":
		code = 2
	}
	return code
}

func (c *Course) GetNextCourseId() (id int) {
	return c.Id + 1
}
