package core

import (
	"log/slog"
	"room-decorator/internal/infra"
	"room-decorator/internal/models"
	"time"

	"github.com/google/uuid"
)

func CreateJob(repo *infra.InMemoryJobRepo, queue *infra.InMemoryQueue, payload string) *models.Job {
	job := &models.Job{
		ID:        uuid.NewString(),
		Status:    models.Pending,
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	repo.Insert(job)
	queue.Enqueue(job.ID)
	return job
}

func ProcessJob(job *models.Job) error {
	slog.Info("processing job", "job_id", job.ID)
	time.Sleep(500 * time.Millisecond)
	slog.Info("finished processing job", "job_id", job.ID)

	return nil
}
