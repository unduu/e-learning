package auth

type Usecase interface {
	Login(username string, password string) string
}
