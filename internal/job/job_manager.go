package job

import (
	"sync"
	"time"
)

// JobStatus represents the state of a crawl job.
type JobStatus string

const (
	JobStatusRunning  JobStatus = "running"
	JobStatusFinished JobStatus = "finished"
	JobStatusFailed   JobStatus = "failed"
)

// Job holds the details and results of a crawl job.
type Job struct {
	ID        string     `json:"id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Status    JobStatus  `json:"status"`
	StartURL  string     `json:"start_url"`
	Error     string     `json:"error,omitempty"`
	Results   []string   `json:"results,omitempty"` // e.g., crawled URLs
}

// Manager tracks jobs concurrently.
type Manager struct {
	mu   sync.Mutex
	jobs map[string]*Job
}

// NewManager creates a new job manager.
func NewManager() *Manager {
	return &Manager{
		jobs: make(map[string]*Job),
	}
}

// Add registers a new job.
func (m *Manager) Add(job *Job) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs[job.ID] = job
}

// Get retrieves a job by its ID.
func (m *Manager) Get(id string) (*Job, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	job, ok := m.jobs[id]
	return job, ok
}

// List returns all jobs.
func (m *Manager) List() []*Job {
	m.mu.Lock()
	defer m.mu.Unlock()
	jobs := make([]*Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// Update replaces an existing job.
func (m *Manager) Update(job *Job) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs[job.ID] = job
}
