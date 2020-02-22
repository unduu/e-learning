package http

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	User  User
	Token string `json:"token"`
}
