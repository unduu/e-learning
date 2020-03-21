package model

type User struct {
	Username string `db:"username"`
	Fullname string
}
