package http

// Request data
type RequestLogin struct {
	Username string `form:"username" json:"username" xml:"username"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// Request data for register api
type RequestRegister struct {
	Fullname string `form:"fullname" json:"fullname" xml:"fullname"  binding:"required"`
	Phone    string `form:"phone" json:"phone" xml:"phone" binding:"required,isphonenumber,phoneExists"`
	Email    string `form:"email" json:"email" xml:"email" binding:"required,email,emailExists"`
	Username string `form:"username" json:"username" xml:"username" binding:"required,usernameExists"`
	Password string `form:"password" json:"password" xml:"password" binding:"required,min=8"`
}

// Request verify
type RequestVerify struct {
	Code string `form:"code" json:"code" xml:"code"  binding:"required"`
}

// Request forgot password
type RequestForgotPassword struct {
	Phone string `form:"phone" json:"phone" xml:"phone"  binding:"required,isphonenumber,isValidPhone"`
}

// Request forgot password
type RequestResetPassword struct {
	PasswordKey string `form:"password_key" json:"password_key" xml:"password_key"  binding:"required,isValidPaswordKey"`
	PasswordNew string `form:"new_password" json:"new_password" xml:"new_password"  binding:"required,min=8"`
}

// Request resend verification code
type RequestResendVerifCode struct {
	Phone string `form:"phone" json:"phone" xml:"phone"  binding:"required,isphonenumber,isValidPhone"`
}
