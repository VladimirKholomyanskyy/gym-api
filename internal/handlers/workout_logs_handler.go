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

type WorkoutLogsHandler struct {
	service *service.TrainingProgramService
}

func NewWorkoutLogsHandler(service *service.TrainingProgramService) *WorkoutLogsHandler {
	return &WorkoutLogsHandler{service: service}
}

func (h *WorkoutLogsHandler) HandleCreateWorkoutSession(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to start workout session")
	var request models.StartWorkoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(UserIDKey).(uint)
	workoutSession, err := h.service.StartWorkoutSession(userID, request)
	if err != nil {
		http.Error(w, "Failed to start workout session", http.StatusInternalServerError)
	}
	var snapshot map[string]interface{}

	err = json.Unmarshal(workoutSession.Snapshot, &snapshot)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := models.StartWorkoutSessionResponse{SessionID: workoutSession.ID, StartedAt: workoutSession.StartedAt.String()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *WorkoutLogsHandler) HandleFinishWorkoutSession(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to finish workout session")
	userID := r.Context().Value(UserIDKey).(uint)
	params := mux.Vars(r)
	sessionID, err := strconv.Atoi(params["workout_session_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	workoutSession, err := h.service.FinishWorkoutSession(userID, uint(sessionID))
	if err != nil {
		http.Error(w, "Failed to start workout session", http.StatusInternalServerError)
	}

	response := models.StartWorkoutSessionResponse{SessionID: workoutSession.ID, StartedAt: workoutSession.StartedAt.String()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *WorkoutLogsHandler) HandleGetWorkoutSessions(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get workout sessions")
	userID := r.Context().Value(UserIDKey).(uint)

	sessions, err := h.service.GetAllWorkoutSessions(userID)
	if err != nil {
		http.Error(w, "Failed to fetch workout session", http.StatusInternalServerError)
		return
	}
	var response []models.GetWorkoutSessionResponse
	for _, session := range sessions {
		var workoutSnapshot models.Workout
		err = json.Unmarshal(session.Snapshot, &workoutSnapshot)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		var completedAt string
		if session.CompletedAt == nil {
			completedAt = ""
		} else {
			completedAt = session.CompletedAt.String()
		}
		item := models.GetWorkoutSessionResponse{SessionID: session.ID, WorkoutSnapshot: workoutSnapshot, StartedAt: session.StartedAt.String(), CompletedAt: completedAt}
		response = append(response, item)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func (h *WorkoutLogsHandler) HandleGetWorkoutSession(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get workout session")
	userID := r.Context().Value(UserIDKey).(uint)
	params := mux.Vars(r)
	sessionID, err := strconv.Atoi(params["workout_session_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	session, err := h.service.GetAllWorkoutSession(userID, uint(sessionID))
	if err != nil {
		http.Error(w, "Failed to fetch workout session", http.StatusInternalServerError)
		return
	}
	var workoutSnapshot models.Workout

	err = json.Unmarshal(session.Snapshot, &workoutSnapshot)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var completedAt string
	if session.CompletedAt == nil {
		completedAt = ""
	} else {
		completedAt = session.CompletedAt.String()
	}
	response := models.GetWorkoutSessionResponse{SessionID: session.ID, WorkoutSnapshot: workoutSnapshot, StartedAt: session.StartedAt.String(), CompletedAt: completedAt}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *WorkoutLogsHandler) HandleCreateExerciseLog(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to log exercise")
	userID := r.Context().Value(UserIDKey).(uint)
	params := mux.Vars(r)
	sessionID, err := strconv.Atoi(params["workout_session_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	var request models.LogExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log, err := h.service.LogExercise(userID, uint(sessionID), request)
	if err != nil {
		http.Error(w, "Failed to log exercise", http.StatusInternalServerError)
	}
	response := models.LogExerciseResponse{LogID: log.ID, ExerciseID: log.ExerciseID, SetNumber: log.SetNumber, RepsCompleted: log.Reps, WeightUsed: log.Weight, LoggedAt: log.LoggedAt.String()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *WorkoutLogsHandler) HandleGetExerciseLog(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get log exercise")
	userID := r.Context().Value(UserIDKey).(uint)
	params := mux.Vars(r)
	sessionID, err := strconv.Atoi(params["workout_session_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	logID, err := strconv.Atoi(params["log_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	log, err := h.service.GetExerciseLog(userID, uint(sessionID), uint(logID))
	if err != nil {
		http.Error(w, "failed to fetch exercise log", http.StatusInternalServerError)
	}
	response := models.LogExerciseResponse{LogID: log.ID, ExerciseID: log.ExerciseID, SetNumber: log.SetNumber, RepsCompleted: log.Reps, WeightUsed: log.Weight, LoggedAt: log.LoggedAt.String()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *WorkoutLogsHandler) HandleGetExerciseLogs(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get logs exercise")
	userID := r.Context().Value(UserIDKey).(uint)
	params := mux.Vars(r)
	sessionID, err := strconv.Atoi(params["workout_session_id"])
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	logs, err := h.service.GetExerciseLogs(userID, uint(sessionID))
	var response []models.LogExerciseResponse
	if err != nil {
		http.Error(w, "failed to fetch exercise log", http.StatusInternalServerError)
	}
	for _, log := range logs {
		response = append(response, models.LogExerciseResponse{LogID: log.ID, ExerciseID: log.ExerciseID, SetNumber: log.SetNumber, RepsCompleted: log.Reps, WeightUsed: log.Weight, LoggedAt: log.LoggedAt.String()})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
