package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/VladimirKholomyanskyy/gym-api/internal/service"
	"github.com/gorilla/mux"
)

type TrainingProgramHandler struct {
	service *service.TrainingProgramService
}

func NewTrainingProgramHandler(service *service.TrainingProgramService) *TrainingProgramHandler {
	return &TrainingProgramHandler{service: service}
}
func (h *TrainingProgramHandler) HandleCreateProgram(w http.ResponseWriter, r *http.Request) {
	var request models.CreateTrainingProgramRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(r.Context().Value(UserIDKey).(string))

	response, err := h.service.CreateTrainingProgram(request, uint(userID))
	if err != nil {
		http.Error(w, "Failed to create training program", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TrainingProgramHandler) HandleGetAllUserPrograms(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Context().Value(UserIDKey).(string))
	userPrograms, err := h.service.GetAllTrainingPrograms(uint(userID))
	if err != nil {
		http.Error(w, "Failed to fetch training programs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userPrograms)
}

func (h *TrainingProgramHandler) HandleDeleteProgram(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Context().Value(UserIDKey).(string))
	params := mux.Vars(r)

	programID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteTrainingProgram(uint(userID), uint(programID))
	if err != nil {
		http.Error(w, "Failed to delete training program", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *TrainingProgramHandler) HandleUpdateProgram(w http.ResponseWriter, r *http.Request) {
	var request models.CreateTrainingProgramRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(r.Context().Value(UserIDKey).(string))

	responce, err := h.service.UpdateTrainingProgram(request, uint(userID))
	if err != nil {
		http.Error(w, "Failed to update training program", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responce)
}

func (h *TrainingProgramHandler) HandleAddWorkoutToProgram(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Context().Value(UserIDKey).(string))
	params := mux.Vars(r)
	programID, err := strconv.Atoi(params["program_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}
	var request models.CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var response *models.Workout
	response, err = h.service.AddWorkoutToProgram(request, uint(userID), uint(programID))
	if err != nil {
		http.Error(w, "Failed to add workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TrainingProgramHandler) HandleRemoveWorkoutFromProgram(w http.ResponseWriter, r *http.Request) {

}
func (h *TrainingProgramHandler) HandleUpdateWorkoutOfProgram(w http.ResponseWriter, r *http.Request) {

}
