package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/service"
	"github.com/gorilla/mux"
)

type ExerciseHandler struct {
	Service *service.ExerciseService
}

func (h *ExerciseHandler) GetAllExercises(w http.ResponseWriter, r *http.Request) {
	log.Println("Request exercises")
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

func (h *ExerciseHandler) HandleGetExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("Request exercise")
	params := mux.Vars(r)
	exerciseId, err := strconv.Atoi(params["exercise_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}
	exercise, err := h.Service.GetExercise(uint(exerciseId))
	if err != nil {
		http.Error(w, "Failed to fetch exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exercise)
}
