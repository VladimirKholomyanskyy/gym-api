package server

import (
	"context"
	"net/http"

	"github.com/VladimirKholomyanskyy/gym-api/internal/handlers"
	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/users", s.UserHandler.CreateUser).Methods("POST")

	// r.Handle("/users", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.UserHandler.GetAllUsers))).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", s.UserHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", s.UserHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", s.UserHandler.DeleteUser).Methods("DELETE")

	r.HandleFunc("/exercises", s.ExerciseHandler.GetAllExercises).Methods("GET")
	// r.Handle("/exercises", s.KeycloakMiddleware.Authenticate(http.HandlerFunc(s.ExerciseHandler.GetAllExercises))).Methods("GET")

	r.Handle("/training-programs", s.HandleUserId(http.HandlerFunc(s.TrainingProgram.HandleCreateProgram))).Methods("POST")
	r.Handle("/training-programs/{id:[0-9]+}", s.HandleUserId(http.HandlerFunc(s.TrainingProgram.HandleDeleteProgram))).Methods("DELETE")
	r.Handle("/training-programs/{id:[0-9]+}", s.HandleUserId(http.HandlerFunc(s.TrainingProgram.HandleGetAllUserPrograms))).Methods("GET")
	r.Handle("/training-programs/{program_id:[0-9]+}/workouts", s.HandleUserId(http.HandlerFunc(s.TrainingProgram.HandleAddWorkoutToProgram))).Methods("POST")
	r.Handle("/training-programs/{program_id:[0-9]+}/workouts/{workout_id:[0-9]+}", s.HandleUserId(http.HandlerFunc(s.TrainingProgram.HandleRemoveWorkoutFromProgram))).Methods("DELETE")
	r.Handle("/training-programs/{program_id:[0-9]+}/workouts/{workout_id:[0-9]+}", s.HandleUserId(http.HandlerFunc(s.TrainingProgram.HandleUpdateWorkoutOfProgram))).Methods("PUT")

	r.Handle("/workout-exercises", s.HandleUserId(http.HandlerFunc(s.WorkoutExerciseHandler.HandleCreateWorkoutExercise))).Methods("POST")
	r.Handle("/workout-exercises", s.HandleUserId(http.HandlerFunc(s.WorkoutExerciseHandler.HandleListWorkoutExercises))).Methods("GET")
	return r
}

func (s *Server) HandleUserId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_id := r.Header.Get("UserId")
		ctx := context.WithValue(r.Context(), handlers.UserIDKey, user_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
