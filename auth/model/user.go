package model

type User struct {
	Username string `db:"username"`
	Role     string `db:"role"`
}
