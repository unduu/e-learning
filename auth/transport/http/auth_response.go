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

// Temp response after register a new user contain activation code data
type LoginResponseTemp struct {
	User       User
	Token      string `json:"token"`
	Activation string `json:"activation"`
}

// Temp response after forgot password contain confirmation code data
type ForgotPasswordResponseTemp struct {
	ConfirmationCode string `json:"code"`
}
