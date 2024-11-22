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
	var request models.TrainingProgramCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user_id, _ := strconv.Atoi(r.Context().Value(models.UserIDKey).(string))
	program := models.TrainingProgramInput{
		Name:        request.Name,
		Description: request.Description,
		UserID:      uint(user_id),
	}
	responce, err := h.service.CreateTrainingProgram(program)
	if err != nil {
		http.Error(w, "Failed to create training program", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responce)
}

func (h *TrainingProgramHandler) HandleGetAllUserPrograms(w http.ResponseWriter, r *http.Request) {
	user_id, _ := strconv.Atoi(r.Context().Value(models.UserIDKey).(string))
	user_programs, err := h.service.GetAllTrainingPrograms(uint(user_id))
	if err != nil {
		http.Error(w, "Failed to fetch training programs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user_programs)
}

func (h *TrainingProgramHandler) HandleDeleteProgram(w http.ResponseWriter, r *http.Request) {
	user_id, _ := strconv.Atoi(r.Context().Value(models.UserIDKey).(string))
	params := mux.Vars(r)

	program_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteTrainingProgram(uint(user_id), uint(program_id))
	if err != nil {
		http.Error(w, "Failed to delete training program", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *TrainingProgramHandler) HandleUpdateProgram(w http.ResponseWriter, r *http.Request) {
	var request models.TrainingProgramCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user_id, _ := strconv.Atoi(r.Context().Value(models.UserIDKey).(string))
	program := models.TrainingProgramInput{
		Name:        request.Name,
		Description: request.Description,
		UserID:      uint(user_id),
	}
	responce, err := h.service.UpdateTrainingProgram(program)
	if err != nil {
		http.Error(w, "Failed to update training program", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responce)
}

func (h *TrainingProgramHandler) HandleAddWorkoutToProgram(w http.ResponseWriter, r *http.Request) {
	user_id, _ := strconv.Atoi(r.Context().Value(models.UserIDKey).(string))
	params := mux.Vars(r)
	program_id, err := strconv.Atoi(params["program_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
	}
	var request models.CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	workoutInput := models.WorkoutInput{Name: request.Name, TrainingProgramID: uint(program_id), UserID: uint(user_id)}
	var response *models.Workout
	response, err = h.service.AddWorkoutToProgram(workoutInput)
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
