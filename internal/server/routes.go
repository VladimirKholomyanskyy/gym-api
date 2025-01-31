package server

import (
	"log"
	"net/http"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (s *Server) RegisterRoutes2() http.Handler {
	ExercisesAPIController := openapi.NewExercisesAPIController(s.ExerciseHandler)
	TrainingProgramsAPIController := openapi.NewTrainingProgramsAPIController(s.TrainingProgram)
	WorkoutsAPIController := openapi.NewWorkoutsAPIController(s.TrainingProgram)
	WorkoutExercisesAPIController := openapi.NewWorkoutExercisesAPIController(s.WorkoutExerciseHandler)
	WorkoutSessionsAPIController := openapi.NewWorkoutSessionsAPIController(s.WorkoutLogsHandler)

	r := openapi.NewRouter(ExercisesAPIController, TrainingProgramsAPIController, WorkoutExercisesAPIController, WorkoutSessionsAPIController, WorkoutsAPIController)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Allow specific origin
		AllowedMethods:   []string{"GET", "PUT", "DELETE", "POST", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(s.KeycloakMiddleware.Authenticate(r))
	return handler

}
func (s *Server) RegisterRoutes() http.Handler {
	// Initialize controllers
	ExercisesAPIController := openapi.NewExercisesAPIController(s.ExerciseHandler)
	TrainingProgramsAPIController := openapi.NewTrainingProgramsAPIController(s.TrainingProgram)
	WorkoutsAPIController := openapi.NewWorkoutsAPIController(s.TrainingProgram)
	WorkoutExercisesAPIController := openapi.NewWorkoutExercisesAPIController(s.WorkoutExerciseHandler)
	WorkoutSessionsAPIController := openapi.NewWorkoutSessionsAPIController(s.WorkoutLogsHandler)
	ExerciseLogsApiController := openapi.NewExerciseLogsAPIController(s.WorkoutLogsHandler)

	// Create a new router
	router := mux.NewRouter()
	r := openapi.NewRouter(
		ExercisesAPIController,
		TrainingProgramsAPIController,
		WorkoutExercisesAPIController,
		WorkoutSessionsAPIController,
		WorkoutsAPIController,
		ExerciseLogsApiController,
	)
	r.Use(s.KeycloakMiddleware.Authenticate)
	router.PathPrefix("/api").Handler(r)

	// Serve Swagger UI
	swaggerDir := "../../swagger-ui" // Path to Swagger UI files
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir(swaggerDir))))

	// Serve OpenAPI YAML file
	router.HandleFunc("/swagger-ui/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		filePath := "../../api/openapi-spec.yaml"
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

// fsys, _ := fs.Sub(swaggerContent, "../../swagger-ui")
// 	router.StaticFS("/swagger", http.FS(fsys))
// 	r.
// 	r.HandleFunc("/users/me", s.UserHandler.GetUser).Methods("GET")
// 	r.HandleFunc("/users/me", s.UserHandler.UpdateUser).Methods("PUT")
