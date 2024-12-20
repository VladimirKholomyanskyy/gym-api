package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/service"
	"github.com/gorilla/mux"
)

type WorkoutExerciseHandler struct {
	service *service.TrainingProgramService
}

func NewWorkoutExerciseHandler(service *service.TrainingProgramService) *WorkoutExerciseHandler {
	return &WorkoutExerciseHandler{service: service}
}

func (h *WorkoutExerciseHandler) HandleCreateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("Request add exercise to a workout")
	var request models.CreateWorkoutExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(UserIDKey).(uint)

	workoutExercise, err := h.service.AddExerciseToWorkout(request, userID)
	if err != nil {
		http.Error(w, "Failed to create workout exercise", http.StatusInternalServerError)
		return
	}
	response := models.WorkoutExerciseResponse{ID: workoutExercise.ID, WorkoutID: workoutExercise.WorkoutID,
		ExerciseID: workoutExercise.ExerciseID, Sets: workoutExercise.Sets, Reps: workoutExercise.Reps}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *WorkoutExerciseHandler) HandleListWorkoutExercises(w http.ResponseWriter, r *http.Request) {
	log.Println("Request get all exercises of a workout")
	parent := r.URL.Query().Get("parent")
	if parent == "" {
		http.Error(w, "Missing parent parameter", http.StatusBadRequest)
	}

	parts := strings.Split(parent, "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid parent format", http.StatusBadRequest)
	}

	userID := r.Context().Value(UserIDKey).(uint)

	parentType, parentID := parts[0], parts[1]

	if parentType != "workout" {
		http.Error(w, "Not valid parent type", http.StatusBadRequest)
	}

	workoutID, err := strconv.Atoi(parentID)
	if err != nil {
		http.Error(w, "Invalid parent id", http.StatusBadRequest)
	}

	workoutExercises, err := h.service.GetAllWorkoutExercisesByWorkout(userID, uint(workoutID))
	if err != nil {
		http.Error(w, "Failed to fetch workout exercise", http.StatusInternalServerError)
		return
	}

	var response []models.WorkoutExerciseResponse

	for _, workoutExercise := range workoutExercises {
		response = append(response, models.WorkoutExerciseResponse{ID: workoutExercise.ID, WorkoutID: workoutExercise.WorkoutID,
			ExerciseID: workoutExercise.ExerciseID, Sets: workoutExercise.Sets, Reps: workoutExercise.Reps})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *WorkoutExerciseHandler) HandlePatchWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to patch workout exercise")
	userID := r.Context().Value(UserIDKey).(uint)
	params := mux.Vars(r)
	workoutExerciseId, err := strconv.Atoi(params["workout_exercise_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}
	var request models.CreateWorkoutExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	workoutExercise, err := h.service.UpdateWorkoutExercise(request, userID, uint(workoutExerciseId))
	if err != nil {
		http.Error(w, "Failed to update workout exercise", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.WorkoutExerciseResponse{ID: workoutExercise.ID, WorkoutID: workoutExercise.WorkoutID,
		ExerciseID: workoutExercise.ExerciseID, Sets: workoutExercise.Sets, Reps: workoutExercise.Reps})
}

func (h *WorkoutExerciseHandler) HandleDeletehWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("Request delete workout exercise")
	userID := r.Context().Value(UserIDKey).(uint)
	params := mux.Vars(r)
	workoutExerciseId, err := strconv.Atoi(params["workout_exercise_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}
	err = h.service.DeleteWorkoutExercise(userID, uint(workoutExerciseId))
	if err != nil {
		http.Error(w, "Failed to delete workout exercise", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
