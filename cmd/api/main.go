package main

import (
	"log"
	"net/http"

	"github.com/Archi00/go_microservice/internal/api"
	"github.com/Archi00/go_microservice/internal/job"
	"github.com/gorilla/mux"
)

func main() {
	jobManager := job.NewManager()
	handler := &api.Handler{JobManager: jobManager}

	router := mux.NewRouter()
	router.HandleFunc("/jobs", handler.CreateJobHandler).Methods("POST")
	router.HandleFunc("/jobs", handler.ListJobsHandler).Methods("GET")
	router.HandleFunc("/jobs/{id}", handler.GetJobHandler).Methods("GET")

	addr := ":8080"
	log.Printf("Starting API server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
