package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/unduu/e-learning/auth/model"
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
	query, err := a.conn.PrepareNamed(`SELECT username,role,status,status_code FROM users WHERE username=:username AND password = :password`)
	if err != nil {
		fmt.Println("Error db GetByUsernamePassword->PrepareNamed : ", err)
	}

	err = query.Get(&menthor, queryParams)
	if err != nil {
		fmt.Println("Error db GetByUsernamePassword->query.Get : ", err)
	}

	return &menthor, err
}

func (a *AuthRepository) InsertNewUser(user model.User) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"username":    user.Username,
		"password":    user.Password,
		"role":        user.Role,
		"fullname":    user.Fullname,
		"phone":       user.Phone,
		"email":       user.Email,
		"status":      user.Status,
		"status_code": user.StatusCode,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`INSERT INTO users 
							SET username = :username, password = :password, role = :role, fullname = :fullname, 
								phone = :phone, email = :email, status = :status, status_code = :status_code`)
	if err != nil {
		fmt.Println("Error db InsertAnswer->PrepareNamed : ", err)
	}

	// Execute query
	result, err := query.Exec(queryParams)
	if err != nil {
		fmt.Println("Error db InsertAnswer->query.Get : ", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		fmt.Println("Error db InsertAnswer->RowsAffected : ", err)
	}

	return affected
}
