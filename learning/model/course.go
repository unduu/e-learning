package model

type Course struct {
	Id           int
	Title        string
	Subtitle     string
	Alias        string
	TotalLesson  int
	Thumbnail    string
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
		total = len(section.Lessons)
	}
	return total
}
