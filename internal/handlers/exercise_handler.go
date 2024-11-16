package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/service"
)

type ExerciseHandler struct {
	Service *service.ExerciseService
}

func (h *ExerciseHandler) GetAllExercises(w http.ResponseWriter, r *http.Request) {
	primaryMuscle := r.URL.Query().Get("primary_muscle")
	var exercises []models.Exercise
	var err error
	if primaryMuscle == "" {
		exercises, err = h.Service.GetAllExercises()
	} else {
		exercises, err = h.Service.GetExercisesByPrimaryMuscle(primaryMuscle)

	}
	if err != nil {
		http.Error(w, "Failed to fetch exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exercises)
}
