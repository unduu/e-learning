package http

type User struct {
	Username   string `json:"username"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}

type LoginResponse struct {
	User  User
	Token string `json:"token"`
}
