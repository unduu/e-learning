package repository

import (
	"fmt"
	"github.com/e-learning/auth/model"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	conn *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		conn: db,
	}
}

func (a *AuthRepository) GetByUsernamePassword(username string, password string) (*model.User, error) {

	menthor := model.User{}

	queryParams := map[string]interface{}{
		"username": username,
		"password": password,
	}
	query, err := a.conn.PrepareNamed(`SELECT username,role FROM users WHERE username=:username AND password = :password`)
	if err != nil {
		fmt.Println("Error db GetByUsernamePassword->PrepareNamed : ", err)
	}

	err = query.Get(&menthor, queryParams)
	if err != nil {
		fmt.Println("Error db GetByUsernamePassword->query.Get : ", err)
	}

	return &menthor, err
}
