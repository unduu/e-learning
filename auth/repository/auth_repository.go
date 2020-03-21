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
func (a *AuthRepository) InsertNewUser(user model.User, verifCode string) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"username":          user.Username,
		"password":          user.Password,
		"role":              user.Role,
		"fullname":          user.Fullname,
		"phone":             user.Phone,
		"email":             user.Email,
		"status":            user.Status,
		"status_code":       user.StatusCode,
		"verification_code": verifCode,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`INSERT INTO users 
							SET username = :username, password = :password, role = :role, fullname = :fullname, 
								phone = :phone, email = :email, status = :status, status_code = :status_code, verification_code = :verification_code`)
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
func (a *AuthRepository) UpdateUserStatus(username string, code string) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"status":            "active",
		"status_code":       1,
		"verification_code": code,
		"username":          username,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`UPDATE users SET status = :status, status_code = :status_code WHERE verification_code = :verification_code AND username = :username`)
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
func (a *AuthRepository) InsertPasswordKey(phone string, passkey string) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"phone":        phone,
		"password_key": passkey,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`UPDATE users SET password_key = :password_key WHERE phone = :phone`)
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
func (a *AuthRepository) UpdateNewPassword(password string, passkey string) (affected int64) {
	// Data for query
	queryParams := map[string]interface{}{
		"password":     password,
		"password_key": passkey,
	}

	// Compose query
	query, err := a.conn.PrepareNamed(`UPDATE users SET password = :password, password_key = "" WHERE password_key = :password_key`)
	if err != nil {
		fmt.Println("Error db UpdateNewPassword->PrepareNamed : ", err)
	}

	// Execute query
	result, err := query.Exec(queryParams)
	if err != nil {
		fmt.Println("Error db UpdateNewPassword->query.Get : ", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		fmt.Println("Error db UpdateNewPassword->RowsAffected : ", err)
	}

	return affected
}
