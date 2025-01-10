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
		AllowedMethods:   []string{"GET", "PUT", "DELETE", "POST", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	r.HandleFunc("/users/me", s.UserHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/me", s.UserHandler.UpdateUser).Methods("PUT")

	r.HandleFunc("/exercises", s.ExerciseHandler.GetAllExercises).Methods("GET")
	r.HandleFunc("/exercises/{exercise_id:[0-9]+}", s.ExerciseHandler.HandleGetExercise).Methods("GET")
	r.HandleFunc("/training-programs", s.TrainingProgram.HandleCreateProgram).Methods("POST")
	r.HandleFunc("/training-programs/{id:[0-9]+}", s.TrainingProgram.HandleGetProgram).Methods("GET")
	r.HandleFunc("/training-programs/{id:[0-9]+}", s.TrainingProgram.HandleDeleteProgram).Methods("DELETE")
	r.HandleFunc("/training-programs/{program_id:[0-9]+}", s.TrainingProgram.HandleUpdateProgram).Methods("PATCH")
	r.HandleFunc("/training-programs", s.TrainingProgram.HandleGetAllUserPrograms).Methods("GET")
	r.HandleFunc("/training-programs/{program_id:[0-9]+}/workouts", s.TrainingProgram.HandleAddWorkoutToProgram).Methods("POST")
	r.HandleFunc("/training-programs/{program_id:[0-9]+}/workouts", s.TrainingProgram.HandleGetAllWorkouts).Methods("GET")
	r.HandleFunc("/training-programs/{program_id:[0-9]+}/workouts/{workout_id:[0-9]+}", s.TrainingProgram.HandleRemoveWorkoutFromProgram).Methods("DELETE")
	r.HandleFunc("/training-programs/{program_id:[0-9]+}/workouts/{workout_id:[0-9]+}", s.TrainingProgram.HandleUpdateWorkoutOfProgram).Methods("PUT")
	r.HandleFunc("/training-programs/{program_id:[0-9]+}/workouts/{workout_id:[0-9]+}", s.TrainingProgram.HandleGetWorkoutForProgram).Methods("GET")
	r.HandleFunc("/workout-exercises", s.WorkoutExerciseHandler.HandleCreateWorkoutExercise).Methods("POST")
	r.HandleFunc("/workout-exercises", s.WorkoutExerciseHandler.HandleListWorkoutExercises).Methods("GET")
	r.HandleFunc("/workout-exercises/{workout_exercise_id:[0-9]+}", s.WorkoutExerciseHandler.HandlePatchWorkoutExercise).Methods("PATCH")
	r.HandleFunc("/workout-exercises/{workout_exercise_id:[0-9]+}", s.WorkoutExerciseHandler.HandleDeletehWorkoutExercise).Methods("DELETE")

	r.HandleFunc("/workout-sessions", s.WorkoutLogsHandler.HandleCreateWorkoutSession).Methods("POST")
	r.HandleFunc("/workout-sessions", s.WorkoutLogsHandler.HandleGetWorkoutSessions).Methods("GET")
	r.HandleFunc("/workout-sessions/{workout_session_id:[0-9]+}/finish", s.WorkoutLogsHandler.HandleFinishWorkoutSession).Methods("POST")
	r.HandleFunc("/workout-sessions/{workout_session_id:[0-9]+}", s.WorkoutLogsHandler.HandleGetWorkoutSession).Methods("GET")
	r.HandleFunc("/workout-sessions/{workout_session_id:[0-9]+}/logs", s.WorkoutLogsHandler.HandleCreateExerciseLog).Methods("POST")
	r.HandleFunc("/workout-sessions/{workout_session_id:[0-9]+}/logs", s.WorkoutLogsHandler.HandleGetExerciseLogs).Methods("GET")
	r.HandleFunc("/workout-sessions/{workout_session_id:[0-9]+}/logs/{log_id:[0-9]+}", s.WorkoutLogsHandler.HandleGetExerciseLogs).Methods("GET")

	handler := c.Handler(s.KeycloakMiddleware.Authenticate(r))
	return handler
}
