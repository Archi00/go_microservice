package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Archi00/go_microservice/internal/crawler"
	"github.com/Archi00/go_microservice/internal/job"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler bundles our job manager for HTTP endpoints.
type Handler struct {
	JobManager *job.Manager
	// In a complete system, you might also pass configuration, logger, etc.
}

// CreateJobHandler starts a new crawl job in the background.
func (h *Handler) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body. Expect a JSON payload with a "start_url" field.
	var req struct {
		StartURL string `json:"start_url"`
	}
	body, err := os.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.StartURL == "" {
		http.Error(w, "start_url is required", http.StatusBadRequest)
		return
	}

	// Create a new job.
	jobID := uuid.New().String()
	newJob := &job.Job{
		ID:        jobID,
		StartTime: time.Now().UTC(),
		Status:    job.JobStatusRunning,
		StartURL:  req.StartURL,
	}
	h.JobManager.Add(newJob)

	// Start the background job.
	go func(jobID, startURL string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		c := crawler.NewCrawler(startURL)
		results, err := c.Crawl(ctx)
		now := time.Now().UTC()
		j, ok := h.JobManager.Get(jobID)
		if !ok {
			log.Printf("Job %s not found", jobID)
			return
		}
		if err != nil {
			j.Status = job.JobStatusFailed
			j.Error = err.Error()
		} else {
			j.Status = job.JobStatusFinished
			j.Results = results
		}
		j.EndTime = &now
		h.JobManager.Update(j)
	}(jobID, req.StartURL)

	// Respond with the job ID.
	resp := map[string]string{"job_id": jobID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ListJobsHandler returns a list of all jobs.
func (h *Handler) ListJobsHandler(w http.ResponseWriter, r *http.Request) {
	jobs := h.JobManager.List()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// GetJobHandler returns details for a specific job.
func (h *Handler) GetJobHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	j, ok := h.JobManager.Get(id)
	if !ok {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j)
}
