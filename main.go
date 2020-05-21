package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	_authRepo "github.com/unduu/e-learning/auth/repository"
	_authHandler "github.com/unduu/e-learning/auth/transport/http"
	_authUsecase "github.com/unduu/e-learning/auth/usecase"

	_evaluationRepo "github.com/unduu/e-learning/evaluation/repository"
	_evaluationHandler "github.com/unduu/e-learning/evaluation/transport/http"
	_evaluationUsecase "github.com/unduu/e-learning/evaluation/usecase"

	_learningRepo "github.com/unduu/e-learning/learning/repository"
	_learningHandler "github.com/unduu/e-learning/learning/transport/http"
	_learningUsecase "github.com/unduu/e-learning/learning/usecase"

	customValidator "github.com/unduu/e-learning/helper/validator"

	"github.com/joho/godotenv"
	"github.com/unduu/e-learning/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConn := initDB()
	router := gin.Default()
	m := middleware.NewMiddleware(dbConn)
	validator := customValidator.NewCustomValidator(dbConn)

	authRepo := _authRepo.NewAuthRepository(dbConn)
	evalauationRepo := _evaluationRepo.NewEvaluationRepository(dbConn)
	learningRepo := _learningRepo.NewLearningRepository(dbConn)

	router.Use(Cors())

	authUsecase := _authUsecase.NewAuthUsecase(authRepo)
	evaluationUsecase := _evaluationUsecase.NewEvaluationUsecase(evalauationRepo)
	learningUsecase := _learningUsecase.NewLearningUsecase(learningRepo, evaluationUsecase)

	s := router.Group("")
	{
		_authHandler.NewHttpAuthHandler(s, m, validator, authUsecase)

		_learningHandler.NewHttpLearningHandler(s, m, validator, learningUsecase)

		_evaluationHandler.NewHttpAuthHandler(s, m, validator, evaluationUsecase, learningUsecase)
	}

	config := cors.DefaultConfig()
	config.AllowMethods = []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "access-control-allow-origin", "access-control-allow-headers", "x-api-key", "x-mock-match-request-body"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))
	router.Use(cors.Default())

	err = http.ListenAndServe(os.Getenv("PORT"), router)
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

	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbName := os.Getenv("DBNAME")
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
