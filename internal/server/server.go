package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/account"
	"github.com/VladimirKholomyanskyy/gym-api/internal/auth"
	"github.com/VladimirKholomyanskyy/gym-api/internal/progress"
	"github.com/VladimirKholomyanskyy/gym-api/internal/seed"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training"
	"github.com/joho/godotenv"

	// _ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	port                   int
	KeycloakMiddleware     *auth.KeycloakMiddleware
	UserHandler            *account.UserHandler
	ExerciseHandler        *training.ExerciseHandler
	TrainingProgram        *training.TrainingProgramHandler
	WorkoutExerciseHandler *training.WorkoutExerciseHandler
	WorkoutLogsHandler     *progress.WorkoutProgressHandler
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
	userRepo := account.NewUserRepository(db)
	exerciseRepo := training.NewExerciseRepository(db)
	trainingProgramRepo := training.NewTrainingProgramRepository(db)
	workoutRepo := training.NewWorkoutRepository(db)
	workoutExerciseRepo := training.NewWorkoutExerciseRepository(db)
	workoutSessionRepo := progress.NewWorkoutLogRepository(db)
	exerciseLogsRepo := progress.NewExerciseLogRepository(db)

	// Initializing service layer
	userService := account.NewUserService(userRepo)
	trainingProgramService := training.NewTrainingManager(trainingProgramRepo, workoutRepo, workoutExerciseRepo, exerciseRepo)
	workoutProgressManager := progress.NewWorkoutProgressManager(trainingProgramService, workoutSessionRepo, exerciseLogsRepo)

	// Initializing application layer
	userHandler := &account.UserHandler{Service: userService}
	exerciseHandler := training.NewExerciseHandler(trainingProgramService)
	trainingProgramHandler := training.NewTrainingProgramHandler(trainingProgramService)
	workoutExerciseHandler := training.NewWorkoutExerciseHandler(trainingProgramService)
	workoutLogsHandler := progress.NewWorkoutProgressHandler(workoutProgressManager)

	dataSeed := seed.NewDatabaseSeed(exerciseRepo, workoutRepo, trainingProgramRepo, workoutExerciseRepo)
	dataSeed.Seed()

	client_id := os.Getenv("CLIENT_ID")
	issuer := os.Getenv("ISSUER")

	KeycloakMiddleware, err := auth.NewKeycloakMiddleware(userService, issuer, client_id)
	if err != nil {
		log.Fatal("Failed to init keycloak")
	}

	NewServer := &Server{
		port:                   port,
		UserHandler:            userHandler,
		ExerciseHandler:        exerciseHandler,
		TrainingProgram:        trainingProgramHandler,
		WorkoutExerciseHandler: workoutExerciseHandler,
		KeycloakMiddleware:     KeycloakMiddleware,
		WorkoutLogsHandler:     workoutLogsHandler,
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
