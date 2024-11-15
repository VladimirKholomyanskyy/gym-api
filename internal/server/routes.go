package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/VladimirKholomyanskyy/gym-api/internal/models"
	"github.com/gorilla/mux"
)

var exercises = []models.Exercise{
	{Id: 1, Name: "Bench Press", Description: "Chest workout", PrimaryMuscle: "Chest", SecondaryMuscles: []string{"Triceps", "Shoulders"}, Equipment: "Barbell"},
	{Id: 2, Name: "Deadlift", Description: "Back workout", PrimaryMuscle: "Back", SecondaryMuscles: []string{"Legs", "Core"}, Equipment: "Barbell"},
	// Add more exercises as needed
}

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.HelloWorldHandler)

	// r.HandleFunc("/health", s.healthHandler)

	// r.HandleFunc("/exercises/{id}", s.GetExercise)
	// r.HandleFunc("/exercises", s.GetAllExercises)
	r.HandleFunc("/users", s.UserHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", s.UserHandler.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", s.UserHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", s.UserHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", s.UserHandler.DeleteUser).Methods("DELETE")
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

// func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
// 	jsonResp, err := json.Marshal(s.db.Health())

// 	if err != nil {
// 		log.Fatalf("error handling JSON marshal. Err: %v", err)
// 	}

// 	_, _ = w.Write(jsonResp)
// }

func (s *Server) GetExercise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, exercise := range exercises {
		if int(exercise.Id) == id {
			json.NewEncoder(w).Encode(exercise)
			return
		}
	}

	http.Error(w, "Exercise not found", http.StatusNotFound)
}

func (s *Server) GetAllExercises(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}
	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default page size
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(exercises) {
		json.NewEncoder(w).Encode([]models.Exercise{})
		return
	}
	if end > len(exercises) {
		end = len(exercises)
	}

	paginatedExercises := exercises[start:end]
	json.NewEncoder(w).Encode(paginatedExercises)
}

func (s *Server) CreateTrainingProgram(w http.ResponseWriter, r *http.Request) {
	var newTrainingProgram models.TrainingProgram

	if err := json.NewDecoder(r.Body).Decode(&newTrainingProgram); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newTrainingProgram.Id = rand.Intn(1000)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrainingProgram)
}

func (s *Server) UpdateTrainingProgram(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) GetTrainingProgram(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) GetAllTrainingPrograms(w http.ResponseWriter, r *http.Request) {

}
