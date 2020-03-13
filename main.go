package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	_authRepo "github.com/unduu/e-learning/auth/repository"
	_authHandler "github.com/unduu/e-learning/auth/transport/http"
	_authUsecase "github.com/unduu/e-learning/auth/usecase"

	_evaluationRepo "github.com/unduu/e-learning/evaluation/repository"
	_evaluationHandler "github.com/unduu/e-learning/evaluation/transport/http"
	_evaluationUsecase "github.com/unduu/e-learning/evaluation/usecase"

	customValidator "github.com/unduu/e-learning/helper/validator"

	"github.com/unduu/e-learning/middleware"
)

func main() {
	dbConn := initDB()
	router := gin.Default()
	m := middleware.NewMiddleware()
	validator := customValidator.NewCustomValidator()

	authRepo := _authRepo.NewAuthRepository(dbConn)
	evalauationRepo := _evaluationRepo.NewEvaluationRepository(dbConn)

	router.Use(Cors())

	s := router.Group("")
	{
		_authUsecase := _authUsecase.NewAuthUsecase(authRepo)
		_authHandler.NewHttpAuthHandler(s, m, validator, _authUsecase)

		_evaluationUsecase := _evaluationUsecase.NewEvaluationUsecase(evalauationRepo)
		_evaluationHandler.NewHttpAuthHandler(s, m, validator, _evaluationUsecase)
	}

	config := cors.DefaultConfig()
	config.AllowMethods = []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "access-control-allow-origin", "access-control-allow-headers", "x-api-key", "x-mock-match-request-body"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))
	router.Use(cors.Default())

	err := http.ListenAndServe(":6000", router)
	if err != nil {
		fmt.Println(err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		c.Next()
	}
}

func initDB() *sqlx.DB {
	var err error

	dbHost := "localhost"
	dbPort := "3306"
	dbUser := "root"
	dbPass := "biteme10"
	dbName := "elearning"
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	db, err := sqlx.Connect(`mysql`, dsn)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	return db
}
