package usecase

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/unduu/e-learning/auth"
	"github.com/unduu/e-learning/auth/model"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type AuthUsecase struct {
	repository auth.Repository
}

func NewAuthUsecase(repository auth.Repository) *AuthUsecase {
	return &AuthUsecase{
		repository: repository,
	}
}

// Login user to system
func (a *AuthUsecase) Login(username string, password string) (user *model.User, tokenString string) {
	passMD5 := md5.Sum([]byte(password))
	passMD5String := fmt.Sprintf("%x", passMD5)
	user, err := a.repository.GetByUsernamePassword(username, passMD5String)

	if err == nil {
		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		expirationTime := time.Now().Add(120 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &model.Claims{
			Username:   user.Username,
			Role:       user.Role,
			StatusCode: user.StatusCode,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}
		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, _ = token.SignedString([]byte(os.Getenv("token_password")))
	}

	return user, tokenString
}

// Register new user
func (a *AuthUsecase) Register(fullname string, phone string, email string, username string, password string) (verifivationCode string, affected int64) {
	vcode := EncodeToString(6)
	passMD5 := md5.Sum([]byte(password))
	user := model.User{
		Username:   username,
		Password:   fmt.Sprintf("%x", passMD5),
		Fullname:   fullname,
		Phone:      phone,
		Email:      email,
		Status:     "inactive",
		StatusCode: 0,
		Role:       "menthor",
	}

	ok := a.repository.InsertNewUser(user, vcode)
	if ok > 0 {
		body := "Thank you for registering to Menthorsip Program. Please enter this verification code to activate your account. Your verification code is " + vcode
		a.SendVerificationCode(vcode, phone, body)
	}
	return vcode, ok
}

// Verify make sure registered email or phone is valid
func (a *AuthUsecase) Verify(username, code string) (affected int64) {
	isValid := a.repository.UpdateUserStatus(username, code)
	return isValid
}

// ForgotPassword send confirmation to reset or create new password
func (a *AuthUsecase) ForgotPassword(phone string) (affected int64, passKey string) {
	passKey = EncodeToString(6)
	isValid := a.repository.InsertPasswordKey(phone, passKey)
	return isValid, passKey
}

// ResetPassword allow user to create new password
func (a *AuthUsecase) ResetPassword(password string, passkey string) (affected int64) {
	passMD5 := md5.Sum([]byte(password))
	isValid := a.repository.UpdateNewPassword(fmt.Sprintf("%x", passMD5), passkey)
	return isValid
}

func (a *AuthUsecase) SendVerificationCode(code string, phone string, body string) {
	// Set account keys & information
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", "+"+phone)
	msgData.Set("From", "+15404015748")
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}

func (a *AuthUsecase) ResendVerificationCode(username string) bool {
	user, _ := a.repository.GetByUsername(username)

	if user.Username == username {
		body := "Thank you for registering to Menthorsip Program. Please enter this verification code to activate your account. Your verification code is " + user.VerificationCode
		a.SendVerificationCode(user.VerificationCode, user.Phone, body)
		return true
	}
	return false
}

// EncodeToString return auto generated 6 digit number
func EncodeToString(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
