package infra

import (
	"room-decorator/internal/models"
	"sync"
	"time"
)

type InMemoryJobRepo struct {
	mu   sync.Mutex
	jobs map[string]*models.Job
}

func NewInMemoryJobRepo() *InMemoryJobRepo {
	return &InMemoryJobRepo{jobs: make(map[string]*models.Job)}
}

func (r *InMemoryJobRepo) Get(jobID string) (*models.Job, bool) {
	r.mu.Lock()
	job, ok := r.jobs[jobID]
	r.mu.Unlock()
	return job, ok
}

func (r *InMemoryJobRepo) Insert(job *models.Job) {
	r.mu.Lock()
	r.jobs[job.ID] = job
	r.mu.Unlock()
}

func (r *InMemoryJobRepo) UpdateStatus(jobID string, status models.JobStatus) {
	r.mu.Lock()
	if job, ok := r.jobs[jobID]; ok {
		job.Status = status
		job.UpdatedAt = time.Now().UTC()
	}
	r.mu.Unlock()
}
