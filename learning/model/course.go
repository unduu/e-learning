package model

type Course struct {
	Id           int
	Title        string
	Subtitle     string
	Alias        string
	TotalLesson  int
	Duration     int
	Participants []*Participant
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
