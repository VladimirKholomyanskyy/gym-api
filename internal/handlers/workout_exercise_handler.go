package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/service"
)

type WorkoutExerciseHandler struct {
	service *service.TrainingProgramService
}

func NewWorkoutExerciseHandler(service *service.TrainingProgramService) *WorkoutExerciseHandler {
	return &WorkoutExerciseHandler{service: service}
}

func (h *WorkoutExerciseHandler) HandleCreateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	var request models.CreateWorkoutExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(UserIDKey).(uint)

	response, err := h.service.AddExerciseToWorkout(request, userID)
	if err != nil {
		http.Error(w, "Failed to create workout exercise", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *WorkoutExerciseHandler) HandleListWorkoutExercises(w http.ResponseWriter, r *http.Request) {
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
	workoutID, err := strconv.Atoi(parentID)
	if err != nil {
		http.Error(w, "Invalid parent id", http.StatusBadRequest)
	}

	if parentType != "workout" {
		http.Error(w, "Not valid parent type", http.StatusBadRequest)
	}
	workoutExercises, err := h.service.GetAllWorkoutExercisesByWorkout(userID, uint(workoutID))
	if err != nil {
		http.Error(w, "Failed to fetch workout exercise", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workoutExercises)
}
