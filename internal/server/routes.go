package server

import (
	"log"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	// Initialize controllers
	AuthAPIController := openapi.NewAuthAPIController(s.AuthHandler)

	ProfileAPIController := openapi.NewProfileAPIController(s.ProfilesHandler)
	SettingsAPIController := openapi.NewSettingsAPIController(s.SettingsHandler)

	TrainingProgramsAPIController := openapi.NewTrainingProgramsAPIController(s.TrainingProgramsHandler)
	WorkoutsAPIController := openapi.NewWorkoutsAPIController(s.WorkoutsHandler)
	WorkoutExercisesAPIController := openapi.NewWorkoutExercisesAPIController(s.WorkoutExercisesHandler)
	ExercisesAPIController := openapi.NewExercisesAPIController(s.ExercisesHandler)
	ScheduledWorkoutsAPIController := openapi.NewScheduledWorkoutsAPIController(s.ScheduledWorkoutsHandler)

	WorkoutSessionsAPIController := openapi.NewWorkoutSessionsAPIController(s.WorkoutSessionsHandler)
	ExerciseLogsApiController := openapi.NewExerciseLogsAPIController(s.ExerciseLogsHandler)

	// Create a new router
	router := mux.NewRouter()
	// Create a subrouter for public endpoints (no authentication required)
	publicRouter := router.PathPrefix("/api/v1").Subrouter()

	// Register the auth config endpoint
	publicRouter.HandleFunc("/auth/config", AuthAPIController.GetAuthConfig).Methods("GET")

	// Create a subrouter for authenticated endpoints
	// All other API endpoints
	authenticatedRouter := openapi.NewRouter(
		ProfileAPIController,
		SettingsAPIController,
		TrainingProgramsAPIController,
		WorkoutsAPIController,
		WorkoutExercisesAPIController,
		ExercisesAPIController,
		ScheduledWorkoutsAPIController,
		WorkoutSessionsAPIController,
		ExerciseLogsApiController,
	)

	// Apply the authentication middleware only to the authenticated router
	authenticatedRouter.Use(s.AuthMiddleware.Authenticate)

	// Mount the authenticated router to the main router
	router.PathPrefix("/api").Handler(authenticatedRouter)

	// Serve Swagger UI
	swaggerDir := "../../swagger-ui" // Path to Swagger UI files
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir(swaggerDir))))

	// Serve OpenAPI YAML file
	router.HandleFunc("/swagger-ui/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		filePath := "../../internal/api/api/openapi.yaml"
		log.Println("Serving OpenAPI YAML from:", filePath) // Debugging the file path
		http.ServeFile(w, r, filePath)
	}).Methods("GET")

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Allow specific origin
		AllowedMethods:   []string{"GET", "PUT", "DELETE", "POST", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Apply CORS middleware to the router
	handler := c.Handler(router)

	return handler
}
