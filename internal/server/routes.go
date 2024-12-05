package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Allow specific origin
		AllowedMethods:   []string{"GET", "PUT", "DELETE", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	r.Handle("/users/me", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.UserHandler.GetUser))).Methods("GET")
	r.Handle("/users/me", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.UserHandler.UpdateUser))).Methods("PUT")
	r.HandleFunc("/users/me", s.UserHandler.DeleteUser).Methods("DELETE")

	r.Handle("/exercises", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.ExerciseHandler.GetAllExercises))).Methods("GET")
	// r.Handle("/exercises", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.ExerciseHandler.GetAllExercises))).Methods("GET")

	r.Handle("/training-programs", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.TrainingProgram.HandleCreateProgram))).Methods("POST")
	r.Handle("/training-programs/{id:[0-9]+}", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.TrainingProgram.HandleDeleteProgram))).Methods("DELETE")
	r.Handle("/training-programs", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.TrainingProgram.HandleGetAllUserPrograms))).Methods("GET")
	r.Handle("/training-programs/{program_id:[0-9]+}/workouts", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.TrainingProgram.HandleAddWorkoutToProgram))).Methods("POST")
	r.Handle("/training-programs/{program_id:[0-9]+}/workouts/{workout_id:[0-9]+}", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.TrainingProgram.HandleRemoveWorkoutFromProgram))).Methods("DELETE")
	r.Handle("/training-programs/{program_id:[0-9]+}/workouts/{workout_id:[0-9]+}", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.TrainingProgram.HandleUpdateWorkoutOfProgram))).Methods("PUT")

	r.Handle("/workout-exercises", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.WorkoutExerciseHandler.HandleCreateWorkoutExercise))).Methods("POST")
	r.Handle("/workout-exercises", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.WorkoutExerciseHandler.HandleListWorkoutExercises))).Methods("GET")

	handler := c.Handler(r)
	return handler
}
