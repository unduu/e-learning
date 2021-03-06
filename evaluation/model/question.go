package model

type Choice struct {
	Type    string
	Options []string
}

type Question struct {
	Id         int    `db:"id"`
	Module     string `db:"module"`
	Type       string `db:"type"`
	AttachType string `db:"attachment_type"`
	Attachment string `db:"attachment"`
	Text       string `db:"question"`
	Choices    Choice
	ChoicesDB  string `db:"choices"`
	Answer     string `db:"answer"`
}
