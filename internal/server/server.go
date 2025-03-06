package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/account"
	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/auth"
	progresshandlers "github.com/VladimirKholomyanskyy/gym-api/internal/progress/handlers"
	progressrepos "github.com/VladimirKholomyanskyy/gym-api/internal/progress/repository"
	progressusecase "github.com/VladimirKholomyanskyy/gym-api/internal/progress/usecase"
	"github.com/VladimirKholomyanskyy/gym-api/internal/seed"
	traininghandlers "github.com/VladimirKholomyanskyy/gym-api/internal/training/handlers"
	trainingrepos "github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
	trainingusecases "github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
	"github.com/joho/godotenv"

	// _ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	port                     int
	KeycloakMiddleware       *auth.KeycloakMiddleware
	ProfilesHandler          openapi.ProfileAPIServicer
	SettingsHandler          openapi.SettingsAPIServicer
	TrainingProgramsHandler  openapi.TrainingProgramsAPIServicer
	WorkoutsHandler          openapi.WorkoutsAPIServicer
	WorkoutExercisesHandler  openapi.WorkoutExercisesAPIServicer
	ExercisesHandler         openapi.ExercisesAPIServicer
	ScheduledWorkoutsHandler openapi.ScheduledWorkoutsAPIServicer
	WorkoutSessionsHandler   openapi.WorkoutSessionsAPIServicer
	ExerciseLogsHandler      openapi.ExerciseLogsAPIServicer
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
	profilesRepo := account.NewProfileRepository(db)
	settingsRepo := account.NewSettingRepository(db)
	trainingProgramRepo := trainingrepos.NewTrainingProgramRepository(db)
	workoutRepo := trainingrepos.NewWorkoutRepository(db)
	workoutExerciseRepo := trainingrepos.NewWorkoutExerciseRepository(db)
	exerciseRepo := trainingrepos.NewExerciseRepository(db)
	scheduledWorkoutsRepo := trainingrepos.NewScheduledWorkoutRepository(db)
	workoutSessionRepo := progressrepos.NewWorkoutSessionRepository(db)
	exerciseLogsRepo := progressrepos.NewExerciseLogRepository(db)

	// Initializing service layer
	authorization := auth.NewAuthorization(trainingProgramRepo, workoutRepo)
	trainingProgramUseCase := trainingusecases.NewTrainingProgramUseCase(trainingProgramRepo)
	workoutsUseCase := trainingusecases.NewWorkoutUseCase(workoutRepo, authorization)
	workoutExercisesUseCase := trainingusecases.NewWorkoutExerciseUseCase(workoutExerciseRepo, authorization)
	exercisesUseCase := trainingusecases.NewExerciseUseCase(exerciseRepo)
	scheduledWorkoutsUseCase := trainingusecases.NewScheduledWorkoutUseCase(scheduledWorkoutsRepo, authorization)
	workoutSessionsUseCases := progressusecase.NewWorkoutSessionUseCase(workoutSessionRepo, workoutsUseCase)
	exerciseLogsUseCase := progressusecase.NewLogExerciseUseCase(exerciseLogsRepo, exercisesUseCase)
	// Initializing application layer
	profilesHandler := account.NewProfileHandler(profilesRepo)
	settingsHandler := account.NewSettingsHandler(settingsRepo)
	trainingProgramsHandler := traininghandlers.NewTrainingProgramHandler(trainingProgramUseCase)
	workoutsHandler := traininghandlers.NewWorkoutHandler(workoutsUseCase)
	workoutExercisesHandler := traininghandlers.NewWorkoutExerciseHandler(workoutExercisesUseCase)
	exercisesHandler := traininghandlers.NewExerciseHandler(exercisesUseCase)
	scheduledWorkoutsHandler := traininghandlers.NewScheduledWorkoutsHandler(scheduledWorkoutsUseCase)
	workoutSessionsHandler := progresshandlers.NewWorkoutSessionHandler(workoutSessionsUseCases)
	exerciseLogsHandler := progresshandlers.NewExerciseLogHandler(exerciseLogsUseCase)

	dataSeed := seed.NewDatabaseSeed(exerciseRepo, workoutRepo, trainingProgramRepo, workoutExerciseRepo, profilesRepo, settingsRepo)
	dataSeed.Seed()

	client_id := os.Getenv("CLIENT_ID")
	issuer := os.Getenv("ISSUER")

	KeycloakMiddleware, err := auth.NewKeycloakMiddleware(profilesRepo, issuer, client_id)
	if err != nil {
		log.Fatal("Failed to init keycloak")
	}

	NewServer := &Server{
		port:                     port,
		KeycloakMiddleware:       KeycloakMiddleware,
		ProfilesHandler:          profilesHandler,
		SettingsHandler:          settingsHandler,
		TrainingProgramsHandler:  trainingProgramsHandler,
		WorkoutsHandler:          workoutsHandler,
		WorkoutExercisesHandler:  workoutExercisesHandler,
		ExercisesHandler:         exercisesHandler,
		ScheduledWorkoutsHandler: scheduledWorkoutsHandler,
		WorkoutSessionsHandler:   workoutSessionsHandler,
		ExerciseLogsHandler:      exerciseLogsHandler,
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
