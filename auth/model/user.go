package model

type User struct {
	Username   string `db:"username"`
	Password   string `db:"password"`
	Fullname   string `db:"fullname"`
	Phone      string `db:"phone"`
	Email      string `db:"email"`
	Status     string `db:"status"`
	StatusCode int    `db:"status_code"`
	Role       string `db:"role"`
}
