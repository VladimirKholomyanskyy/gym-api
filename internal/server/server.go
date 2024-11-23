package server

import (
	"fmt"
	// "log"
	"net/http"
	"os"
	"strconv"
	"time"

	// "github.com/VladimirKholomyanskyy/gym-api/internal/auth"
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
	port int
	// KeycloakMiddleware *auth.KeycloakMiddleware
	UserHandler            *handlers.UserHandler
	ExerciseHandler        *handlers.ExerciseHandler
	TrainingProgram        *handlers.TrainingProgramHandler
	WorkoutExerciseHandler *handlers.WorkoutExerciseHandler
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

	// Initializing data layer
	userRepo := repository.NewUserRepository(db)
	exerciseRepo := repository.NewExerciseRepository(db)
	trainingProgramRepo := repository.NewTrainingProgramRepository(db)
	workoutRepo := repository.NewWorkoutRepository(db)
	workoutExerciseRepo := repository.NewWorkoutExerciseRepository(db)

	// Initializing service layer
	userService := service.NewUserService(userRepo)
	exerciseService := service.NewExerciseService(exerciseRepo)
	trainingProgramService := service.NewTrainingProgramService(trainingProgramRepo, workoutRepo, workoutExerciseRepo)

	// Initializing application layer
	userHandler := &handlers.UserHandler{Service: userService}
	exerciseHandler := &handlers.ExerciseHandler{Service: exerciseService}
	trainingProgramHandler := handlers.NewTrainingProgramHandler(trainingProgramService)
	workoutExerciseHandler := handlers.NewWorkoutExerciseHandler(trainingProgramService)

	seed.SeedExercises(exerciseRepo)

	// client_id := os.Getenv("CLIENT_ID")
	// issuer := os.Getenv("ISSUER")

	// var KeycloakMiddleware *auth.KeycloakMiddleware
	// KeycloakMiddleware, err = auth.NewKeycloakMiddleware(issuer, client_id)
	// if err != nil {
	// 	log.Fatal("Failed to init keycloak")
	// }

	NewServer := &Server{
		port:                   port,
		UserHandler:            userHandler,
		ExerciseHandler:        exerciseHandler,
		TrainingProgram:        trainingProgramHandler,
		WorkoutExerciseHandler: workoutExerciseHandler,
		// KeycloakMiddleware: KeycloakMiddleware,
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
