package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
}

type ExerciseRequest struct {
	Name             string   `json:"name"`              // Name of the exercise
	Description      string   `json:"description"`       // Optional description of the exercise
	PrimaryMuscle    string   `json:"primary_muscle"`    // Primary muscle targeted
	SecondaryMuscles []string `json:"secondary_muscles"` // List of secondary muscles targeted
	Equipment        string   `json:"equipment"`
}

type httpServer struct {
}

func newHTTPServer() *httpServer {
	return &httpServer{}
}

func NewHTTPServer(address string) *http.Server {
	httpServer := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/exercises", httpServer.getExercise).Methods("GET")
	return &http.Server{Addr: address, Handler: r}
}

func (server *httpServer) getExercise(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}
