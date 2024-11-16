package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/handlers"
	"github.com/VladimirKholomyanskyy/gym-api/internal/repository"
	"github.com/VladimirKholomyanskyy/gym-api/internal/seed"
	"github.com/VladimirKholomyanskyy/gym-api/internal/service"
	"github.com/joho/godotenv"

	// _ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	port            int
	UserHandler     *handlers.UserHandler
	ExerciseHandler *handlers.ExerciseHandler
}

func NewServer() *http.Server {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Can't load")
	}
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	fmt.Println(port)
	var (
		database = os.Getenv("BLUEPRINT_DB_DATABASE")
		password = os.Getenv("BLUEPRINT_DB_PASSWORD")
		username = os.Getenv("BLUEPRINT_DB_USERNAME")
		db_port  = os.Getenv("BLUEPRINT_DB_PORT")
		host     = os.Getenv("BLUEPRINT_DB_HOST")
		schema   = os.Getenv("BLUEPRINT_DB_SCHEMA")
	)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, db_port, database, schema)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	user_repo := &repository.UserRepository{DB: db}
	user_svc := &service.UserService{Repo: user_repo}
	user_handler := &handlers.UserHandler{Service: user_svc}

	exercise_repo := &repository.ExerciseRepository{DB: db}
	exercise_svc := &service.ExerciseService{Repo: exercise_repo}
	exercise_handler := &handlers.ExerciseHandler{Service: exercise_svc}

	seed.SeedExercises(exercise_repo)
	NewServer := &Server{
		port:            port,
		UserHandler:     user_handler,
		ExerciseHandler: exercise_handler,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
